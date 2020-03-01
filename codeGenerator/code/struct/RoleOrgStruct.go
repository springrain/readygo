package permstruct
import (
	"time"

	"readygo/orm"
)

//表名常量,方便直接调用
const RoleOrgStructTableName = "角色部门中间表"

// 角色部门中间表
type RoleOrgStruct struct {
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
    CreateUserId sql.NullString `column:"createUserId"`
    
    // <no value>
    UpdateTime time.Time `column:"updateTime"`
    
    // <no value>
    UpdateUserId sql.NullString `column:"updateUserId"`
    
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
func (entity *RoleOrgStruct) GetTableName() string {
	return RoleOrgStructTableName
}

//获取数据库表的主键字段名称.因为要兼容Map,只能是数据库的字段名称.对应的struct 属性field
func (entity *RoleOrgStruct) GetPKColumnName() string {
	return "id"
}

