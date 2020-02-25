package permstruct

import (
	"readygo/orm"
)

// 微信号需要的配置信息
type WxCpconfigStructStruct struct {
	//引入默认的struct,隔离IEntityStruct的方法改动
	orm.EntityStruct

	// <no value>
	Id string `column:"id"`

	// 站点Id
	OrgId string `column:"orgId"`

	// 开发者Id
	AppId string `column:"appId"`

	// 应用密钥
	Secret string `column:"secret"`

	// 状态 0不可用,1可用
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
func (entity *WxCpconfigStructStruct) GetTableName() string {
	return "wx_cpconfig"
}

//获取数据库表的主键字段名称.因为要兼容Map,只能是数据库的字段名称.对应的struct 属性field
func (entity *WxCpconfigStructStruct) GetPKColumnName() string {
	return "id"
}
