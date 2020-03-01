package permstruct
import (
	"time"

	"readygo/orm"
)

//表名常量,方便直接调用
const WxPayconfigStructTableName = "微信号需要的配置信息"

// 微信号需要的配置信息
type WxPayconfigStruct struct {
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
    
    // 微信支付商户号
    MchId string `column:"mchId"`
    
    // 交易过程生成签名的密钥，仅保留在商户系统和微信支付后台，不会在网络中传播
    Key string `column:"key"`
    
    // 证书地址
    CertificateFile string `column:"certificateFile"`
    
    // 通知地址
    NotifyUrl sql.NullString `column:"notifyUrl"`
    
    // 加密方式,MD5和HMAC-SHA256
    SignType sql.NullString `column:"signType"`
    
    // 状态 0不可用,1可用
    Active int `column:"active"`
    
    // <no value>
    Bak1 sql.NullString `column:"bak1"`
    
    // <no value>
    Bak2 sql.NullString `column:"bak2"`
    
    // <no value>
    Bak3 sql.NullString `column:"bak3"`
    
    // <no value>
    Bak4 sql.NullString `column:"bak4"`
    
    // <no value>
    Bak5 sql.NullString `column:"bak5"`
    
	//------------------数据库字段结束,自定义字段写在下面---------------//


}


//获取表名称
func (entity *WxPayconfigStruct) GetTableName() string {
	return WxPayconfigStructTableName
}

//获取数据库表的主键字段名称.因为要兼容Map,只能是数据库的字段名称.对应的struct 属性field
func (entity *WxPayconfigStruct) GetPKColumnName() string {
	return "id"
}

