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

//缓存主键名称
var cacheStructPKFieldNameMap = make(map[string]string)

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

func (baseDao *BaseDao) Query(sql string) {
	rows, err := baseDao.dataSource.Query(sql)
	if err != nil {
		return
	}
	for rows.Next() {
		var id string
		var account string
		rows.Scan(&id, &account)
		fmt.Println(id + ":" + account)
	}
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

//保存对象
func (baseDao *BaseDao) Update(entity IEntityStruct, onlyupdatenotnull bool) error {
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

		//缓存主键的属性名
		if tagValue == entity.GetPKColumnName() {
			cacheStructPKFieldNameMap[entityName] = field.Name
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
	cacheKey := reflect.TypeOf(entity).Elem().String()

	fieldName := cacheStructPKFieldNameMap[cacheKey]
	if len(fieldName) > 0 {
		return fieldName
	}
	columnAndValue(entity)
	pkName := cacheStructPKFieldNameMap[cacheKey]

	return pkName

}
