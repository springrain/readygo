package permstruct

import (
	"time"

	"gitee.com/chunanyong/zorm"
)

// OrgStructTableName 表名常量,方便直接调用
const OrgStructTableName = "t_org"

// OrgStruct 部门
type OrgStruct struct {
	// 引入默认的struct,隔离IEntityStruct的方法改动
	zorm.EntityStruct

	// Id 编号
	Id string `column:"id"`

	// Name 名称
	Name string `column:"name"`

	// Comcode 代码
	Comcode string `column:"comcode"`

	// Pid 上级部门ID
	Pid string `column:"pid"`

	// OrgType 0-99门店,100-199部门,200-299,分公司,300-399集团公司,900-999总平台
	OrgType int `column:"org_type"`

	// Sortno 排序,查询时倒叙排列
	Sortno int `column:"sortno"`

	// Remark 备注
	Remark string `column:"remark"`

	// CreateTime <no value>
	CreateTime time.Time `column:"create_time"`

	// CreateUserId <no value>
	CreateUserId string `column:"create_user_id"`

	// UpdateTime <no value>
	UpdateTime time.Time `column:"update_time"`

	// UpdateUserId <no value>
	UpdateUserId string `column:"update_user_id"`

	// Active 是否有效(0否,1是)
	Active int `column:"active"`

	//------------------数据库字段结束,自定义字段写在下面---------------//

	// Children 子部门
	Children []OrgStruct
}

// GetTableName 获取表名称
func (entity *OrgStruct) GetTableName() string {
	return OrgStructTableName
}

// GetPKColumnName 获取数据库表的主键字段名称.因为要兼容Map,只能是数据库的字段名称.
func (entity *OrgStruct) GetPKColumnName() string {
	return "id"
}
