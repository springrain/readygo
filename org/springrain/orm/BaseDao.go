package orm

import (
	"context"
	"errors"
	"reflect"
)

//允许的Type
//bug(chunanyong) 1.需要完善支持的数据类型和赋值接口,例如sql.NullString.
//废弃,是否支持让数据库自己抛错吧
/*
var allowTypeMap = map[reflect.Kind]bool{
	reflect.Float32: true,
	reflect.Float64: true,
	reflect.Int:     true,
	reflect.String:  true,
}
*/

//数据库操作基类,隔离原生操作数据库API入口,所有数据库操作必须通过BaseDao进行.
type BaseDao struct {
	config     *DataSourceConfig
	dataSource *dataSource
}

//创建baseDao
func NewBaseDao(config *DataSourceConfig) (*BaseDao, error) {
	dataSource, err := newDataSource(config)
	return &BaseDao{config, dataSource}, err
}

//根据Finder和封装为指定的entity类型,entity必须是*struct类型.把查询的数据赋值给entity,所以要求指针类型
func (baseDao *BaseDao) QueryStruct(finder *Finder, entity interface{}) error {

	checkerr := checkEntityKind(entity)
	if checkerr != nil {
		return checkerr
	}
	//获取到sql语句
	sqlstr, err := wrapQuerySQL(baseDao.config.DBType, finder, nil)
	if err != nil {
		return err
	}

	//根据语句和参数查询
	rows, e := baseDao.dataSource.Query(sqlstr, finder.Values...)
	if e != nil {
		return e
	}

	//数据库返回的列名
	columns, cne := rows.Columns()
	if cne != nil {
		return cne
	}

	typeOf := reflect.TypeOf(entity).Elem()
	valueOf := reflect.ValueOf(entity).Elem()
	//获取到类型的字段缓存
	dbColumnFieldMap, e := getDBColumnFieldMap(typeOf)
	if e != nil {
		return err
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
		err = rows.Scan(values...)
		if err != nil {
			return err
		}

	}

	return nil
}

//bug(chunanyong)数据库字段为null,映射异常
//bug(chunanyong)需要处理查询总条数的逻辑
//bug(chunanyong)需要处理查询一个基础类型的情况,例如 int,[]int
//根据Finder和封装为指定的entity类型,entity必须是*[]struct类型,已经初始化好的数组,此方法只Append元素,这样调用方就不需要强制类型转换了.
func (baseDao *BaseDao) QueryStructList(finder *Finder, rowsSlicePtr interface{}, page *Page) error {

	if rowsSlicePtr == nil { //如果为nil
		return errors.New("数组必须是&[]stuct类型")
	}

	pv1 := reflect.ValueOf(rowsSlicePtr)
	if pv1.Kind() != reflect.Ptr { //如果不是指针
		return errors.New("数组必须是&[]stuct类型")
	}

	//获取数组元素
	sliceValue := reflect.Indirect(pv1)

	//如果不是数组
	if sliceValue.Kind() != reflect.Slice {
		return errors.New("数组必须是&[]stuct类型")
	}
	//获取数组内的元素类型
	sliceElementType := sliceValue.Type().Elem()

	//如果不是struct
	if sliceElementType.Kind() != reflect.Struct {
		return errors.New("数组必须是&[]stuct类型")
	}

	sqlstr, err := wrapQuerySQL(baseDao.config.DBType, finder, nil)
	if err != nil {
		return err
	}
	//根据语句和参数查询
	rows, e := baseDao.dataSource.Query(sqlstr, finder.Values...)
	if e != nil {
		return e
	}

	//数据库返回的列名
	columns, cne := rows.Columns()
	if cne != nil {
		return cne
	}

	//获取到类型的字段缓存
	dbColumnFieldMap, e := getDBColumnFieldMap(sliceElementType)
	if e != nil {
		return err
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
		err = rows.Scan(values...)
		if err != nil {
			return err
		}

		//values[i] = f.Addr().Interface()
		//通过反射给slice添加元素
		sliceValue.Set(reflect.Append(sliceValue, pv))
	}
	return nil

}

//根据Finder查询,封装Map.获取具体的值,需要自己根据类型调用ColumnValue的转化方法,例如ColumnValue.String()
//golang的sql驱动不支持获取到数据字段的metadata......垃圾.....
//bug(chunanyong)需要测试一下 in 数组, like ,还有查询一个基础类型(例如 string)的功能
func (baseDao *BaseDao) QueryMap(finder *Finder) (map[string]interface{}, error) {
	resultMapList, err := baseDao.QueryMapList(finder, nil)
	if err != nil {
		return nil, err
	}
	if resultMapList == nil {
		return nil, err
	}
	if len(resultMapList) > 1 {
		return resultMapList[0], errors.New("查询出多条数据")
	}
	return resultMapList[0], nil
}

//根据Finder查询,封装Map数组.获取具体的值,需要自己根据类型调用ColumnValue的转化方法,例如ColumnValue.String()
//bug(chunanyong)需要处理查询总条数的逻辑
//golang的sql驱动不支持获取到数据字段的metadata......垃圾.....
func (baseDao *BaseDao) QueryMapList(finder *Finder, page *Page) ([]map[string]interface{}, error) {

	sqlstr, err := wrapQuerySQL(baseDao.config.DBType, finder, nil)
	if err != nil {
		return nil, err
	}

	//根据语句和参数查询
	rows, e := baseDao.dataSource.Query(sqlstr, finder.Values...)
	if e != nil {
		return nil, e
	}

	//数据库返回的列名
	columns, cne := rows.Columns()
	if cne != nil {
		return nil, cne
	}
	resultMapList := make([]map[string]interface{}, 0)
	//循环遍历结果集
	for rows.Next() {
		//接收数据库返回的数据,需要使用指针接收,以前使用[]byte接收,无法接收NULL值.无法获取sql的metadata,比较恶心......
		values := make([]interface{}, len(columns))
		//使用指针类型接收字段值,需要使用interface{}包装一下
		result := make(map[string]interface{})
		//给数据赋值初始化变量
		for i := range values {
			values[i] = new(interface{})
		}
		//scan赋值
		err = rows.Scan(values...)
		if err != nil {
			return nil, err
		}
		//获取每一列的值
		for i, column := range columns {
			//获取指针下的真实值,赋值到map
			result[column] = *(values[i].(*interface{}))
		}

		//添加Map到数组
		resultMapList = append(resultMapList, result)

	}

	//bug(chunanyong) 还缺少查询总条数的逻辑
	return resultMapList, nil
}

//更新Finder
func (baseDao *BaseDao) UpdateFinder(finder *Finder) error {
	if finder == nil {
		return errors.New("finder不能为空")
	}
	sqlstr, err := finder.GetSQL()
	if err != nil {
		return err
	}
	sqlstr, err = wrapSQL(baseDao.config.DBType, sqlstr)
	if err != nil {
		return err
	}

	tx, err := baseDao.dataSource.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	//流弊的...,把数组展开变成多个参数的形式
	tx.Exec(sqlstr, finder.Values...)

	tx.Commit()

	//fmt.Println(entity.GetTableName() + " save success")
	return nil
}

//保存Struct对象,必须是IEntityStruct类型
//bug(chunanuyong) 如果是自增主键,需要返回.需要sql驱动支持
func (baseDao *BaseDao) SaveStruct(entity IEntityStruct) error {
	if entity == nil {
		return errors.New("对象不能为空")
	}
	columns, values, err := columnAndValue(entity)
	if err != nil {
		return err
	}
	if len(columns) < 1 {
		return errors.New("没有tag信息,请检查struct中 column 的tag")
	}
	//SQL语句
	sqlstr, err := wrapSaveStructSQL(baseDao.config.DBType, entity, columns, values)
	if err != nil {
		return err
	}

	tx, err := baseDao.dataSource.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	//流弊的...,把数组展开变成多个参数的形式
	tx.Exec(sqlstr, values...)

	tx.Commit()

	return nil

}

//更新struct所有属性,必须是IEntityStruct类型
func (baseDao *BaseDao) UpdateStruct(entity IEntityStruct) error {
	return updateStructFunc(baseDao, entity, false)
}

//更新struct不为nil的属性,必须是IEntityStruct类型
func (baseDao *BaseDao) UpdateStructNotNil(entity IEntityStruct) error {
	return updateStructFunc(baseDao, entity, true)
}

// 根据主键删除一个对象.必须是IEntityStruct类型
func (baseDao *BaseDao) DeleteStruct(entity IEntityStruct) error {
	if entity == nil {
		return errors.New("对象不能为空")
	}
	pkName, err := entityPKFieldName(entity)
	if err != nil {
		return err
	}

	value, e := structFieldValue(entity, pkName)
	if e != nil {
		return e
	}
	//SQL语句
	sqlstr, err := wrapDeleteStructSQL(baseDao.config.DBType, entity)
	if err != nil {
		return err
	}

	tx, err := baseDao.dataSource.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	tx.Exec(sqlstr, value)

	tx.Commit()

	return nil

}

//保存IEntityMap对象.使用Map保存数据,需要在数据中封装好包括Id在内的所有数据.不适用于复杂情况
func (baseDao *BaseDao) SaveMap(entity IEntityMap) error {
	if entity == nil {
		return errors.New("对象不能为空")
	}
	//SQL语句
	sqlstr, values, err := wrapSaveMapSQL(baseDao.config.DBType, entity)
	if err != nil {
		return err
	}

	tx, err := baseDao.dataSource.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	//流弊的...,把数组展开变成多个参数的形式
	tx.Exec(sqlstr, values...)

	tx.Commit()

	return nil

}

//更新IEntityMap对象.使用Map修改数据,需要在数据中封装好包括Id在内的所有数据.不适用于复杂情况
func (baseDao *BaseDao) UpdateMap(entity IEntityMap) error {
	if entity == nil {
		return errors.New("对象不能为空")
	}
	//SQL语句
	sqlstr, values, err := wrapUpdateMapSQL(baseDao.config.DBType, entity)
	if err != nil {
		return err
	}
	//fmt.Println(sqlstr)

	tx, err := baseDao.dataSource.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	//流弊的...,把数组展开变成多个参数的形式
	tx.Exec(sqlstr, values...)

	tx.Commit()

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

//检查entity类型必须是*struct类型
func checkEntityKind(entity interface{}) error {
	if entity == nil {
		return errors.New("参数不能为空,必须是*struct类型")
	}
	typeOf := reflect.TypeOf(entity)
	if typeOf.Kind() != reflect.Ptr { //如果不是指针
		return errors.New("必须是*struct类型")
	}
	typeOf = typeOf.Elem()
	if typeOf.Kind() != reflect.Struct { //如果不是指针
		return errors.New("必须是*struct类型")
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
			//bug(chunanyong)这个地方还要添加其他类型的判断,参照ColumnValue.go文件

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
func updateStructFunc(baseDao *BaseDao, entity IEntityStruct, onlyupdatenotnull bool) error {
	if entity == nil {
		return errors.New("对象不能为空")
	}
	columns, values, err := columnAndValue(entity)
	if err != nil {
		return err
	}
	//SQL语句
	sqlstr, err := wrapUpdateStructSQL(baseDao.config.DBType, entity, columns, values, onlyupdatenotnull)
	if err != nil {
		return err
	}

	tx, err := baseDao.dataSource.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	//流弊的...,把数组展开变成多个参数的形式
	tx.Exec(sqlstr, values...)

	tx.Commit()

	//fmt.Println(entity.GetTableName() + " update success")
	return nil

}
