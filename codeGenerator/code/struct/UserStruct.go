package permstruct
import (
	"time"

	"readygo/orm"
)

//表名常量,方便直接调用
const UserStructTableName = "用户"

// 用户
type UserStruct struct {
	//引入默认的struct,隔离IEntityStruct的方法改动
	orm.EntityStruct
	
    //  
    Id string `column:"id"`
    
    // 姓名
    UserName sql.NullString `column:"userName"`
    
    // 账号
    Account string `column:"account"`
    
    // 密码
    Password string `column:"password"`
    
    // 性别
    Sex sql.NullString `column:"sex"`
    
    // 手机号码
    Mobile sql.NullString `column:"mobile"`
    
    // 邮箱
    Email sql.NullString `column:"email"`
    
    // 微信openId
    OpenId sql.NullString `column:"openId"`
    
    // 微信UnionID
    UnionID sql.NullString `column:"unionID"`
    
    // 头像地址
    Avatar sql.NullString `column:"avatar"`
    
    // 0会员,1员工,2店长收银,9系统管理员
    UserType int `column:"userType"`
    
    // <no value>
    CreateTime time.Time `column:"createTime"`
    
    // <no value>
    CreateUserId sql.NullString `column:"createUserId"`
    
    // <no value>
    UpdateTime time.Time `column:"updateTime"`
    
    // <no value>
    UpdateUserId sql.NullString `column:"updateUserId"`
    
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
func (entity *UserStruct) GetTableName() string {
	return UserStructTableName
}

//获取数据库表的主键字段名称.因为要兼容Map,只能是数据库的字段名称.对应的struct 属性field
func (entity *UserStruct) GetPKColumnName() string {
	return "id"
}

