package peristruct

// 公共字典
type DicDataStruct struct {

	// <no value>
	Id string `column:"id"`

	// 名称
	Name string `column:"name"`

	// 编码
	Code string `column:"code"`

	// 值
	Val string `column:"val"`

	// 父ID
	Pid string `column:"pid"`

	// 排序
	Sortno int `column:"sortno"`

	// 描述
	Remark string `column:"remark"`

	// 是否有效(0否,1是)
	Active int `column:"active"`

	// 类型
	Typekey string `column:"typekey"`

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
func (entity *DicDataStruct) GetTableName() string {
	return "t_dic_data"
}

//获取数据库表的主键字段名称.因为要兼容Map,只能是数据库的字段名称.对应的struct 属性field,使用cacheStructPKFieldNameMap缓存
func (entity *DicDataStruct) GetPKColumnName() string {
	return "id"
}

//Oracle和pgsql没有自增,主键使用序列.优先级高于GetPKColumnName方法
func (entity *DicDataStruct) GetPkSequence() string {
	return ""
}
