package peristruct

import (
	"time"
)

// 角色
type RoleStruct struct {

	// 角色ID
	Id string `column:"id"`

	// 角色名称
	Name string `column:"name"`

	// 权限编码
	Code string `column:"code"`

	// 上级角色ID,暂时不实现
	Pid string `column:"pid"`

	// 角色的部门是否私有,0否,1是,默认0.当角色私有时,菜单只使用此角色的部门权限,不再扩散到全局角色权限,用于设置特殊的菜单权限.公共权限时部门主管有所管理部门的数据全权限,无论角色是否分配. 私有部门权限时,严格按照配置的数据执行,部门主管可能没有部门权限.
	PrivateOrg int `column:"privateOrg"`

	// 0自己的数据,1所在部门,2所在部门及子部门数据,3.自定义部门数据.
	RoleOrgType int `column:"roleOrgType"`

	// 角色的归属部门,只有归属部门的主管和上级主管才可以管理角色,其他人员只能增加归属到角色的人员.不能选择部门或则其他操作,只能添加人员,不然存在提权风险,例如 员工角色下有1000人, 如果给 角色 设置了部门,那这1000人都起效了.
	OrgId string `column:"orgId"`

	// 角色是否共享,0否 1是,默认0,共享的角色可以被下级部门直接使用,但是下级只能添加人员,不能设置其他属性.共享的角色一般只设置roleOrgType,并不设定部门.
	ShareRole int `column:"shareRole"`

	// <no value>
	CreateTime time.Time `column:"createTime"`

	// <no value>
	CreateUserId string `column:"createUserId"`

	// <no value>
	UpdateTime time.Time `column:"updateTime"`

	// <no value>
	UpdateUserId string `column:"updateUserId"`

	// 排序,查询时倒叙排列
	Sortno int `column:"sortno"`

	// 备注
	Remark string `column:"remark"`

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
func (entity *RoleStruct) GetTableName() string {
	return "t_role"
}

//获取数据库表的主键字段名称.因为要兼容Map,只能是数据库的字段名称.对应的struct 属性field,使用cacheStructPKFieldNameMap缓存
func (entity *RoleStruct) GetPKColumnName() string {
	return "id"
}

//Oracle和pgsql没有自增,主键使用序列.优先级高于GetPKColumnName方法
func (entity *RoleStruct) GetPkSequence() string {
	return ""
}
