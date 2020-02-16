package util

import (
	"bytes"
	"encoding/gob"
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

// 用于缓存反射的信息,sync.Map内部处理了并发锁.全局变量不能使用:=简写声明
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
		return nil, nil, errors.New("必须是Struct或者*Struct类型")
	}

	if kind == reflect.Ptr {
		//获取指针下的Struct类型
		typeOf = typeOf.Elem()
		if typeOf.Kind() != reflect.Struct {
			return nil, nil, errors.New("必须是Struct或者*Struct类型")
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

	// 声明所有字段的载体
	var allFieldMap sync.Map
	anonymous := make([]reflect.StructField, 0)

	//遍历所有字段,记录匿名属性
	for i := 0; i < fieldNum; i++ {
		field := typeOf.Field(i)
		if _, ok := allFieldMap.Load(field.Name); !ok {
			allFieldMap.Store(field.Name, field)
		}
		if field.Anonymous { //如果是匿名的
			anonymous = append(anonymous, field)
		}
	}
	//调用匿名struct的递归方法
	recursiveAnonymousStruct(&allFieldMap, anonymous)

	//遍历sync.Map,要求输入一个func作为参数
	//这个函数的入参、出参的类型都已经固定，不能修改
	//可以在函数体内编写自己的代码,调用map中的k,v
	f := func(k, v interface{}) bool {
		// fmt.Println(k, ":", v)
		field := v.(reflect.StructField)
		if ast.IsExported(field.Name) { //如果是可以输出的
			exPortStructFields = append(exPortStructFields, field)
		} else {
			privateStructFields = append(privateStructFields, field)
		}

		return true
	}
	allFieldMap.Range(f)

	//加入缓存
	cacheStructFieldMap.Store(exPortCacheKey, exPortStructFields)
	cacheStructFieldMap.Store(privateCacheKey, privateStructFields)

	return exPortStructFields, privateStructFields, nil
}

//递归调用struct的匿名属性,就近覆盖属性.
func recursiveAnonymousStruct(allFieldMap *sync.Map, anonymous []reflect.StructField) {

	for i := 0; i < len(anonymous); i++ {
		field := anonymous[i]
		typeOf := field.Type

		if typeOf.Kind() == reflect.Ptr {
			//获取指针下的Struct类型
			typeOf = typeOf.Elem()
		}

		//只处理Struct类型
		if typeOf.Kind() != reflect.Struct {
			continue
		}

		//获取字段长度
		fieldNum := typeOf.NumField()
		//如果没有字段
		if fieldNum < 1 {
			continue
		}

		// 匿名struct里自身又有匿名struct
		anonymousField := make([]reflect.StructField, 0)

		//遍历所有字段
		for i := 0; i < fieldNum; i++ {
			field := typeOf.Field(i)
			if _, ok := allFieldMap.Load(field.Name); ok { //如果存在属性名
				continue
			} else { //不存在属性名,加入到allFieldMap
				allFieldMap.Store(field.Name, field)
			}

			if field.Anonymous { //匿名struct里自身又有匿名struct
				anonymousField = append(anonymousField, field)
			}
		}

		//递归调用匿名struct
		recursiveAnonymousStruct(allFieldMap, anonymousField)

	}

}

//获取指定字段的值
func StructFieldValue(s interface{}, fieldName string) (interface{}, error) {

	if s == nil {
		return nil, errors.New("数据为空")
	}
	//entity的s类型
	valueOf := reflect.ValueOf(s)

	kind := valueOf.Kind()
	if !(kind == reflect.Ptr || kind == reflect.Struct) {
		return nil, errors.New("必须是Struct或者*Struct类型")
	}

	if kind == reflect.Ptr {
		//获取指针下的Struct类型
		valueOf = valueOf.Elem()
		if valueOf.Kind() != reflect.Struct {
			return nil, errors.New("必须是Struct或者*Struct类型")
		}
	}

	//FieldByName方法返回的是reflect.Value类型,调用Interface()方法,返回原始类型的数据值
	value := valueOf.FieldByName(fieldName).Interface()

	return value, nil

}

//深度拷贝对象.golang没有构造函数,反射复制对象时,对象中struct类型的属性无法初始化,指针属性也会收到影响.使用深度对象拷贝
func DeepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}
