package util

import (
	"errors"
	"go/ast"
	"reflect"
	"sync"
)

//缓存map的key前缀
const (
	exPortPrefix  = "_exPortStructFields_"
	privatePrefix = "_privateStructFields_"
)

// 用于缓存反射的信息.全局变量不能使用:=简写声明
var cacheStructFieldMap sync.Map

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

	//缓存的key
	exPortCacheKey := exPortPrefix + typeOf.String()
	privateCacheKey := privatePrefix + typeOf.String()
	//缓存的值
	cacheExPortStructFields, exportOk := cacheStructFieldMap.Load(exPortCacheKey)
	cachePrivateStructFields, privateOk := cacheStructFieldMap.Load(privateCacheKey)

	//如果存在值,返回值为true
	if exportOk && privateOk {
		//把 interface{} 类型 转为 []reflect.StructField 类型
		return cacheExPortStructFields.([]reflect.StructField), cachePrivateStructFields.([]reflect.StructField), nil
	}
	//获取字段长度
	fieldNum := typeOf.NumField()
	//如果没有字段
	if fieldNum < 1 {
		return nil, nil, errors.New("entity没有属性")
	}
	//声明承接数组
	exPortStructFields := make([]reflect.StructField, 0)
	privateStructFields := make([]reflect.StructField, 0)
	//遍历所有字段
	for i := 0; i < fieldNum; i++ {
		field := typeOf.Field(i)
		if ast.IsExported(field.Name) { //如果是可以输出的
			exPortStructFields = append(exPortStructFields, field)
		} else {
			privateStructFields = append(privateStructFields, field)
		}
	}

	//加入缓存
	cacheStructFieldMap.Store(exPortCacheKey, exPortStructFields)
	cacheStructFieldMap.Store(privateCacheKey, privateStructFields)

	return exPortStructFields, privateStructFields, nil
}
