package permstruct

import (
	"readygo/orm"
)

// 公共字典
type DicDataStructStruct struct {
	//引入默认的struct,隔离IEntityStruct的方法改动
	orm.EntityStruct

	// <no value>
	Id string `column:"id"`

	// 名称
	Name string `column:"name"`

	// 编码
	Code string `column:"code"`

	// 值
	Val string `column:"val"`

	// 父ID
	Pid string `column:"pid"`

	// 排序
	Sortno int `column:"sortno"`

	// 描述
	Remark string `column:"remark"`

	// 是否有效(0否,1是)
	Active int `column:"active"`

	// 类型
	Typekey string `column:"typekey"`

	// <no value>
	Bak1 string `column:"bak1"`

	// <no value>
	Bak2 string `column:"bak2"`

	// <no value>
	Bak3 string `column:"bak3"`

	// <no value>
	Bak4 string `column:"bak4"`

	// <no value>
	Bak5 string `column:"bak5"`

	//------------------数据库字段结束,自定义字段写在下面---------------//

}

//获取表名称
func (entity *DicDataStructStruct) GetTableName() string {
	return "t_dic_data"
}

//获取数据库表的主键字段名称.因为要兼容Map,只能是数据库的字段名称.对应的struct 属性field
func (entity *DicDataStructStruct) GetPKColumnName() string {
	return "id"
}
