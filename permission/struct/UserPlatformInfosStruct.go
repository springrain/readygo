package permstruct

// 用户平台信息表
type UserPlatformInfosStruct struct {

	// 主键id
	Id string `column:"id"`

	// 公众号openId,企业号userId,小程序openId,APP推送deviceToken
	OpenId string `column:"openId"`

	// 设备/应用类型：1公众号2小程序3企业号4APP IOS消息推送5APP安卓消息推送6web
	DeviceType int `column:"deviceType"`

	// 所属站点ID
	SiteId string `column:"siteId"`

	// t_user表中ID
	UserId string `column:"userId"`

	// <no value>
	Bak1 string `column:"bak1"`

	// <no value>
	Bak2 string `column:"bak2"`

	// <no value>
	Bak3 string `column:"bak3"`

	// <no value>
	Bak4 string `column:"bak4"`

	//------------------数据库字段结束,自定义字段写在下面---------------//

}

//获取表名称
func (entity *UserPlatformInfosStruct) GetTableName() string {
	return "t_user_platform_infos"
}

//获取数据库表的主键字段名称.因为要兼容Map,只能是数据库的字段名称.对应的struct 属性field
func (entity *UserPlatformInfosStruct) GetPKColumnName() string {
	return "id"
}

//Oracle和pgsql没有自增,主键使用序列.优先级高于GetPKColumnName方法
func (entity *UserPlatformInfosStruct) GetPkSequence() string {
	return ""
}
