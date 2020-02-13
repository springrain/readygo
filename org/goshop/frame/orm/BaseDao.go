package orm

import (
	"context"
	"errors"
	"fmt"
	"goshop/org/goshop/frame/util"
	"reflect"
	"strings"
)

//数据库操作基类,隔离原生操作数据库API入口,所有数据库操作必须通过BaseDao进行.
type BaseDao struct {
	dataSource *dataSource
}

func NewBaseDao(config *DataSourceConfig) (*BaseDao, error) {
	dataSource, err := newDataSource(config)
	return &BaseDao{dataSource}, err
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
func (baseDao *BaseDao) Save(entity IBaseEntity) (bool, error) {

	// 获取实体类的反射
	valueOf := reflect.ValueOf(entity)

	//获取Kind,验证是否是指针,只能是*Struct结构
	if valueOf.Kind() != reflect.Ptr {
		return false, errors.New("只能保存*Struct类型")
	}
	//获取实体类的输出字段和私有字段
	exPortStructFields, _, err := util.StructFieldInfo(entity)
	if err != nil {
		return false, err
	}

	//实体类公开字段的长度
	fLen := len(exPortStructFields)
	//声明接收值的数组
	values := make([]interface{}, 0, fLen)
	//获取指针下struct的反射
	valueOf = valueOf.Elem()

	//SQL语句的构造器
	var sqlBuilder strings.Builder
	sqlBuilder.WriteString("INSERT INTO ")
	sqlBuilder.WriteString(entity.GetTableName())
	sqlBuilder.WriteString("(")

	//SQL语句中,VALUES(?,?,...)语句的构造器
	var valueSQLBuilder strings.Builder
	valueSQLBuilder.WriteString(" VALUES (")

	//遍历所有公共属性
	for i := 0; i < fLen; i++ {
		//获取字段类型的Kind
		fieldKind := exPortStructFields[i].Type.Kind()
		if fieldKind != reflect.String {
			continue
		}
		pname := exPortStructFields[i].Name
		//FieldByName方法返回的是reflect.Value类型,调用Interface()方法,返回原始类型的数据值
		value := valueOf.FieldByName(pname).Interface()
		//添加到记录值的数组
		values = append(values, value)
		//拼接SQL语句
		sqlBuilder.WriteString(pname)
		valueSQLBuilder.WriteString("?")
		//如果是最后一个字段
		if i+1 == fLen {
			sqlBuilder.WriteString(")")
			valueSQLBuilder.WriteString(")")
		} else {
			sqlBuilder.WriteString(",")
			valueSQLBuilder.WriteString(",")
		}

	}

	//SQL拼接VALUES语句
	sqlBuilder.WriteString(valueSQLBuilder.String())

	tx, err := baseDao.dataSource.BeginTx(context.Background(), nil)
	if err != nil {
		return false, err
	}
	defer tx.Rollback()

	// 测试,执行删除
	tx.Exec("DELETE FROM " + entity.GetTableName())
	//完成的SQL语句
	sqlstr := sqlBuilder.String()
	fmt.Println(sqlstr)

	//流弊的...,把数组展开变成多个参数的形式
	tx.Exec(sqlstr, values...)

	tx.Commit()

	fmt.Println(entity.GetTableName() + " save success")
	return true, nil

}
