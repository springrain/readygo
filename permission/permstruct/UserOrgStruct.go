package permstruct

import (
	"readygo/orm"
)

//UserOrgStructTableName 表名常量,方便直接调用
const UserOrgStructTableName = "t_user_org"

// UserOrgStruct 用户部门中间表
type UserOrgStruct struct {
	//引入默认的struct,隔离IEntityStruct的方法改动
	orm.EntityStruct

	// 编号
	Id string `column:"id"`

	// 用户编号
	UserId string `column:"userId"`

	// 机构编号
	OrgId string `column:"orgId"`

	// 0会员,1员工,2主管
	ManagerType int `column:"managerType"`

	//------------------数据库字段结束,自定义字段写在下面---------------//

}

//获取表名称
func (entity *UserOrgStruct) GetTableName() string {
	return UserOrgStructTableName
}

//获取数据库表的主键字段名称.因为要兼容Map,只能是数据库的字段名称.对应的struct 属性field
func (entity *UserOrgStruct) GetPKColumnName() string {
	return "id"
}
