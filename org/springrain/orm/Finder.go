package orm

import (
	"reflect"
	"strconv"
	"strings"
)

//查询数据库的载体,所有的sql语句都要通过Finder执行.
type Finder struct {
	//拼接SQL
	sqlBuilder strings.Builder
	//SQL的参数值
	Values []interface{}
	//默认false 不允许SQL注入的 ' 单引号
	InjectionSQL bool
	// 设置总条数查询的finder.Struct不能为nil,自己引用自己,go无法初始化Finder struct,使用可以为nil的指针,就可以了.
	//CountFinder Finder
	CountFinder *Finder
}

// 初始化一个Finder,生成一个空的Finder
func NewFinder() *Finder {
	finder := Finder{}
	finder.Values = make([]interface{}, 0)
	return &finder
}

//根据表名初始化查询的Finder
//NewSelectFinder("tableName") SELECT * FROM tableName
//NewSelectFinder("tableName", "id") SELECT id FROM tableName
func NewSelectFinder(tableName string, strs ...string) *Finder {
	finder := NewFinder()
	finder.sqlBuilder.WriteString("SELECT ")
	if len(strs) > 0 {
		for _, str := range strs {
			finder.sqlBuilder.WriteString(str)
		}
	} else {
		finder.sqlBuilder.WriteString("*")
	}
	finder.sqlBuilder.WriteString(" FROM ")
	finder.sqlBuilder.WriteString(tableName)
	return finder
}

//根据表名初始化更新的Finder,  UPDATE tableName SET
func NewUpdateFinder(tableName string) *Finder {
	finder := NewFinder()
	finder.sqlBuilder.WriteString("UPDATE ")
	finder.sqlBuilder.WriteString(tableName)
	finder.sqlBuilder.WriteString(" SET ")
	return finder
}

//根据表名初始化删除的Finder,  DELETE FROM tableName WHERE
func NewDeleteFinder(tableName string) *Finder {
	finder := NewFinder()
	finder.sqlBuilder.WriteString("DELETE FROM ")
	finder.sqlBuilder.WriteString(tableName)
	finder.sqlBuilder.WriteString(" WHERE ")
	return finder
}

//添加SQL和参数的值,第一个参数是语句,后面的参数[可选]是参数的值,顺序要正确.
//例如: finder.Append(" and id=? and name=? ",23123,"abc")
//只拼接SQL,例如: finder.Append(" and name=123 ")
func (finder *Finder) Append(s string, values ...interface{}) *Finder {

	if len(s) > 0 {
		finder.sqlBuilder.WriteString(s)
	}
	if values == nil || len(values) < 1 {
		return finder
	}
	//for _, v := range values {
	//	finder.Values = append(finder.Values, v)
	//}
	finder.Values = append(finder.Values, values...)
	return finder
}

//添加另一个Finder finder.AppendFinder(f)
func (finder *Finder) AppendFinder(f *Finder) *Finder {
	if f == nil {
		return nil
	}
	//添加f的SQL
	sqlstr := f.GetSQL()
	finder.sqlBuilder.WriteString(sqlstr)
	//添加f的值
	finder.Values = append(finder.Values, f.Values...)
	return finder
}

// 返回Finder封装的SQL语句
func (finder *Finder) GetSQL() string {
	sqlstr := finder.sqlBuilder.String()
	//包含单引号,属于非法字符串
	if !finder.InjectionSQL && (strings.Index(sqlstr, "'") >= 0) {
		return "SQL语句请不要直接拼接字符串参数!!!使用标准的占位符实现,例如  finder.Append(' and id=? and name=? ','123','abc')"
	}

	//处理sql语句中的in,实际就是把数组变量展开,例如 id in(?) ["1","2","3"] 语句变更为 id in (?,?,?) 参数也展开到参数数组里
	//这里认为 slice类型的参数就是in
	if finder.Values == nil || len(finder.Values) < 1 { //如果没有参数
		return sqlstr
	}

	//?问号切割的数组
	questions := strings.Split(sqlstr, "?")

	//语句中没有?问号
	if len(questions) < 1 {
		return sqlstr
	}

	//重新记录参数值
	newValues := make([]interface{}, 0)
	//新的sql
	var newSqlStr strings.Builder
	newSqlStr.WriteString(questions[0])
	for i, v := range finder.Values {

		//拼接?
		newSqlStr.WriteString("?")

		valueOf := reflect.ValueOf(v)
		typeOf := reflect.TypeOf(v)
		kind := valueOf.Kind()
		if kind == reflect.Ptr { //如果是指针
			valueOf = valueOf.Elem()
			typeOf = typeOf.Elem()
			kind = valueOf.Kind()
		}
		//获取数组长度
		sliceLen := valueOf.Len()
		//没有长度
		if sliceLen < 1 {
			return "语句:" + sqlstr + ",第" + strconv.Itoa(i+1) + "个参数,类型是Array或者Slice,值的长度为0,请检查sql参数有效性"
		}

		//如果不是数组或者slice
		if !(kind == reflect.Array || kind == reflect.Slice) {
			//记录新值
			newValues = append(newValues, v)
			//记录SQL
			newSqlStr.WriteString(questions[i+1])
			continue
		}
		//字节数组是特殊的情况
		if typeOf == reflect.TypeOf([]byte{}) {
			//记录新值
			newValues = append(newValues, v)
			//记录SQL
			newSqlStr.WriteString(questions[i+1])
			continue
		}
		for j := 0; j < sliceLen; j++ {
			//每多一个参数,对应",?" 两个符号.增加的问号长度总计是(sliceLen-1)*2.
			if j >= 1 {
				//记录SQL
				newSqlStr.WriteString(",?")
			}
			//记录新值
			sliceValue := valueOf.Index(j).Interface()
			newValues = append(newValues, sliceValue)
		}
		//记录SQL
		newSqlStr.WriteString(questions[i+1])
	}
	//重新赋值
	sqlstr = newSqlStr.String()
	finder.Values = newValues
	return sqlstr
}
