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

	valueOf := reflect.ValueOf(entity)

	if valueOf.Kind() != reflect.Ptr {
		return false, errors.New("只能保存*Struct类型")
	}

	exPortStructFields, _, err := util.StructFieldInfo(entity)
	if err != nil {
		return false, err
	}

	fLen := len(exPortStructFields)

	values := make([]interface{}, 0, fLen)
	valueOf = valueOf.Elem()

	var sqlBuilder strings.Builder
	sqlBuilder.WriteString("INSERT INTO ")
	sqlBuilder.WriteString(entity.GetTableName())
	sqlBuilder.WriteString("(")

	var valueSQLBuilder strings.Builder
	valueSQLBuilder.WriteString(" VALUES (")

	for i := 0; i < fLen; i++ {
		fieldKind := exPortStructFields[i].Type.Kind()

		if fieldKind != reflect.String {
			continue
		}
		pname := exPortStructFields[i].Name
		//FieldByName方法返回的是reflect.Value类型,调用Interface()方法,返回原始类型的数据值
		value := valueOf.FieldByName(pname).Interface()
		fmt.Println("value:", value)
		values = append(values, value)
		//values = append(values, "id")
		sqlBuilder.WriteString(pname)
		valueSQLBuilder.WriteString("?")
		if i+1 == fLen {
			sqlBuilder.WriteString(")")
			valueSQLBuilder.WriteString(")")
		} else {
			sqlBuilder.WriteString(",")
			valueSQLBuilder.WriteString(",")
		}

	}

	sqlBuilder.WriteString(valueSQLBuilder.String())

	tx, err := baseDao.dataSource.BeginTx(context.Background(), nil)
	if err != nil {
		return false, err
	}
	defer tx.Rollback()

	tx.Exec("DELETE FROM " + entity.GetTableName())
	sqlstr := sqlBuilder.String()
	fmt.Println(sqlstr)

	//流弊的...,把数组重新变成多个参数模式
	tx.Exec(sqlstr, values...)
	//strconv.Atoi("a")
	tx.Commit()

	fmt.Println(entity.GetTableName() + " save success")
	return true, nil

}
