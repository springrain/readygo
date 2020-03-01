package permstruct
import (
	"time"

	"readygo/orm"
)

//表名常量,方便直接调用
const MenuStructTableName = "菜单"

// 菜单
type MenuStruct struct {
	//引入默认的struct,隔离IEntityStruct的方法改动
	orm.EntityStruct
	
    // <no value>
    Id string `column:"id"`
    
    // 菜单名称
    Name string `column:"name"`
    
    // 代码
    Comcode string `column:"comcode"`
    
    // vue使用 meta.title
    Title sql.NullString `column:"title"`
    
    // <no value>
    Pid string `column:"pid"`
    
    // 备注
    Remark sql.NullString `column:"remark"`
    
    // <no value>
    Pageurl sql.NullString `column:"pageurl"`
    
    // 0.功能按钮,1.导航菜单
    MenuType int `column:"menuType"`
    
    // vue路由地址
    Path sql.NullString `column:"path"`
    
    // vue组件使用
    KeepAlive TINYINT `column:"keepAlive"`
    
    // vue组件使用
    Component sql.NullString `column:"component"`
    
    // vue组件使用
    Permission sql.NullString `column:"permission"`
    
    // vue组件使用
    Redirect sql.NullString `column:"redirect"`
    
    // <no value>
    Icon sql.NullString `column:"icon"`
    
    // <no value>
    CreateTime time.Time `column:"createTime"`
    
    // <no value>
    CreateUserId sql.NullString `column:"createUserId"`
    
    // <no value>
    UpdateTime time.Time `column:"updateTime"`
    
    // <no value>
    UpdateUserId sql.NullString `column:"updateUserId"`
    
    // 排序,查询时倒叙排列
    Sortno int `column:"sortno"`
    
    // 是否有效(0否,1是)
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
func (entity *MenuStruct) GetTableName() string {
	return MenuStructTableName
}

//获取数据库表的主键字段名称.因为要兼容Map,只能是数据库的字段名称.对应的struct 属性field
func (entity *MenuStruct) GetPKColumnName() string {
	return "id"
}

