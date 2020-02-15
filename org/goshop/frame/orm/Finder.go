package orm

import (
	"errors"
	"strings"
)

//查询数据库的载体,所有的sql语句都要通过Finder执行.
type Finder struct {
	//拼接SQL
	sqlBuilder strings.Builder
	//SQL的参数值
	Values []interface{}
	//默认不允许SQL注入的 ' 单引号
	InjectionSQL bool
	// 设置总条数查询的finder.Struct不能为nil,自己引用自己,go无法初始化Finder struct,使用可以为nil的指针,就可以了.
	//CountFinder Finder
	CountFinder *Finder
}

// 初始化一个Finder,生成一个空的Finder
func NewFinder() *Finder {
	finder := Finder{}
	finder.Values = []interface{}{}
	return &finder
}

//根据表名初始化查询的Finder, SELECT * FROM tableName
func NewSelectFinder(tableName string) *Finder {
	finder := Finder{}
	finder.Values = []interface{}{}
	finder.sqlBuilder.WriteString("SELECT * FROM ")
	finder.sqlBuilder.WriteString(tableName)
	return &finder
}

//根据表名初始化更新的Finder,  UPDATE tableName SET
func NewUpdateFinder(tableName string) *Finder {
	finder := Finder{}
	finder.Values = []interface{}{}
	finder.sqlBuilder.WriteString("UPDATE ")
	finder.sqlBuilder.WriteString(tableName)
	finder.sqlBuilder.WriteString(" SET ")
	return &finder
}

//根据表名初始化删除的Finder,  DELETE FROM tableName WHERE
func NewDeleteFinder(tableName string) *Finder {
	finder := Finder{}
	finder.Values = []interface{}{}
	finder.sqlBuilder.WriteString("DELETE FROM ")
	finder.sqlBuilder.WriteString(tableName)
	finder.sqlBuilder.WriteString(" WHERE ")
	return &finder
}

//添加SQL和参数的值,第一个参数是语句,后面的参数[可选]是参数的值,顺序要正确.
//例如: finder.Append(" and id=? and name=? ",23123,"abc")
//只拼接SQL,例如: finder.Append(" and name=123 ")
func (finder *Finder) Append(s string, v ...interface{}) *Finder {

	if len(s) > 0 {
		finder.sqlBuilder.WriteString(s)
	}
	if v != nil {
		finder.Values = append(finder.Values, v)
	}
	return finder
}

//添加另一个Finder finder.AppendFinder(f)
func (finder *Finder) AppendFinder(f *Finder) *Finder {
	if f == nil {
		return nil
	}
	//添加f的SQL
	sqlstr, err := f.GetSQL()
	if err != nil {
		return nil
	}
	finder.sqlBuilder.WriteString(sqlstr)
	//添加f的值
	finder.Values = append(finder.Values, f.Values...)
	return finder
}

// 返回Finder封装的SQL语句
func (finder *Finder) GetSQL() (string, error) {
	sqlstr := finder.sqlBuilder.String()
	if !finder.InjectionSQL && (strings.Index(sqlstr, "'") >= 0) { //包含单引号,属于非法字符串
		return "", errors.New("SQL语句请不要直接拼接字符串参数!!!使用标准的占位符实现,例如  finder.Append(' and id=? and name=? ','123','abc')")
	}

	return sqlstr, nil
}
