package permstruct
import (
	"time"

	"readygo/orm"
)

//表名常量,方便直接调用
const WxMiniappconfigStructTableName = "小程序配置表"

// 小程序配置表
type WxMiniappconfigStruct struct {
	//引入默认的struct,隔离IEntityStruct的方法改动
	orm.EntityStruct
	
    // 主键id
    Id string `column:"id"`
    
    // 站点Id
    OrgId string `column:"orgId"`
    
    // 开发者Id
    AppId string `column:"appId"`
    
    // 应用密钥
    Secret string `column:"secret"`
    
    // 签约模板Id
    PlanId sql.NullString `column:"planId"`
    
    // 签约请求序列号
    RequestSerial sql.NullString `column:"requestSerial"`
    
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
func (entity *WxMiniappconfigStruct) GetTableName() string {
	return WxMiniappconfigStructTableName
}

//获取数据库表的主键字段名称.因为要兼容Map,只能是数据库的字段名称.对应的struct 属性field
func (entity *WxMiniappconfigStruct) GetPKColumnName() string {
	return "id"
}

