package peristruct

// 小程序配置表
type WxMiniappconfigStruct struct {

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
func (entity *WxMiniappconfigStruct) GetTableName() string {
	return "wx_miniappconfig"
}

//获取数据库表的主键字段名称.因为要兼容Map,只能是数据库的字段名称.对应的struct 属性field,使用cacheStructPKFieldNameMap缓存
func (entity *WxMiniappconfigStruct) GetPKColumnName() string {
	return "id"
}

//Oracle和pgsql没有自增,主键使用序列.优先级高于GetPKColumnName方法
func (entity *WxMiniappconfigStruct) GetPkSequence() string {
	return ""
}
