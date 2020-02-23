package orm

import (
	"database/sql"
	"errors"
	"fmt"
	"goshop/org/springrain/logger"
	"reflect"
	"strconv"
	"strings"
	"time"
)

//默认的零时时间1970-01-01 00:00:00 +0000 UTC,兼容数据库,避免0001-01-01 00:00:00 +0000 UTC的零值.数据库不让存值,加上1秒,跪了
//因为mysql 5.7后,The TIMESTAMP data type is used for values that contain both date and time parts. TIMESTAMP has a range of '1970-01-01 00:00:01' UTC to '2038-01-19 03:14:07' UTC.
var defaultZeroTime = time.Date(1970, time.January, 1, 0, 0, 1, 0, time.UTC)

//var defaultZeroTime = time.Now()

//注释如果是 . 句号结尾,IDE的提示就截止了,注释结尾不要用 . 结束
//允许的Type
//bug(springrain) 1.需要完善支持的数据类型和赋值接口,例如sql.NullString
//处理基础类型的查询
var allowBaseTypeMap = map[reflect.Kind]bool{
	reflect.String: true,

	reflect.Int:   true,
	reflect.Int8:  true,
	reflect.Int16: true,
	reflect.Int32: true,
	reflect.Int64: true,

	reflect.Uint:   true,
	reflect.Uint8:  true,
	reflect.Uint16: true,
	reflect.Uint32: true,
	reflect.Uint64: true,

	reflect.Float32: true,
	reflect.Float64: true,
}

//数据库操作基类,隔离原生操作数据库API入口,所有数据库操作必须通过BaseDao进行.
type BaseDao struct {
	config     *DataSourceConfig
	dataSource *dataSource
}

//代码只执行一次
//var once sync.Once
//创建baseDao
func NewBaseDao(config *DataSourceConfig) (*BaseDao, error) {
	dataSource, err := newDataSource(config)

	if err != nil {
		err = fmt.Errorf("创建dataSource失败:%w", err)
		logger.Error(err)
		return nil, err
	}

	return &BaseDao{config, dataSource}, nil
}

/*
//代码示例
baseDao.Transaction(func(session *orm.Session) (interface{}, error) {

	//在此编写业务代码


	//return的error如果不为nil,事务就会回滚
    return nil, nil
})
*/
//事务方法,隔离session相关的API.必须通过这个方法进行事务处理,统一事务方式
//return的error如果不为nil,事务就会回滚
func (baseDao *BaseDao) Transaction(doTransaction func(sesion *Session) (interface{}, error)) (interface{}, error) {
	session := baseDao.dataSource.getSession()
	beginerr := session.begin()
	if beginerr != nil {
		beginerr = fmt.Errorf("事务开启失败:%w ", beginerr)
		logger.Error(beginerr)
		return nil, beginerr
	}
	defer func() {
		if r := recover(); r != nil {
			//err = fmt.Errorf("事务开启失败:%w ", err)

			rberr := session.rollback()
			if rberr != nil {
				rberr = fmt.Errorf("recover内事务回滚失败:%w", rberr)
				logger.Error(rberr)
			}

			//记录异常日志
			//if _, ok := r.(runtime.Error); ok {
			//	panic(r)
			//}
			err, errOk := r.(error)
			if errOk {
				err = fmt.Errorf("recover异常:%w", err)
				logger.Panic(err)
			}

		}
	}()

	info, err := doTransaction(session)
	if err != nil {
		err = fmt.Errorf("事务执行失败:%w", err)
		logger.Error(err)
		rberr := session.rollback()
		if rberr != nil {
			rberr = fmt.Errorf("事务回滚失败:%w", rberr)
			logger.Error(rberr)
		}
		return info, err
	}
	commitError := session.commit()
	if commitError != nil {
		commitError = fmt.Errorf("事务提交失败:%w", commitError)
		logger.Error(commitError)
		return info, commitError
	}
	return nil, nil
}

//bug(springrain)日期类型需要测试一下
//根据Finder和封装为指定的entity类型,entity必须是*struct类型或者基础类型的指针.把查询的数据赋值给entity,所以要求指针类型
//如果没有事务,session传入nil.如果有事务,参照使用BaseDao.Transaction方法传入session.请不要自己构建Session
func (baseDao *BaseDao) QueryStruct(session *Session, finder *Finder, entity interface{}) error {
	//禁止外部构建
	if session != nil && session.db == nil {
		return errors.New("如果没有事务,session传入nil.如果有事务,参照使用BaseDao.Transaction方法传入session.请不要自己构建Session")
	}

	checkerr := checkEntityKind(entity)
	if checkerr != nil {
		checkerr = fmt.Errorf("类型检查错误:%w", checkerr)
		logger.Error(checkerr)
		return checkerr
	}
	//获取到sql语句
	sqlstr, err := wrapQuerySQL(baseDao.config.DBType, finder, nil)
	if err != nil {
		err = fmt.Errorf("获取查询语句错误:%w", err)
		logger.Error(err)
		return err
	}

	//根据语句和参数查询
	var rows *sql.Rows
	var e error
	if session != nil {
		rows, e = session.query(sqlstr, finder.Values...)
	} else {
		rows, e = baseDao.dataSource.Query(sqlstr, finder.Values...)
	}

	if e != nil {
		e = fmt.Errorf("查询数据库错误:%w", e)
		logger.Error(e)
		return e
	}

	typeOf := reflect.TypeOf(entity).Elem()

	//数据库返回的列名
	columns, cne := rows.Columns()
	if cne != nil {
		cne = fmt.Errorf("数据库返回列名错误:%w", cne)
		logger.Error(cne)
		return cne
	}

	//如果是基础类型,就查询一个字段
	if allowBaseTypeMap[typeOf.Kind()] && len(columns) == 1 {
		i := 0
		//循环遍历结果集
		for rows.Next() {
			if i > 1 {
				return errors.New("查询出多条数据")
			}
			i++
			scanerr := rows.Scan(entity)
			if scanerr != nil {
				scanerr = fmt.Errorf("rows.Scan异常:%w", scanerr)
				logger.Error(scanerr)
				return scanerr
			}
		}

		return nil
	}

	valueOf := reflect.ValueOf(entity).Elem()
	//获取到类型的字段缓存
	dbColumnFieldMap, dbe := getDBColumnFieldMap(typeOf)
	if dbe != nil {
		dbe = fmt.Errorf("获取字段缓存错误:%w", dbe)
		logger.Error(dbe)
		return dbe
	}
	//声明载体数组,用于存放struct的属性指针
	values := make([]interface{}, len(columns))
	i := 0
	//循环遍历结果集
	for rows.Next() {

		if i > 1 {
			return errors.New("查询出多条数据")
		}
		i++
		//遍历数据库的列名
		for i, column := range columns {
			//从缓存中获取列名的file字段
			field, fok := dbColumnFieldMap[column]
			if !fok { //如果列名不存在,就初始化一个空值
				values[i] = new(interface{})
				continue
			}
			//获取struct的属性值的指针地址
			value := valueOf.FieldByName(field.Name).Addr().Interface()
			//把指针地址放到数组
			values[i] = value
		}
		//scan赋值.是一个指针数组,已经根据struct的属性类型初始化了,sql驱动能感知到参数类型,所以可以直接赋值给struct的指针.这样struct的属性就有值了
		//困扰了我2天,sql驱动真恶心......
		//再说一遍,sql驱动垃圾......
		scanerr := rows.Scan(values...)
		if scanerr != nil {
			scanerr = fmt.Errorf("rows.Scan错误:%w", scanerr)
			logger.Error(scanerr)
			return scanerr
		}

	}

	return nil
}

//根据Finder和封装为指定的entity类型,entity必须是*[]struct类型,已经初始化好的数组,此方法只Append元素,这样调用方就不需要强制类型转换了
//如果没有事务,session传入nil.如果有事务,参照使用BaseDao.Transaction方法传入session.请不要自己构建Session
func (baseDao *BaseDao) QueryStructList(session *Session, finder *Finder, rowsSlicePtr interface{}, page *Page) error {
	if session != nil && session.db == nil { //禁止外部构建
		return errors.New("如果没有事务,session传入nil.如果有事务,参照使用BaseDao.Transaction方法传入session.请不要自己构建Session")
	}
	if rowsSlicePtr == nil { //如果为nil
		return errors.New("数组必须是&[]stuct类型或者基础类型数组的指针")
	}

	pv1 := reflect.ValueOf(rowsSlicePtr)
	if pv1.Kind() != reflect.Ptr { //如果不是指针
		return errors.New("数组必须是&[]stuct类型或者基础类型数组的指针")
	}

	//获取数组元素
	sliceValue := reflect.Indirect(pv1)

	//如果不是数组
	if sliceValue.Kind() != reflect.Slice {
		return errors.New("数组必须是&[]stuct类型或者基础类型数组的指针")
	}
	//获取数组内的元素类型
	sliceElementType := sliceValue.Type().Elem()

	//如果不是struct
	if !(sliceElementType.Kind() == reflect.Struct || allowBaseTypeMap[sliceElementType.Kind()]) {
		return errors.New("数组必须是&[]stuct类型或者基础类型数组的指针")
	}

	sqlstr, err := wrapQuerySQL(baseDao.config.DBType, finder, nil)
	if err != nil {
		err = fmt.Errorf("获取查询语句异常:%w", err)
		logger.Error(err)
		return err
	}
	//根据语句和参数查询
	var rows *sql.Rows
	var e error
	if session != nil {
		rows, e = session.query(sqlstr, finder.Values...)
	} else {
		rows, e = baseDao.dataSource.Query(sqlstr, finder.Values...)
	}
	if e != nil {
		e = fmt.Errorf("查询rows异常:%w", e)
		logger.Error(e)
		return e
	}
	//数据库返回的列名
	columns, cne := rows.Columns()
	if cne != nil {
		cne = fmt.Errorf("数据库返回列名错误:%w", cne)
		logger.Error(cne)
		return cne
	}

	//如果是基础类型,就查询一个字段
	if allowBaseTypeMap[sliceElementType.Kind()] {

		//循环遍历结果集
		for rows.Next() {
			//初始化一个基本类型,new出来的是指针.
			pv := reflect.New(sliceElementType)
			//把数据库值赋给指针
			scanerr := rows.Scan(pv.Interface())
			if scanerr != nil {
				scanerr = fmt.Errorf("rows.Scan异常:%w", scanerr)
				logger.Error(scanerr)
				return scanerr
			}
			//通过反射给slice添加元素.添加指针下的真实元素
			sliceValue.Set(reflect.Append(sliceValue, pv.Elem()))
		}

		//查询总条数
		if page != nil && finder.SelectPageCount {
			count, counterr := baseDao.selectCount(session, finder)
			if counterr != nil {
				counterr = fmt.Errorf("查询总条数错误:%w", counterr)
				logger.Error(counterr)
				return counterr
			}
			page.SetTotalCount(count)
		}
		return nil
	}

	//获取到类型的字段缓存
	dbColumnFieldMap, dbe := getDBColumnFieldMap(sliceElementType)
	if dbe != nil {
		dbe = fmt.Errorf("获取字段缓存错误:%w", dbe)
		logger.Error(dbe)
		return dbe
	}
	//声明载体数组,用于存放struct的属性指针
	values := make([]interface{}, len(columns))
	//循环遍历结果集
	for rows.Next() {
		//deepCopy(a, entity)
		//反射初始化一个数组内的元素
		//new 出来的为什么是个指针啊????
		pv := reflect.New(sliceElementType).Elem()
		//遍历数据库的列名
		for i, column := range columns {
			//从缓存中获取列名的file字段
			field, fok := dbColumnFieldMap[column]
			if !fok { //如果列名不存在,就初始化一个空值
				values[i] = new(interface{})
				continue
			}
			//获取struct的属性值的指针地址
			value := pv.FieldByName(field.Name).Addr().Interface()
			//把指针地址放到数组
			values[i] = value
		}
		//scan赋值.是一个指针数组,已经根据struct的属性类型初始化了,sql驱动能感知到参数类型,所以可以直接赋值给struct的指针.这样struct的属性就有值了
		//困扰了我2天,sql驱动真恶心......
		//再说一遍,sql驱动垃圾......
		scanerr := rows.Scan(values...)
		if scanerr != nil {
			scanerr = fmt.Errorf("rows.Scan异常:%w", scanerr)
			logger.Error(scanerr)
			return scanerr
		}

		//values[i] = f.Addr().Interface()
		//通过反射给slice添加元素
		sliceValue.Set(reflect.Append(sliceValue, pv))
	}

	//查询总条数
	if page != nil && finder.SelectPageCount {
		count, counterr := baseDao.selectCount(session, finder)
		if counterr != nil {
			counterr = fmt.Errorf("查询总条数错误:%w", counterr)
			logger.Error(counterr)
			return counterr
		}
		page.SetTotalCount(count)
	}

	return nil

}

//根据Finder查询,封装Map
//bug(springrain)需要测试一下 in 数组, like ,还有查询一个基础类型(例如 string)的功能
//如果没有事务,session传入nil.如果有事务,参照使用BaseDao.Transaction方法传入session.请不要自己构建Session
func (baseDao *BaseDao) QueryMap(session *Session, finder *Finder) (map[string]interface{}, error) {

	if session != nil && session.db == nil { //禁止外部构建
		return nil, errors.New("如果没有事务,session传入nil.如果有事务,参照使用BaseDao.Transaction方法传入session.请不要自己构建Session")
	}
	resultMapList, listerr := baseDao.QueryMapList(session, finder, nil)
	if listerr != nil {
		listerr = fmt.Errorf("QueryMapList查询错误:%w", listerr)
		logger.Error(listerr)
		return nil, listerr
	}
	if resultMapList == nil {
		return nil, nil
	}
	if len(resultMapList) > 1 {
		return resultMapList[0], errors.New("查询出多条数据")
	}
	return resultMapList[0], nil
}

//根据Finder查询,封装Map数组
//根据数据库字段的类型,完成从[]byte到golang类型的映射,理论上其他方法都可以调用这一个方法(需要处理sql.Nullxxx等驱动支持的类型)
//如果没有事务,session传入nil.如果有事务,参照使用BaseDao.Transaction方法传入session.请不要自己构建Session
func (baseDao *BaseDao) QueryMapList(session *Session, finder *Finder, page *Page) ([]map[string]interface{}, error) {

	if session != nil && session.db == nil { //禁止外部构建
		return nil, errors.New("如果没有事务,session传入nil.如果有事务,参照使用BaseDao.Transaction方法传入session.请不要自己构建Session")
	}
	sqlstr, err := wrapQuerySQL(baseDao.config.DBType, finder, nil)
	if err != nil {
		err = fmt.Errorf("QueryMapList查询错误:%w", err)
		logger.Error(err)
		return nil, err
	}

	//根据语句和参数查询
	rows, e := baseDao.dataSource.Query(sqlstr, finder.Values...)
	if e != nil {
		e = fmt.Errorf("查询rows错误:%w", e)
		logger.Error(e)
		return nil, e
	}

	//数据库返回的列名
	//columns, cne := rows.Columns()
	//妈蛋,columnType.scanType返回的类型都是[]byte......
	//如果使用columnType.databaseType挨个判断也是去了意义,只要字段允许null,都是[]byte
	columnTypes, cne := rows.ColumnTypes()
	if cne != nil {
		cne = fmt.Errorf("数据库返回列名错误:%w", cne)
		logger.Error(cne)
		return nil, cne
	}
	resultMapList := make([]map[string]interface{}, 0)
	//循环遍历结果集
	for rows.Next() {
		//接收数据库返回的数据,需要使用指针接收,以前使用[]byte接收,无法接收NULL值.无法获取sql的metadata,比较恶心......
		values := make([]interface{}, len(columnTypes))
		//使用指针类型接收字段值,需要使用interface{}包装一下
		result := make(map[string]interface{})
		//给数据赋值初始化变量
		for i := range values {
			values[i] = new(interface{})
		}
		//scan赋值
		scanerr := rows.Scan(values...)
		if scanerr != nil {
			scanerr = fmt.Errorf("rows.Scan异常:%w", scanerr)
			logger.Error(scanerr)
			return nil, scanerr
		}
		//获取每一列的值
		for i, columnType := range columnTypes {

			//取到指针下的值,[]byte格式
			v := *(values[i].(*interface{}))
			//从[]byte转化成实际的类型值,例如string,int
			v = converValueColumnType(v, columnType)
			//赋值到Map
			result[columnType.Name()] = v

		}

		//添加Map到数组
		resultMapList = append(resultMapList, result)

	}

	//bug(springrain) 还缺少查询总条数的逻辑
	//查询总条数
	if page != nil && finder.SelectPageCount {
		count, counterr := baseDao.selectCount(session, finder)
		if counterr != nil {
			counterr = fmt.Errorf("查询总条数错误:%w", counterr)
			logger.Error(counterr)
			return resultMapList, counterr
		}
		page.SetTotalCount(count)
	}

	return resultMapList, nil
}

//更新Finder
//session不能为nil,参照使用BaseDao.Transaction方法传入session.请不要自己构建Session
func (baseDao *BaseDao) UpdateFinder(session *Session, finder *Finder) error {

	//必须要有session和事务
	if session == nil || session.tx == nil {
		return errors.New("session不能为nil,参照使用BaseDao.Transaction方法传入session.请不要自己构建Session")
	}

	if finder == nil {
		return errors.New("finder不能为空")
	}
	sqlstr, err := finder.GetSQL()
	if err != nil {
		err = fmt.Errorf("finder.GetSQL()错误:%w", err)
		logger.Error(err)
		return err
	}
	sqlstr, err = wrapSQL(baseDao.config.DBType, sqlstr)
	if err != nil {
		err = fmt.Errorf("UpdateFinder-->wrapSQL错误:%w", err)
		logger.Error(err)
		return err
	}
	//流弊的...,把数组展开变成多个参数的形式
	_, errexec := session.exec(sqlstr, finder.Values...)

	if errexec != nil {
		errexec = fmt.Errorf("执行更新错误:%w", errexec)
		logger.Error(errexec)
		return errexec
	}
	return nil
}

//保存Struct对象,必须是IEntityStruct类型
//bug(chunanuyong) 如果是自增主键,需要返回.需要sql驱动支持
//session不能为nil,参照使用BaseDao.Transaction方法传入session.请不要自己构建Session
func (baseDao *BaseDao) SaveStruct(session *Session, entity IEntityStruct) error {
	//必须要有session和事务
	if session == nil || session.tx == nil {
		return errors.New("session不能为nil,参照使用BaseDao.Transaction方法传入session.请不要自己构建Session")
	}
	if entity == nil {
		return errors.New("对象不能为空")
	}
	columns, values, columnAndValueErr := columnAndValue(entity)
	if columnAndValueErr != nil {
		columnAndValueErr = fmt.Errorf("SaveStruct-->columnAndValue获取实体类的列和值异常:%w", columnAndValueErr)
		logger.Error(columnAndValueErr)
		return columnAndValueErr
	}
	if len(columns) < 1 {
		return errors.New("没有tag信息,请检查struct中 column 的tag")
	}
	//SQL语句
	sqlstr, autoIncrement, err := wrapSaveStructSQL(baseDao.config.DBType, entity, &columns, &values)
	if err != nil {
		err = fmt.Errorf("SaveStruct-->wrapSaveStructSQL获取保存语句错误:%w", err)
		logger.Error(err)
		return err
	}

	//流弊的...,把数组展开变成多个参数的形式
	res, errexec := session.exec(sqlstr, values...)

	if errexec != nil {
		errexec = fmt.Errorf("SaveStruct执行保存错误:%w", errexec)
		logger.Error(errexec)
		return errexec
	}
	//如果是自增主键
	if autoIncrement {
		//需要数据库支持,获取自增主键
		autoIncrementIdInt64, e := res.LastInsertId()
		if e != nil { //数据库不支持自增主键,不再赋值给struct属性
			e = fmt.Errorf("数据库不支持自增主键,不再赋值给struct属性:%w", e)
			logger.Error(e)
			return nil
		}
		pkName := entity.GetPKColumnName()
		//int64 转 int
		strInt64 := strconv.FormatInt(autoIncrementIdInt64, 10)
		autoIncrementIdInt, _ := strconv.Atoi(strInt64)
		//设置自增主键的值
		seterr := setFieldValueByColumnName(entity, pkName, autoIncrementIdInt)
		if seterr != nil {
			seterr = fmt.Errorf("反射赋值数据库返回的自增主键错误:%w", seterr)
			logger.Error(seterr)
			return seterr
		}
	}

	return nil

}

//更新struct所有属性,必须是IEntityStruct类型
//session不能为nil,参照使用BaseDao.Transaction方法传入session.请不要自己构建Session
func (baseDao *BaseDao) UpdateStruct(session *Session, entity IEntityStruct) error {
	err := baseDao.updateStructFunc(session, entity, false)
	if err != nil {
		err = fmt.Errorf("UpdateStruct-->updateStructFunc更新错误:%w", err)
		return err
	}
	return nil
}

//更新struct不为nil的属性,必须是IEntityStruct类型
//session不能为nil,参照使用BaseDao.Transaction方法传入session.请不要自己构建Session
func (baseDao *BaseDao) UpdateStructNotNil(session *Session, entity IEntityStruct) error {
	err := baseDao.updateStructFunc(session, entity, true)
	if err != nil {
		err = fmt.Errorf("UpdateStructNotNil-->updateStructFunc更新错误:%w", err)
		return err
	}
	return nil
}

// 根据主键删除一个对象.必须是IEntityStruct类型
//session不能为nil,参照使用BaseDao.Transaction方法传入session.请不要自己构建Session
func (baseDao *BaseDao) DeleteStruct(session *Session, entity IEntityStruct) error {
	//必须要有session和事务
	if session == nil || session.tx == nil {
		return errors.New("session不能为nil,参照使用BaseDao.Transaction方法传入session.请不要自己构建Session")
	}
	if entity == nil {
		return errors.New("对象不能为空")
	}
	pkName, pkNameErr := entityPKFieldName(entity)
	if pkNameErr != nil {
		pkNameErr = fmt.Errorf("DeleteStruct-->entityPKFieldName获取主键名称错误:%w", pkNameErr)
		logger.Error(pkNameErr)
		return pkNameErr
	}

	value, e := structFieldValue(entity, pkName)
	if e != nil {
		e = fmt.Errorf("DeleteStruct-->structFieldValue获取主键值错误:%w", e)
		logger.Error(e)
		return e
	}
	//SQL语句
	sqlstr, err := wrapDeleteStructSQL(baseDao.config.DBType, entity)
	if err != nil {
		err = fmt.Errorf("DeleteStruct-->wrapDeleteStructSQL获取SQL语句错误:%w", err)
		logger.Error(err)
		return err
	}

	_, errexec := session.exec(sqlstr, value)

	if errexec != nil {
		errexec = fmt.Errorf("DeleteStruct执行删除错误:%w", errexec)
		logger.Error(errexec)
		return errexec
	}

	return nil

}

//保存IEntityMap对象.使用Map保存数据,需要在数据中封装好包括Id在内的所有数据.不适用于复杂情况
//session不能为nil,参照使用BaseDao.Transaction方法传入session.请不要自己构建Session
func (baseDao *BaseDao) SaveMap(session *Session, entity IEntityMap) error {
	//必须要有session和事务
	if session == nil || session.tx == nil {
		return errors.New("session不能为nil,参照使用BaseDao.Transaction方法传入session.请不要自己构建Session")
	}
	if entity == nil {
		return errors.New("对象不能为空")
	}
	//SQL语句
	sqlstr, values, err := wrapSaveMapSQL(baseDao.config.DBType, entity)
	if err != nil {
		err = fmt.Errorf("SaveMap-->wrapSaveMapSQL获取SQL语句错误:%w", err)
		logger.Error(err)
		return err
	}

	//流弊的...,把数组展开变成多个参数的形式

	_, errexec := session.exec(sqlstr, values...)

	if errexec != nil {
		errexec = fmt.Errorf("SaveMap执行保存错误:%w", errexec)
		logger.Error(errexec)
		return errexec
	}

	return nil

}

//更新IEntityMap对象.使用Map修改数据,需要在数据中封装好包括Id在内的所有数据.不适用于复杂情况
//session不能为nil,参照使用BaseDao.Transaction方法传入session.请不要自己构建Session
func (baseDao *BaseDao) UpdateMap(session *Session, entity IEntityMap) error {
	//必须要有session和事务
	if session == nil || session.tx == nil {
		return errors.New("session不能为nil,参照使用BaseDao.Transaction方法传入session.请不要自己构建Session")
	}
	if entity == nil {
		return errors.New("对象不能为空")
	}
	//SQL语句
	sqlstr, values, err := wrapUpdateMapSQL(baseDao.config.DBType, entity)
	if err != nil {
		err = fmt.Errorf("UpdateMap-->wrapUpdateMapSQL获取SQL语句错误:%w", err)
		logger.Error(err)
		return err
	}
	//fmt.Println(sqlstr)
	//流弊的...,把数组展开变成多个参数的形式
	_, errexec := session.exec(sqlstr, values...)

	if errexec != nil {
		errexec = fmt.Errorf("UpdateMap执行更新错误:%w", errexec)
		logger.Error(errexec)
		return errexec
	}
	//fmt.Println(entity.GetTableName() + " update success")
	return nil

}

//根据保存的对象,返回插入的语句,需要插入的字段,字段的值.
func columnAndValue(entity interface{}) ([]reflect.StructField, []interface{}, error) {
	checkerr := checkEntityKind(entity)
	if checkerr != nil {
		return nil, nil, checkerr
	}
	// 获取实体类的反射,指针下的struct
	valueOf := reflect.ValueOf(entity).Elem()
	//reflect.Indirect

	//先从本地缓存中查找
	typeOf := reflect.TypeOf(entity).Elem()

	dbMap, err := getDBColumnFieldMap(typeOf)
	if err != nil {
		return nil, nil, err
	}

	//实体类公开字段的长度
	fLen := len(dbMap)
	//接收列的数组
	columns := make([]reflect.StructField, 0, fLen)
	//接收值的数组
	values := make([]interface{}, 0, fLen)

	//遍历所有数据库属性
	for _, field := range dbMap {
		//获取字段类型的Kind
		//	fieldKind := field.Type.Kind()
		//if !allowTypeMap[fieldKind] { //不允许的类型
		//	continue
		//}

		columns = append(columns, field)
		//FieldByName方法返回的是reflect.Value类型,调用Interface()方法,返回原始类型的数据值
		value := valueOf.FieldByName(field.Name).Interface()
		if value != nil { //如果不是nil
			timeValue, ok := value.(time.Time)
			if ok && timeValue.IsZero() { //如果是日期零时,需要设置一个初始值1970-01-01 00:00:01,兼容数据库
				value = defaultZeroTime
			}
		}

		//添加到记录值的数组
		values = append(values, value)

	}

	//缓存数据库的列

	return columns, values, nil

}

//获取实体类主键属性名称
func entityPKFieldName(entity IEntityStruct) (string, error) {
	//缓存的key,TypeOf和ValueOf的String()方法,返回值不一样
	typeOf := reflect.TypeOf(entity).Elem()

	dbMap, err := getDBColumnFieldMap(typeOf)
	if err != nil {
		return "", err
	}
	field := dbMap[entity.GetPKColumnName()]
	return field.Name, nil

}

//检查entity类型必须是*struct类型或者基础类型的指针
func checkEntityKind(entity interface{}) error {
	if entity == nil {
		return errors.New("参数不能为空,必须是*struct类型或者基础类型的指针")
	}
	typeOf := reflect.TypeOf(entity)
	if typeOf.Kind() != reflect.Ptr { //如果不是指针
		return errors.New("必须是*struct类型或者基础类型的指针")
	}
	typeOf = typeOf.Elem()
	if !(typeOf.Kind() == reflect.Struct || allowBaseTypeMap[typeOf.Kind()]) { //如果不是指针
		return errors.New("必须是*struct类型或者基础类型的指针")
	}
	return nil
}

//根据数据库返回的sql.Rows,查询出列名和对应的值.废弃
/*
func columnValueMap2Struct(resultMap map[string]interface{}, typeOf reflect.Type, valueOf reflect.Value) error {


		dbMap, err := getDBColumnFieldMap(typeOf)
		if err != nil {
			return err
		}

		for column, columnValue := range resultMap {
			field, ok := dbMap[column]
			if !ok {
				continue
			}
			fieldName := field.Name
			if len(fieldName) < 1 {
				continue
			}
			//反射获取字段的值对象
			fieldValue := valueOf.FieldByName(fieldName)
			//获取值类型
			kindType := fieldValue.Kind()
			valueType := fieldValue.Type()
			if kindType == reflect.Ptr { //如果是指针类型的属性,查找指针下的类型
				kindType = fieldValue.Elem().Kind()
				valueType = fieldValue.Elem().Type()
			}
			kindTypeStr := kindType.String()
			valueTypeStr := valueType.String()
			var v interface{}
			if kindTypeStr == "string" || valueTypeStr == "string" { //兼容string的扩展类型
				v = columnValue.String()
			} else if kindTypeStr == "int" || valueTypeStr == "int" { //兼容int的扩展类型
				v = columnValue.Int()
			}
			//bug(springrain)这个地方还要添加其他类型的判断,参照ColumnValue.go文件

			fieldValue.Set(reflect.ValueOf(v))

		}

	return nil

}
*/
//根据sql查询结果,返回map.废弃
/*
func wrapMap(columns []string, values []columnValue) (map[string]columnValue, error) {
	columnValueMap := make(map[string]columnValue)
	for i, column := range columns {
		columnValueMap[column] = values[i]
	}
	return columnValueMap, nil
}
*/

//更新对象
//session不能为nil,参照使用BaseDao.Transaction方法传入session.请不要自己构建Session
func (baseDao *BaseDao) updateStructFunc(session *Session, entity IEntityStruct, onlyupdatenotnull bool) error {
	//必须要有session和事务
	if session == nil || session.tx == nil {
		return errors.New("session不能为nil,参照使用BaseDao.Transaction方法传入session.请不要自己构建Session")
	}
	if entity == nil {
		return errors.New("对象不能为空")
	}
	columns, values, columnAndValueErr := columnAndValue(entity)
	if columnAndValueErr != nil {
		return columnAndValueErr
	}
	//SQL语句
	sqlstr, err := wrapUpdateStructSQL(baseDao.config.DBType, entity, columns, values, onlyupdatenotnull)
	if err != nil {
		return err
	}

	//流弊的...,把数组展开变成多个参数的形式
	_, errexec := session.exec(sqlstr, values...)

	if errexec != nil {
		return errexec
	}

	//fmt.Println(entity.GetTableName() + " update success")
	return nil

}

//根据finder查询总条数
//如果没有事务,session传入nil.如果有事务,参照使用BaseDao.Transaction方法传入session.请不要自己构建Session
func (baseDao *BaseDao) selectCount(session *Session, finder *Finder) (int, error) {
	//禁止外部构建
	if session != nil && session.db == nil {
		return -1, errors.New("如果没有事务,session传入nil.如果有事务,参照使用BaseDao.Transaction方法传入session.请不要自己构建Session")
	}
	if finder == nil {
		return -1, errors.New("参数为nil")
	}
	//自定义的查询总条数Finder,主要是为了在group by等复杂情况下,为了性能,手动编写总条数语句
	if finder.CountFinder != nil {
		count := -1
		err := baseDao.QueryStruct(session, finder.CountFinder, &count)
		if err != nil {
			return -1, err
		}
		return count, nil
	}

	countsql, counterr := finder.GetSQL()
	if counterr != nil {
		return -1, counterr
	}

	//查询order by 的位置
	locOrderBy := findOrderByIndex(countsql)
	if len(locOrderBy) > 0 { //如果存在order by
		countsql = countsql[:locOrderBy[0]]
	}
	s := strings.ToLower(countsql)
	gbi := -1
	locGroupBy := findGroupByIndex(countsql)
	if len(locGroupBy) > 0 {
		gbi = locGroupBy[0]
	}
	//特殊关键字,包装SQL
	if strings.Index(s, " distinct ") > -1 || strings.Index(s, " union ") > -1 || gbi > -1 {
		countsql = "SELECT COUNT(*)  frame_row_count FROM (" + countsql + ") temp_frame_noob_table_name WHERE 1=1 "
	} else {
		locFrom := findFromIndex(countsql)
		//没有找到FROM关键字,认为是异常语句
		if len(locFrom) < 0 {
			return -1, errors.New("没有FROM关键字,语句错误")
		}
		countsql = "SELECT COUNT(*) " + countsql[locFrom[0]:]
	}

	countFinder := NewFinder()
	countFinder.Append(countsql)
	countFinder.Values = finder.Values

	count := -1
	cerr := baseDao.QueryStruct(session, countFinder, &count)
	if cerr != nil {
		return -1, cerr
	}
	return count, nil

}
