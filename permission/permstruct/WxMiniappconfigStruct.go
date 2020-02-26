package permstruct

import (
	"time"

	"readygo/orm"
)

//WxMiniappconfigStructTableName 表名常量,方便直接调用
const WxMiniappconfigStructTableName = "wx_miniappconfig"

// WxMiniappconfigStruct 小程序配置表
type WxMiniappconfigStruct struct {
	//引入默认的struct,隔离IEntityStruct的方法改动
	orm.EntityStruct

	// 主键id
	Id string `column:"id"`

	// 站点Id
	OrgId string `column:"orgId"`

	// 开发者Id
	AppId string `column:"appId"`

	// 应用密钥
	Secret string `column:"secret"`

	// 签约模板Id
	PlanId string `column:"planId"`

	// 签约请求序列号
	RequestSerial string `column:"requestSerial"`

	// <no value>
	CreateTime time.Time `column:"createTime"`

	// <no value>
	CreateUserId string `column:"createUserId"`

	// <no value>
	UpdateTime time.Time `column:"updateTime"`

	// <no value>
	UpdateUserId string `column:"updateUserId"`

	// 状态 0不可用,1可用
	Active int `column:"active"`

	//------------------数据库字段结束,自定义字段写在下面---------------//

}

//获取表名称
func (entity *WxMiniappconfigStruct) GetTableName() string {
	return WxMiniappconfigStructTableName
}

//获取数据库表的主键字段名称.因为要兼容Map,只能是数据库的字段名称.对应的struct 属性field
func (entity *WxMiniappconfigStruct) GetPKColumnName() string {
	return "id"
}
