package permstruct

import (
	"time"

	"gitee.com/chunanyong/zorm"
)

//WxCpconfigStructTableName 表名常量,方便直接调用
const WxCpconfigStructTableName = "wx_cpconfig"

// WxCpconfigStruct 微信号需要的配置信息
type WxCpconfigStruct struct {
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
func (entity *WxCpconfigStruct) GetTableName() string {
	return WxCpconfigStructTableName
}

//GetPKColumnName 获取数据库表的主键字段名称.因为要兼容Map,只能是数据库的字段名称.
func (entity *WxCpconfigStruct) GetPKColumnName() string {
	return "id"
}
