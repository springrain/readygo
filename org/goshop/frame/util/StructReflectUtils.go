package util

import (
	"errors"
	"go/ast"
	"reflect"
)

//获取StructField的信息.只对struct或者*struct判断,如果是指针,返回指针下实际的struct类型.
//第一个返回值是可以输出的字段(首字母大写),第二个是不能输出的字段(首字母小写)
func StructFieldInfo(s interface{}) ([]reflect.StructField, []reflect.StructField, error) {

	if s == nil {
		return nil, nil, errors.New("数据为空")
	}
	//entity的s类型
	typeOf := reflect.TypeOf(s)

	kind := typeOf.Kind()

	if !(kind == reflect.Ptr || kind == reflect.Struct) {
		return nil, nil, errors.New("entity必须是Struct或者*Struct类型")
	}

	if kind == reflect.Ptr {
		//获取指针下的Struct类型
		typeOf = typeOf.Elem()
		if typeOf.Kind() != reflect.Struct {
			return nil, nil, errors.New("entity必须是Struct或者*Struct类型")
		}
	}

	fieldNum := typeOf.NumField()

	exPortStructFields := make([]reflect.StructField, 0)
	privateStructFields := make([]reflect.StructField, 0)
	for i := 0; i < fieldNum; i++ {
		field := typeOf.Field(i)
		if ast.IsExported(field.Name) { //如果是可以输出的
			exPortStructFields = append(exPortStructFields, field)
		} else {
			privateStructFields = append(privateStructFields, field)
		}
	}
	return exPortStructFields, privateStructFields, nil
}
