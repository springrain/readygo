package permstruct
import (
	"time"

	"readygo/orm"
)

//表名常量,方便直接调用
const UserPlatformInfosStructTableName = "用户平台信息表"

// 用户平台信息表
type UserPlatformInfosStruct struct {
	//引入默认的struct,隔离IEntityStruct的方法改动
	orm.EntityStruct
	
    // 主键id
    Id string `column:"id"`
    
    // 公众号openId,企业号userId,小程序openId,APP推送deviceToken
    OpenId sql.NullString `column:"openId"`
    
    // 设备/应用类型：1公众号2小程序3企业号4APP IOS消息推送5APP安卓消息推送6web
    DeviceType sql.NullInt32 `column:"deviceType"`
    
    // 所属站点ID
    SiteId sql.NullString `column:"siteId"`
    
    // t_user表中ID
    UserId sql.NullString `column:"userId"`
    
    // <no value>
    Bak1 sql.NullString `column:"bak1"`
    
    // <no value>
    Bak2 sql.NullString `column:"bak2"`
    
    // <no value>
    Bak3 sql.NullString `column:"bak3"`
    
    // <no value>
    Bak4 sql.NullString `column:"bak4"`
    
	//------------------数据库字段结束,自定义字段写在下面---------------//


}


//获取表名称
func (entity *UserPlatformInfosStruct) GetTableName() string {
	return UserPlatformInfosStructTableName
}

//获取数据库表的主键字段名称.因为要兼容Map,只能是数据库的字段名称.对应的struct 属性field
func (entity *UserPlatformInfosStruct) GetPKColumnName() string {
	return "id"
}

