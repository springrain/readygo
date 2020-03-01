package permstruct
import (
	"time"

	"readygo/orm"
)

//表名常量,方便直接调用
const DicDataStructTableName = "公共字典"

// 公共字典
type DicDataStruct struct {
	//引入默认的struct,隔离IEntityStruct的方法改动
	orm.EntityStruct
	
    // <no value>
    Id string `column:"id"`
    
    // 名称
    Name string `column:"name"`
    
    // 编码
    Code sql.NullString `column:"code"`
    
    // 值
    Val sql.NullString `column:"val"`
    
    // 父ID
    Pid sql.NullString `column:"pid"`
    
    // 排序
    Sortno sql.NullInt32 `column:"sortno"`
    
    // 描述
    Remark sql.NullString `column:"remark"`
    
    // 是否有效(0否,1是)
    Active sql.NullInt32 `column:"active"`
    
    // 类型
    Typekey sql.NullString `column:"typekey"`
    
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
func (entity *DicDataStruct) GetTableName() string {
	return DicDataStructTableName
}

//获取数据库表的主键字段名称.因为要兼容Map,只能是数据库的字段名称.对应的struct 属性field
func (entity *DicDataStruct) GetPKColumnName() string {
	return "id"
}

