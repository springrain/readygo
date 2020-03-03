package permstruct

import (
	"readygo/zorm"
	"time"
)

//MenuStructTableName 表名常量,方便直接调用
const MenuStructTableName = "t_menu"

// MenuStruct 菜单
type MenuStruct struct {
	//引入默认的struct,隔离IEntityStruct的方法改动
	zorm.EntityStruct

	//Id <no value>
	Id string `column:"id"`

	//Name 菜单名称
	Name string `column:"name"`

	//Comcode 代码
	Comcode string `column:"comcode"`

	//Pid <no value>
	Pid string `column:"pid"`

	//Remark 备注
	Remark string `column:"remark"`

	//Pageurl <no value>
	Pageurl string `column:"pageurl"`

	//MenuType 0.功能按钮,1.导航菜单
	MenuType int `column:"menuType"`

	//CreateTime <no value>
	CreateTime time.Time `column:"createTime"`

	//CreateUserId <no value>
	CreateUserId string `column:"createUserId"`

	//UpdateTime <no value>
	UpdateTime time.Time `column:"updateTime"`

	//UpdateUserId <no value>
	UpdateUserId string `column:"updateUserId"`

	//Sortno 排序,查询时倒叙排列
	Sortno int `column:"sortno"`

	//Active 是否有效(0否,1是)
	Active int `column:"active"`

	//------------------数据库字段结束,自定义字段写在下面---------------//

	//RoleId 菜单是由哪个角色产生的,用于强制部门权限的判定
	RoleId string

	//Children 菜单下的子菜单
	Children []MenuStruct
}

//GetTableName 获取表名称
func (entity *MenuStruct) GetTableName() string {
	return MenuStructTableName
}

//GetPKColumnName 获取数据库表的主键字段名称.因为要兼容Map,只能是数据库的字段名称.
func (entity *MenuStruct) GetPKColumnName() string {
	return "id"
}
