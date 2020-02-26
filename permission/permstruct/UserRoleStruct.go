package permstruct

import (
	"readygo/orm"
)

//UserRoleStructTableName 表名常量,方便直接调用
const UserRoleStructTableName = "t_user_role"

// UserRoleStruct 用户角色中间表
type UserRoleStruct struct {
	//引入默认的struct,隔离IEntityStruct的方法改动
	orm.EntityStruct

	// 编号
	Id string `column:"id"`

	// 用户编号
	UserId string `column:"userId"`

	// 角色编号
	RoleId string `column:"roleId"`

	//------------------数据库字段结束,自定义字段写在下面---------------//

}

//获取表名称
func (entity *UserRoleStruct) GetTableName() string {
	return UserRoleStructTableName
}

//获取数据库表的主键字段名称.因为要兼容Map,只能是数据库的字段名称.对应的struct 属性field
func (entity *UserRoleStruct) GetPKColumnName() string {
	return "id"
}
