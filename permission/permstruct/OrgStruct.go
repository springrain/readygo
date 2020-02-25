package permstruct

import (
	"time"

	"readygo/orm"
)

// 部门
type OrgStruct struct {
	//引入默认的struct,隔离IEntityStruct的方法改动
	orm.EntityStruct

	// 编号
	Id string `column:"id"`

	// 名称
	Name string `column:"name"`

	// 代码
	Comcode string `column:"comcode"`

	// 上级部门ID
	Pid string `column:"pid"`

	// 0-99门店,100-199部门,200-299,分公司,300-399集团公司,900-999总平台
	OrgType int `column:"orgType"`

	// 排序,查询时倒叙排列
	Sortno int `column:"sortno"`

	// 备注
	Remark string `column:"remark"`

	// <no value>
	CreateTime time.Time `column:"createTime"`

	// <no value>
	CreateUserId string `column:"createUserId"`

	// <no value>
	UpdateTime time.Time `column:"updateTime"`

	// <no value>
	UpdateUserId string `column:"updateUserId"`

	// 是否有效(0否,1是)
	Active int `column:"active"`

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
func (entity *OrgStruct) GetTableName() string {
	return "t_org"
}

//获取数据库表的主键字段名称.因为要兼容Map,只能是数据库的字段名称.对应的struct 属性field
func (entity *OrgStruct) GetPKColumnName() string {
	return "id"
}
