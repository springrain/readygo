package permstruct

import (
	"database/sql"
	"readygo/orm"
)

//表名常量,方便直接调用
const AliPayconfigStructTableName = "支付宝的配置信息"

// 支付宝的配置信息
type AliPayconfigStruct struct {
	//引入默认的struct,隔离IEntityStruct的方法改动
	orm.EntityStruct
	
    // <no value>
    Id string `column:"id"`
    
    // <no value>
    PrivateKey sql.NullString `column:"privateKey"`
    
    // <no value>
    AliPayPublicKey sql.NullString `column:"aliPayPublicKey"`
    
    // <no value>
    AppId sql.NullString `column:"appId"`
    
    // <no value>
    ServiceUrl string `column:"serviceUrl"`
    
    // <no value>
    Charset string `column:"charset"`
    
    // <no value>
    SignType string `column:"signType"`
    
    // <no value>
    Format string `column:"format"`
    
    // <no value>
    CertPath sql.NullString `column:"certPath"`
    
    // <no value>
    AlipayPublicCertPath sql.NullString `column:"alipayPublicCertPath"`
    
    // <no value>
    RootCertPath sql.NullString `column:"rootCertPath"`
    
    // <no value>
    EncryptType sql.NullString `column:"encryptType"`
    
    // <no value>
    AesKey sql.NullString `column:"aesKey"`
    
	//------------------数据库字段结束,自定义字段写在下面---------------//


}


//获取表名称
func (entity *AliPayconfigStruct) GetTableName() string {
	return AliPayconfigStructTableName
}

//获取数据库表的主键字段名称.因为要兼容Map,只能是数据库的字段名称.对应的struct 属性field
func (entity *AliPayconfigStruct) GetPKColumnName() string {
	return "id"
}

