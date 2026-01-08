package permstruct

import (
	"time"

	"gitee.com/chunanyong/zorm"
)

// AliPayconfigStructTableName 表名常量,方便直接调用
const AliPayconfigStructTableName = "ali_payconfig"

// AliPayconfigStruct 支付宝的配置信息
type AliPayconfigStruct struct {
	// 引入默认的struct,隔离IEntityStruct的方法改动
	zorm.EntityStruct

	// Id <no value>
	Id string `column:"id"`

	// PrivateKey <no value>
	PrivateKey string `column:"private_key"`

	// AliPayPublicKey <no value>
	AliPayPublicKey string `column:"ali_pay_public_key"`

	// AppId <no value>
	AppId string `column:"app_id"`

	// ServiceUrl <no value>
	ServiceUrl string `column:"service_url"`

	// Charset <no value>
	Charset string `column:"charset"`

	// SignType <no value>
	SignType string `column:"sign_type"`

	// Format <no value>
	Format string `column:"format"`

	// CertPath <no value>
	CertPath string `column:"cert_path"`

	// AlipayPublicCertPath <no value>
	AlipayPublicCertPath string `column:"alipay_public_cert_path"`

	// RootCertPath <no value>
	RootCertPath string `column:"root_cert_path"`

	// EncryptType <no value>
	EncryptType string `column:"encrypt_type"`

	// AesKey <no value>
	AesKey string `column:"aes_key"`

	// CreateTime <no value>
	CreateTime time.Time `column:"create_time"`

	// CreateUserId <no value>
	CreateUserId string `column:"create_user_id"`

	// UpdateTime <no value>
	UpdateTime time.Time `column:"update_time"`

	// UpdateUserId <no value>
	UpdateUserId string `column:"update_user_id"`

	// Active 状态 0不可用,1可用
	Active int `column:"active"`

	//------------------数据库字段结束,自定义字段写在下面---------------//
}

// GetTableName 获取表名称
func (entity *AliPayconfigStruct) GetTableName() string {
	return AliPayconfigStructTableName
}

// GetPKColumnName 获取数据库表的主键字段名称.因为要兼容Map,只能是数据库的字段名称.
func (entity *AliPayconfigStruct) GetPKColumnName() string {
	return "id"
}
