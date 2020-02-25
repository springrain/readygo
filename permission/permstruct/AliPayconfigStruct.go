package permstruct

// 支付宝的配置信息
type AliPayconfigStruct struct {

	// <no value>
	Id string `column:"id"`

	// <no value>
	PrivateKey string `column:"privateKey"`

	// <no value>
	AliPayPublicKey string `column:"aliPayPublicKey"`

	// <no value>
	AppId string `column:"appId"`

	// <no value>
	ServiceUrl string `column:"serviceUrl"`

	// <no value>
	Charset string `column:"charset"`

	// <no value>
	SignType string `column:"signType"`

	// <no value>
	Format string `column:"format"`

	// <no value>
	CertPath string `column:"certPath"`

	// <no value>
	AlipayPublicCertPath string `column:"alipayPublicCertPath"`

	// <no value>
	RootCertPath string `column:"rootCertPath"`

	// <no value>
	EncryptType string `column:"encryptType"`

	// <no value>
	AesKey string `column:"aesKey"`

	//------------------数据库字段结束,自定义字段写在下面---------------//

}

//获取表名称
func (entity *AliPayconfigStruct) GetTableName() string {
	return "ali_payconfig"
}

//获取数据库表的主键字段名称.因为要兼容Map,只能是数据库的字段名称.对应的struct 属性field
func (entity *AliPayconfigStruct) GetPKColumnName() string {
	return "id"
}

//Oracle和pgsql没有自增,主键使用序列.优先级高于GetPKColumnName方法
func (entity *AliPayconfigStruct) GetPkSequence() string {
	return ""
}
