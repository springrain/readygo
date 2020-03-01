package permstruct
import (
	"time"

	"readygo/orm"
)

//表名常量,方便直接调用
const RoleStructTableName = "角色"

// 角色
type RoleStruct struct {
	//引入默认的struct,隔离IEntityStruct的方法改动
	orm.EntityStruct
	
    // 角色ID
    Id string `column:"id"`
    
    // 角色名称
    Name string `column:"name"`
    
    // 权限编码
    Code sql.NullString `column:"code"`
    
    // 上级角色ID,暂时不实现
    Pid sql.NullString `column:"pid"`
    
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
    CreateUserId sql.NullString `column:"createUserId"`
    
    // <no value>
    UpdateTime time.Time `column:"updateTime"`
    
    // <no value>
    UpdateUserId sql.NullString `column:"updateUserId"`
    
    // 排序,查询时倒叙排列
    Sortno int `column:"sortno"`
    
    // 备注
    Remark sql.NullString `column:"remark"`
    
    // 是否有效(0否,1是)
    Active int `column:"active"`
    
    // <no value>
    Bak1 sql.NullString `column:"bak1"`
    
    // <no value>
    Bak2 sql.NullString `column:"bak2"`
    
    // <no value>
    Bak3 sql.NullString `column:"bak3"`
    
    // <no value>
    Bak4 sql.NullString `column:"bak4"`
    
    // <no value>
    Bak5 sql.NullString `column:"bak5"`
    
	//------------------数据库字段结束,自定义字段写在下面---------------//


}


//获取表名称
func (entity *RoleStruct) GetTableName() string {
	return RoleStructTableName
}

//获取数据库表的主键字段名称.因为要兼容Map,只能是数据库的字段名称.对应的struct 属性field
func (entity *RoleStruct) GetPKColumnName() string {
	return "id"
}

