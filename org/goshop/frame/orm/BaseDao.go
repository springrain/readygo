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

//数据库操作基类,隔离原生操作数据库API入口,所有数据库操作必须通过BaseDao进行.
type BaseDao struct {
	config     *DataSourceConfig
	dataSource *dataSource
}

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

//保存对象
func (baseDao *BaseDao) Save(entity IBaseEntity) error {

	columns, values, err := columnAndValue(entity)
	if err != nil {
		return err
	}
	//SQL语句
	sqlstr, err := wrapsavesql(baseDao.config.DBType, entity, columns, values)
	if err != nil {
		return err
	}
	fmt.Println(sqlstr)

	tx, err := baseDao.dataSource.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 测试,执行删除
	tx.Exec("DELETE FROM " + entity.GetTableName())

	//流弊的...,把数组展开变成多个参数的形式
	tx.Exec(sqlstr, values...)

	tx.Commit()

	fmt.Println(entity.GetTableName() + " save success")
	return nil

}

//根据保存的对象,返回插入的语句,需要插入的字段,字段的值.
func columnAndValue(entity IBaseEntity) ([]reflect.StructField, []interface{}, error) {

	// 获取实体类的反射
	valueOf := reflect.ValueOf(entity)

	//获取Kind,验证是否是指针,只能是*Struct结构
	if valueOf.Kind() != reflect.Ptr {
		return nil, nil, errors.New("只能保存*Struct类型")
	}

	//获取实体类的输出字段和私有字段
	exPortStructFields, _, err := util.StructFieldInfo(entity)
	if err != nil {
		return nil, nil, err
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
		//获取字段类型的Kind
		fieldKind := exPortStructFields[i].Type.Kind()
		if !allowTypeMap[fieldKind] { //不允许的类型
			continue
		}
		columns = append(columns, exPortStructFields[i])
		pname := exPortStructFields[i].Name
		//FieldByName方法返回的是reflect.Value类型,调用Interface()方法,返回原始类型的数据值
		value := valueOf.FieldByName(pname).Interface()
		//添加到记录值的数组
		values = append(values, value)

	}

	return columns, values, nil

}
