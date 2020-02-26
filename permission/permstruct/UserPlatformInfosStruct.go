package permstruct

import (
	"time"

	"readygo/orm"
)

//UserPlatformInfosStructTableName 表名常量,方便直接调用
const UserPlatformInfosStructTableName = "t_user_platform_infos"

// UserPlatformInfosStruct 用户平台信息表
type UserPlatformInfosStruct struct {
	//引入默认的struct,隔离IEntityStruct的方法改动
	orm.EntityStruct

	// 主键id
	Id string `column:"id"`

	// t_user表中ID
	UserId string `column:"userId"`

	// 公众号openId,企业号userId,小程序openId,APP推送deviceToken
	OpenId string `column:"openId"`

	// 设备/应用类型：1公众号2小程序3企业号4APP IOS消息推送5APP安卓消息推送6web
	DeviceType int `column:"deviceType"`

	// 所属组织机构ID
	OrgId string `column:"orgId"`

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

	//------------------数据库字段结束,自定义字段写在下面---------------//

}

//获取表名称
func (entity *UserPlatformInfosStruct) GetTableName() string {
	return UserPlatformInfosStructTableName
}

//获取数据库表的主键字段名称.因为要兼容Map,只能是数据库的字段名称.对应的struct 属性field
func (entity *UserPlatformInfosStruct) GetPKColumnName() string {
	return "id"
}
