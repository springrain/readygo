package peristruct

// 微信号需要的配置信息
type WxPayconfigStruct struct {

	// <no value>
	Id string `column:"id"`

	// 站点Id
	OrgId string `column:"orgId"`

	// 开发者Id
	AppId string `column:"appId"`

	// 应用密钥
	Secret string `column:"secret"`

	// 微信支付商户号
	MchId string `column:"mchId"`

	// 交易过程生成签名的密钥，仅保留在商户系统和微信支付后台，不会在网络中传播
	Key string `column:"key"`

	// 证书地址
	CertificateFile string `column:"certificateFile"`

	// 通知地址
	NotifyUrl string `column:"notifyUrl"`

	// 加密方式,MD5和HMAC-SHA256
	SignType string `column:"signType"`

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
func (entity *WxPayconfigStruct) GetTableName() string {
	return "wx_payconfig"
}

//获取数据库表的主键字段名称.因为要兼容Map,只能是数据库的字段名称.对应的struct 属性field
func (entity *WxPayconfigStruct) GetPKColumnName() string {
	return "id"
}

//Oracle和pgsql没有自增,主键使用序列.优先级高于GetPKColumnName方法
func (entity *WxPayconfigStruct) GetPkSequence() string {
	return ""
}
