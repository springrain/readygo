package permstruct

import (
	"time"

	"readygo/orm"
)

//DicDataStructTableName 表名常量,方便直接调用
const DicDataStructTableName = "t_dic_data"

// DicDataStruct 公共字典
type DicDataStruct struct {
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

	// 描述
	Remark string `column:"remark"`

	// 类型
	Typekey string `column:"typekey"`

	// <no value>
	CreateTime time.Time `column:"createTime"`

	// <no value>
	CreateUserId string `column:"createUserId"`

	// <no value>
	UpdateTime time.Time `column:"updateTime"`

	// <no value>
	UpdateUserId string `column:"updateUserId"`

	// 排序
	Sortno int `column:"sortno"`

	// 是否有效(0否,1是)
	Active int `column:"active"`

	//------------------数据库字段结束,自定义字段写在下面---------------//

}

//获取表名称
func (entity *DicDataStruct) GetTableName() string {
	return DicDataStructTableName
}

//获取数据库表的主键字段名称.因为要兼容Map,只能是数据库的字段名称.对应的struct 属性field
func (entity *DicDataStruct) GetPKColumnName() string {
	return "id"
}
