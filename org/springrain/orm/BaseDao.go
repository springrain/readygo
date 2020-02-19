package orm

import (
	"context"
	"errors"
	"reflect"
)

//允许的Type
//bug(chunanyong) 1.需要完善支持的数据类型和赋值接口,例如sql.NullString.
var allowTypeMap = map[reflect.Kind]bool{
	reflect.Float32: true,
	reflect.Float64: true,
	reflect.Int:     true,
	reflect.String:  true,
}

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

	//获取map对象
	resultMap, err := baseDao.QueryMap(finder)
	if err != nil {
		return err
	}
	e := columnValueMap2Struct(resultMap, entity)

	if e != nil {
		return e
	}

	return nil
}

//根据Finder和封装为指定的entity类型,entity必须是*[]struct类型,已经初始化好的数组,此方法只Append元素,这样调用方就不需要强制类型转换了.
func (baseDao *BaseDao) QueryStructList(finder *Finder, rowsSlicePtr interface{}, page *Page) error {
	mapList, err := baseDao.QueryMapList(finder, page)
	if err != nil {
		return err
	}

	//获取数组内元素的类型
	sliceValue := reflect.Indirect(reflect.ValueOf(rowsSlicePtr))
	sliceElementType := sliceValue.Type().Elem()

	dbMap, err := getDBColumnFieldMap(sliceElementType)

	if err != nil {
		return err
	}

	//	var a []structType = structList.([]structType)
	//valueType := reflect.ValueOf(structList).Elem()
	for _, resultMap := range mapList {
		//deepCopy(a, entity)
		//反射初始化一个元素
		//new 出来的为什么是个指针啊????
		pv := reflect.New(sliceElementType).Elem()

		//bug(chunanyong)需要重新梳理字段缓存
		for column, columnValue := range resultMap {
			field := dbMap[column]
			pv.FieldByName(field.Name).Set(reflect.ValueOf(columnValue.String()))
		}

		sliceValue.Set(reflect.Append(sliceValue, pv))
	}

	return nil

}

//根据Finder查询,封装Map.获取具体的值,需要自己根据类型调用ColumnValue的转化方法,例如ColumnValue.String()
//golang的sql驱动不支持获取到数据字段的metadata......垃圾.....
//bug(chunanyong)需要测试一下 in 数组, like ,还有查询一个基础类型(例如 string)的功能
func (baseDao *BaseDao) QueryMap(finder *Finder) (map[string]ColumnValue, error) {
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
//golang的sql驱动不支持获取到数据字段的metadata......垃圾.....
func (baseDao *BaseDao) QueryMapList(finder *Finder, page *Page) ([]map[string]ColumnValue, error) {

	var sqlstr string
	var err error

	//获取到没有page的sql的语句
	if page == nil {
		sqlstr, err = wrapSQL(baseDao.config.DBType, finder.GetSQL())
	} else {
		sqlstr, err = wrapPageSQL(baseDao.config.DBType, finder.GetSQL(), page)
	}

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
	resultMapList := make([]map[string]ColumnValue, 0)
	//循环遍历结果集
	for rows.Next() {
		//接收数据库返回的值,返回的字段值都是[]byte直接数组,需要使用指针接收.比较恶心......

		values := make([]ColumnValue, len(columns))
		//使用指针类型接收字段值,需要使用interface{}包装一下
		scans := make([]interface{}, len(columns))
		//包装[]byte的指针地址包装
		for j := range values {
			scans[j] = &values[j]
		}
		//接收数据库返回值,之后values就有值了
		err = rows.Scan(scans...)
		if err != nil {
			return nil, err
		}
		result, e := wrapMap(columns, values)
		if e != nil {
			return nil, e
		}
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

	sqlstr := finder.GetSQL()
	var err error
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

	//fmt.Println(entity.GetTableName() + " save success")
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

	//fmt.Println(entity.GetTableName() + " save success")
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
		fieldKind := field.Type.Kind()
		if !allowTypeMap[fieldKind] { //不允许的类型
			continue
		}

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

//根据数据库返回的sql.Rows,查询出列名和对应的值.
func columnValueMap2Struct(resultMap map[string]ColumnValue, entity interface{}) error {

	checkerr := checkEntityKind(entity)
	if checkerr != nil {
		return checkerr
	}

	typeOf := reflect.TypeOf(entity).Elem()
	dbMap, err := getDBColumnFieldMap(typeOf)
	if err != nil {
		return err
	}

	for column, columnValue := range resultMap {
		field := dbMap[column]
		fieldName := field.Name
		if len(fieldName) < 1 {
			continue
		}
		//反射获取字段的值对象
		fieldValue := reflect.ValueOf(entity).Elem().FieldByName(fieldName)
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

//根据sql查询结果,返回map
func wrapMap(columns []string, values []ColumnValue) (map[string]ColumnValue, error) {
	columnValueMap := make(map[string]ColumnValue)
	for i, column := range columns {
		columnValueMap[column] = values[i]
	}
	return columnValueMap, nil
}

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
