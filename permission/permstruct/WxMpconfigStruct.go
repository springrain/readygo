package permstruct

import (
	"time"

	"gitee.com/chunanyong/zorm"
)

//WxMpconfigStructTableName 表名常量,方便直接调用
const WxMpconfigStructTableName = "wx_mpconfig"

// WxMpconfigStruct 微信号需要的配置信息
type WxMpconfigStruct struct {
	//引入默认的struct,隔离IEntityStruct的方法改动
	zorm.EntityStruct

	//Id <no value>
	Id string `column:"id"`

	//OrgId 站点Id
	OrgId string `column:"orgId"`

	//AppId 开发者Id
	AppId string `column:"appId"`

	//Secret 应用密钥
	Secret string `column:"secret"`

	//Token 开发者令牌
	Token string `column:"token"`

	//AesKey 消息加解密密钥
	AesKey string `column:"aesKey"`

	//WxOriginalId 微信原始ID
	WxOriginalId string `column:"wxOriginalId"`

	//Oauth2 是否支持微信oauth2.0协议,0是不支持,1是支持
	Oauth2 int `column:"oauth2"`

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
func (entity *WxMpconfigStruct) GetTableName() string {
	return WxMpconfigStructTableName
}

//GetPKColumnName 获取数据库表的主键字段名称.因为要兼容Map,只能是数据库的字段名称.
func (entity *WxMpconfigStruct) GetPKColumnName() string {
	return "id"
}
