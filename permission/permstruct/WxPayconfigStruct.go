package permstruct

import (
	"readygo/zorm"
	"time"
)

//WxPayconfigStructTableName 表名常量,方便直接调用
const WxPayconfigStructTableName = "wx_payconfig"

// WxPayconfigStruct 微信号需要的配置信息
type WxPayconfigStruct struct {
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

	//MchId 微信支付商户号
	MchId string `column:"mchId"`

	//Key 交易过程生成签名的密钥，仅保留在商户系统和微信支付后台，不会在网络中传播
	Key string `column:"key"`

	//CertificateFile 证书地址
	CertificateFile string `column:"certificateFile"`

	//NotifyUrl 通知地址
	NotifyUrl string `column:"notifyUrl"`

	//SignType 加密方式,MD5和HMAC-SHA256
	SignType string `column:"signType"`

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
func (entity *WxPayconfigStruct) GetTableName() string {
	return WxPayconfigStructTableName
}

//GetPKColumnName 获取数据库表的主键字段名称.因为要兼容Map,只能是数据库的字段名称.
func (entity *WxPayconfigStruct) GetPKColumnName() string {
	return "id"
}
