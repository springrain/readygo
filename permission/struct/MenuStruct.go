package peristruct

import (
	"time"
)

// 菜单
type MenuStruct struct {

	// <no value>
	Id string `column:"id"`

	// 菜单名称
	Name string `column:"name"`

	// 代码
	Comcode string `column:"comcode"`

	// vue使用 meta.title
	Title string `column:"title"`

	// <no value>
	Pid string `column:"pid"`

	// 备注
	Remark string `column:"remark"`

	// <no value>
	Pageurl string `column:"pageurl"`

	// 0.功能按钮,1.导航菜单
	MenuType int `column:"menuType"`

	// vue路由地址
	Path string `column:"path"`

	// <no value>
	CreateTime time.Time `column:"createTime"`

	// <no value>
	CreateUserId string `column:"createUserId"`

	// <no value>
	UpdateTime time.Time `column:"updateTime"`

	// <no value>
	UpdateUserId string `column:"updateUserId"`

	// 排序,查询时倒叙排列
	Sortno int `column:"sortno"`

	// 是否有效(0否,1是)
	Active int `column:"active"`

	// <no value>
	Bak1 string `column:"bak1"`

	// <no value>
	Bak2 string `column:"bak2"`

	// <no value>
	Bak3 string `column:"bak3"`

	// <no value>
	Bak4 string `column:"bak4"`

	// <no value>
	Bak5 string `column:"bak5"`

	//------------------数据库字段结束,自定义字段写在下面---------------//

}

//获取表名称
func (entity *MenuStruct) GetTableName() string {
	return "t_menu"
}

//获取数据库表的主键字段名称.因为要兼容Map,只能是数据库的字段名称.对应的struct 属性field,使用cacheStructPKFieldNameMap缓存
func (entity *MenuStruct) GetPKColumnName() string {
	return "id"
}

//Oracle和pgsql没有自增,主键使用序列.优先级高于GetPKColumnName方法
func (entity *MenuStruct) GetPkSequence() string {
	return ""
}
