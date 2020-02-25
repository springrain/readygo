package permstruct

import (
	"time"

	"readygo/orm"
)

// 角色部门中间表
type RoleOrgStructStruct struct {
	//引入默认的struct,隔离IEntityStruct的方法改动
	orm.EntityStruct

	// 编号
	Id string `column:"id"`

	// 部门编号
	OrgId string `column:"orgId"`

	// 角色编号
	RoleId string `column:"roleId"`

	// 0不包含子部门,1包含.用于表示角色和部门的权限关系
	Children int `column:"children"`

	// <no value>
	CreateTime time.Time `column:"createTime"`

	// <no value>
	CreateUserId string `column:"createUserId"`

	// <no value>
	UpdateTime time.Time `column:"updateTime"`

	// <no value>
	UpdateUserId string `column:"updateUserId"`

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
func (entity *RoleOrgStructStruct) GetTableName() string {
	return "t_role_org"
}

//获取数据库表的主键字段名称.因为要兼容Map,只能是数据库的字段名称.对应的struct 属性field
func (entity *RoleOrgStructStruct) GetPKColumnName() string {
	return "id"
}
