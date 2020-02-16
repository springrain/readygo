package orm

import (
	"context"
	"errors"
	"fmt"
	"goshop/org/goshop/frame/util"
	"reflect"
)

//允许的Type
var allowTypeMap = map[reflect.Kind]bool{
	reflect.Float32: true,
	reflect.Float64: true,
	reflect.Int:     true,
	reflect.String:  true,
}

const (
	tagColumnName = "column"
)

//缓存数据列
var cacheDBColumnMap = make(map[string][]reflect.StructField)

//缓存数据库字段和struct属性名对应的Map
var cacheColumn2FieldNameMap = make(map[string]map[string]string)

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

//根据Finder和封装为指定的entity类型,entity必须是*struct类型
func (baseDao *BaseDao) Query(finder *Finder, entity IEntityStruct) error {
	//检查Kind
	checke := checkEntityKind(entity)
	if checke != nil {
		return checke
	}
	//获取到Finder的语句
	sqlstr, err := wrapSQL(baseDao.config.DBType, finder.GetSQL())
	if err != nil {
		return err
	}
	//根据语句和参数查询
	rows, e := baseDao.dataSource.Query(sqlstr, finder.Values...)
	if e != nil {
		return e
	}
	//记录条数,本方法只能查询一个对象
	i := 0
	//数据库返回的列名
	columns, cne := rows.Columns()
	if cne != nil {
		return cne
	}
	//循环遍历结果集
	for rows.Next() {
		//只能查询出一条,如果查询出多条,只取第一条,然后抛错
		i++
		if i > 1 {
			i++
			break
		}
		i++
		//接收数据库返回的值
		values := make([]interface{}, len(columns))
		rows.Scan(values...)
		//根据列明和值,包装成Struct对象
		wrape := wrapStruct(columns, values, entity)
		if wrape != nil {
			return wrape
		}
		//包装struct对象
		wse := wrapStruct(columns, values, entity)
		if wse != nil {
			return wse
		}

	}

	//如果没有查询出数据
	if i == 0 {
		return errors.New("没有查询出数据")
	} else if i > 1 { //查询出多条数据
		return errors.New("查询出多条数据")
	}

	fmt.Println("Query:", entity)
	return nil
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

//保存Struct对象
func (baseDao *BaseDao) Save(entity IEntityStruct) error {
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
	fmt.Println(sqlstr)

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

//更新struct所有属性
func (baseDao *BaseDao) Update(entity IEntityStruct) error {
	return baseDao.updateStruct(entity, false)
}

//更新struct不为nil的属性
func (baseDao *BaseDao) UpdateNotNil(entity IEntityStruct) error {
	return baseDao.updateStruct(entity, true)
}

//更新对象
func (baseDao *BaseDao) updateStruct(entity IEntityStruct, onlyupdatenotnull bool) error {
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
	fmt.Println(sqlstr)

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

// 根据主键删除一个对象
func (baseDao *BaseDao) Delete(entity IEntityStruct) error {
	if entity == nil {
		return errors.New("对象不能为空")
	}
	pkName := entityPKFieldName(entity)
	value, err := util.StructFieldValue(entity, pkName)

	if err != nil {
		return err
	}
	//SQL语句
	sqlstr, err := wrapDeleteStructSQL(baseDao.config.DBType, entity)
	if err != nil {
		return err
	}
	fmt.Println(sqlstr)

	tx, err := baseDao.dataSource.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	tx.Exec(sqlstr, value)

	tx.Commit()

	return nil

}

//保存对象
func (baseDao *BaseDao) SaveMap(entity IEntityMap) error {
	if entity == nil {
		return errors.New("对象不能为空")
	}
	//SQL语句
	sqlstr, values, err := wrapSaveMapSQL(baseDao.config.DBType, entity)
	if err != nil {
		return err
	}
	fmt.Println(sqlstr)

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

//保存Map
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
func columnAndValue(entity IEntityStruct) ([]reflect.StructField, []interface{}, error) {

	// 获取实体类的反射
	valueOf := reflect.ValueOf(entity)

	//获取Kind,验证是否是指针,只能是*Struct结构
	if valueOf.Kind() != reflect.Ptr {
		return nil, nil, errors.New("只能是*Struct类型")
	}

	//先从本地缓存中查找
	entityName := reflect.TypeOf(entity).Elem().String()
	exPortStructFields := cacheDBColumnMap[entityName]
	if len(exPortStructFields) < 1 { //缓存不存在
		//获取实体类的输出字段和私有字段
		var err error
		exPortStructFields, _, err = util.StructFieldInfo(entity)
		if err != nil {
			return nil, nil, err
		}
	}

	//实体类公开字段的长度
	fLen := len(exPortStructFields)
	//接收列的数组
	columns := make([]reflect.StructField, 0, fLen)
	//接收值的数组
	values := make([]interface{}, 0, fLen)
	//获取指针下struct的反射
	valueOf = valueOf.Elem()

	//获取数据库列名和struct字段的对照缓存
	cacheColumn2Field := true
	column2FieldNameMap := cacheColumn2FieldNameMap[entityName]
	if column2FieldNameMap == nil || len(column2FieldNameMap) < 1 {
		cacheColumn2Field = false
		column2FieldNameMap = make(map[string]string)
		cacheColumn2FieldNameMap[entityName] = column2FieldNameMap
	}

	//遍历所有公共属性
	for i := 0; i < fLen; i++ {
		field := exPortStructFields[i]
		//获取字段类型的Kind
		fieldKind := field.Type.Kind()
		if !allowTypeMap[fieldKind] { //不允许的类型
			continue
		}

		// 只处理tag有column的字段
		tagValue := field.Tag.Get(tagColumnName)
		if len(tagValue) < 1 {
			continue
		}

		//如果没缓存列名和字段对应表
		if cacheColumn2Field {
			column2FieldNameMap[tagValue] = field.Name
		}

		columns = append(columns, field)
		//FieldByName方法返回的是reflect.Value类型,调用Interface()方法,返回原始类型的数据值
		value := valueOf.FieldByName(field.Name).Interface()
		//添加到记录值的数组
		values = append(values, value)

	}

	//缓存数据库的列
	cacheDBColumnMap[entityName] = columns

	return columns, values, nil

}

//获取实体类主键属性名称
func entityPKFieldName(entity IEntityStruct) string {
	//缓存的key,TypeOf和ValueOf的String()方法,返回值不一样
	cacheKey := reflect.TypeOf(entity).Elem().String()
	//列名和属性名的对照缓存
	column2FieldNameMap := cacheColumn2FieldNameMap[cacheKey]
	//如果缓存不存在,调用缓存逻辑
	if column2FieldNameMap == nil || len(column2FieldNameMap) < 1 {
		columnAndValue(entity)
	}

	column2FieldNameMap = cacheColumn2FieldNameMap[cacheKey]
	//获取主键的列名
	pkName := column2FieldNameMap[entity.GetPKColumnName()]

	return pkName

}

//检查entity类型必须是*struct类型
func checkEntityKind(entity IEntityStruct) error {
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

//根据数据库返回的sql.Rows,查询出列名和对应的值
func wrapStruct(columns []string, values []interface{}, entity IEntityStruct) error {
	checke := checkEntityKind(entity)
	if checke != nil {
		return checke
	}
	//缓存的key,TypeOf和ValueOf的String()方法,返回值不一样
	cacheKey := reflect.TypeOf(entity).Elem().String()
	//列名和属性名的对照缓存
	column2FieldNameMap := cacheColumn2FieldNameMap[cacheKey]
	//如果缓存不存在,调用缓存逻辑
	if column2FieldNameMap == nil || len(column2FieldNameMap) < 1 {
		columnAndValue(entity)
	}

	column2FieldNameMap = cacheColumn2FieldNameMap[cacheKey]

	//对象值的操作
	valueOf := reflect.ValueOf(entity).Elem()
	for i, column := range columns {
		fieldName := column2FieldNameMap[column]
		if len(fieldName) < 1 { //不存在列名,可以不接收
			continue
		}
		//给字段赋值
		valueOf.FieldByName(fieldName).Set(reflect.ValueOf(values[i]))

	}

	return nil

}
