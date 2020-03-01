package permstruct
import (
	"time"

	"readygo/orm"
)

//表名常量,方便直接调用
const WxMpconfigStructTableName = "微信号需要的配置信息"

// 微信号需要的配置信息
type WxMpconfigStruct struct {
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
    
    // 开发者令牌
    Token sql.NullString `column:"token"`
    
    // 消息加解密密钥
    AesKey sql.NullString `column:"aesKey"`
    
    // 微信原始ID
    WxOriginalId sql.NullString `column:"wxOriginalId"`
    
    // 是否支持微信oauth2.0协议,0是不支持,1是支持
    Oauth2 sql.NullInt32 `column:"oauth2"`
    
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
func (entity *WxMpconfigStruct) GetTableName() string {
	return WxMpconfigStructTableName
}

//获取数据库表的主键字段名称.因为要兼容Map,只能是数据库的字段名称.对应的struct 属性field
func (entity *WxMpconfigStruct) GetPKColumnName() string {
	return "id"
}

