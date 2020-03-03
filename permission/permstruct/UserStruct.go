package permstruct

import (
	"readygo/zorm"
	"time"
)

//UserStructTableName 表名常量,方便直接调用
const UserStructTableName = "t_user"

// UserStruct 用户
type UserStruct struct {
	//引入默认的struct,隔离IEntityStruct的方法改动
	zorm.EntityStruct

	//Id
	Id string `column:"id"`

	//UserName 姓名
	UserName string `column:"userName"`

	//Account 账号
	Account string `column:"account"`

	//Password 密码
	Password string `column:"password"`

	//Sex 性别
	Sex string `column:"sex"`

	//Mobile 手机号码
	Mobile string `column:"mobile"`

	//Email 邮箱
	Email string `column:"email"`

	//OpenId 微信openId
	OpenId string `column:"openId"`

	//UnionID 微信UnionID
	UnionID string `column:"unionID"`

	//Avatar 头像地址
	Avatar string `column:"avatar"`

	//UserType 0会员,1员工,2店长收银,9系统管理员
	UserType int `column:"userType"`

	//CreateTime <no value>
	CreateTime time.Time `column:"createTime"`

	//CreateUserId <no value>
	CreateUserId string `column:"createUserId"`

	//UpdateTime <no value>
	UpdateTime time.Time `column:"updateTime"`

	//UpdateUserId <no value>
	UpdateUserId string `column:"updateUserId"`

	//Active 是否有效(0否,1是)
	Active int `column:"active"`

	//------------------数据库字段结束,自定义字段写在下面---------------//

	//Roles 用户的角色
	Roles []RoleStruct
}

//GetTableName 获取表名称
func (entity *UserStruct) GetTableName() string {
	return UserStructTableName
}

//GetPKColumnName 获取数据库表的主键字段名称.因为要兼容Map,只能是数据库的字段名称.
func (entity *UserStruct) GetPKColumnName() string {
	return "id"
}
