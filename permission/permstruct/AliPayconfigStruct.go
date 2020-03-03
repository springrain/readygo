package permstruct

import (
	"readygo/zorm"
	"time"
)

//AliPayconfigStructTableName 表名常量,方便直接调用
const AliPayconfigStructTableName = "ali_payconfig"

// AliPayconfigStruct 支付宝的配置信息
type AliPayconfigStruct struct {
	//引入默认的struct,隔离IEntityStruct的方法改动
	zorm.EntityStruct

	//Id <no value>
	Id string `column:"id"`

	//PrivateKey <no value>
	PrivateKey string `column:"privateKey"`

	//AliPayPublicKey <no value>
	AliPayPublicKey string `column:"aliPayPublicKey"`

	//AppId <no value>
	AppId string `column:"appId"`

	//ServiceUrl <no value>
	ServiceUrl string `column:"serviceUrl"`

	//Charset <no value>
	Charset string `column:"charset"`

	//SignType <no value>
	SignType string `column:"signType"`

	//Format <no value>
	Format string `column:"format"`

	//CertPath <no value>
	CertPath string `column:"certPath"`

	//AlipayPublicCertPath <no value>
	AlipayPublicCertPath string `column:"alipayPublicCertPath"`

	//RootCertPath <no value>
	RootCertPath string `column:"rootCertPath"`

	//EncryptType <no value>
	EncryptType string `column:"encryptType"`

	//AesKey <no value>
	AesKey string `column:"aesKey"`

	//CreateTime <no value>
	CreateTime time.Time `column:"createTime"`

	//CreateUserId <no value>
	CreateUserId string `column:"createUserId"`

	//UpdateTime <no value>
	UpdateTime time.Time `column:"updateTime"`

	//UpdateUserId <no value>
	UpdateUserId string `column:"updateUserId"`

	//Active 状态 0不可用,1可用
	Active int `column:"active"`

	//------------------数据库字段结束,自定义字段写在下面---------------//

}

//GetTableName 获取表名称
func (entity *AliPayconfigStruct) GetTableName() string {
	return AliPayconfigStructTableName
}

//GetPKColumnName 获取数据库表的主键字段名称.因为要兼容Map,只能是数据库的字段名称.
func (entity *AliPayconfigStruct) GetPKColumnName() string {
	return "id"
}
