package permstruct

import (
	"readygo/orm"
)

//RoleMenuStructTableName 表名常量,方便直接调用
const RoleMenuStructTableName = "t_role_menu"

// RoleMenuStruct 角色菜单中间表
type RoleMenuStruct struct {
	//引入默认的struct,隔离IEntityStruct的方法改动
	orm.EntityStruct

	//Id 编号
	Id string `column:"id"`

	//RoleId 角色编号
	RoleId string `column:"roleId"`

	//MenuId 菜单编号
	MenuId string `column:"menuId"`

	//------------------数据库字段结束,自定义字段写在下面---------------//

	//Checked 是否选中,未选中就是删除.用于前台操作
	Checked bool
}

//GetTableName 获取表名称
func (entity *RoleMenuStruct) GetTableName() string {
	return RoleMenuStructTableName
}

//GetPKColumnName 获取数据库表的主键字段名称.因为要兼容Map,只能是数据库的字段名称.
func (entity *RoleMenuStruct) GetPKColumnName() string {
	return "id"
}
