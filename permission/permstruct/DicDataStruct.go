package permstruct

import (
	"readygo/zorm"
	"time"
)

//DicDataStructTableName 表名常量,方便直接调用
const DicDataStructTableName = "t_dic_data"

// DicDataStruct 公共字典
type DicDataStruct struct {
	//引入默认的struct,隔离IEntityStruct的方法改动
	zorm.EntityStruct

	//Id <no value>
	Id string `column:"id"`

	//Name 名称
	Name string `column:"name"`

	//Code 编码
	Code string `column:"code"`

	//Val 值
	Val string `column:"val"`

	//Pid 父ID
	Pid string `column:"pid"`

	//Remark 描述
	Remark string `column:"remark"`

	//Typekey 类型
	Typekey string `column:"typekey"`

	//CreateTime <no value>
	CreateTime time.Time `column:"createTime"`

	//CreateUserId <no value>
	CreateUserId string `column:"createUserId"`

	//UpdateTime <no value>
	UpdateTime time.Time `column:"updateTime"`

	//UpdateUserId <no value>
	UpdateUserId string `column:"updateUserId"`

	//Sortno 排序
	Sortno int `column:"sortno"`

	//Active 是否有效(0否,1是)
	Active int `column:"active"`

	//------------------数据库字段结束,自定义字段写在下面---------------//

}

//GetTableName 获取表名称
func (entity *DicDataStruct) GetTableName() string {
	return DicDataStructTableName
}

//GetPKColumnName 获取数据库表的主键字段名称.因为要兼容Map,只能是数据库的字段名称.
func (entity *DicDataStruct) GetPKColumnName() string {
	return "id"
}
