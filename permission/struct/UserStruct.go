package permstruct

import (
	"time"
)

// 用户
type UserStruct struct {

	//
	Id string `column:"id"`

	// 姓名
	UserName string `column:"userName"`

	// 账号
	Account string `column:"account"`

	// 密码
	Password string `column:"password"`

	// 性别
	Sex string `column:"sex"`

	// 手机号码
	Mobile string `column:"mobile"`

	// 邮箱
	Email string `column:"email"`

	// 微信openId
	OpenId string `column:"openId"`

	// 微信UnionID
	UnionID string `column:"unionID"`

	// 头像地址
	Avatar string `column:"avatar"`

	// 0会员,1员工,2店长收银,9系统管理员
	UserType int `column:"userType"`

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
func (entity *UserStruct) GetTableName() string {
	return "t_user"
}

//获取数据库表的主键字段名称.因为要兼容Map,只能是数据库的字段名称.对应的struct 属性field
func (entity *UserStruct) GetPKColumnName() string {
	return "id"
}

//Oracle和pgsql没有自增,主键使用序列.优先级高于GetPKColumnName方法
func (entity *UserStruct) GetPkSequence() string {
	return ""
}
