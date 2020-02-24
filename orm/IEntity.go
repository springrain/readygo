package orm

//Entity实体类接口,所有实体类必须实现,否则baseDao无法执行.baseDao函数形参只有Finder和IBaseEntity
type IEntityStruct interface {
	//获取表名称
	GetTableName() string
	//获取数据库表的主键字段名称.因为要兼容Map,只能是数据库的字段名称.对应的struct 属性field,使用cacheStructPKFieldNameMap缓存
	GetPKColumnName() string
	//兼容主键序列.如果有值,优先级最高
	GetPkSequence() string
}

//Entity实体类接口,所有实体类必须实现,否则baseDao无法执行.baseDao函数形参只有Finder和IBaseEntity
type IEntityMap interface {
	//获取表名称
	GetTableName() string
	//获取数据库表的主键字段名称.因为要兼容Map,只能是数据库的字段名称.对应的struct 属性field,使用cacheStructPKFieldNameMap缓存
	GetPKColumnName() string
	//针对Map类型,记录数据库字段
	GetDBFieldMap() map[string]interface{}
}

//IBaseEntity 的基础实现,所有的实体类都匿名注入.这样就类似实现继承了,如果接口增加方法,调整这个默认实现即可
type EntityStruct struct {
}

//默认数据库的主键列名
const defaultPkName = "id"

//获取表名称
/*
func (entity *EntityStruct) GetTableName() string {
	return ""
}
*/

//获取数据库表的主键字段名称.因为要兼容Map,只能是数据库的字段名称.对应的struct 属性field,使用cacheStructPKFieldNameMap缓存
func (entity *EntityStruct) GetPKColumnName() string {
	return defaultPkName
}

//Oracle和pgsql没有自增,主键使用序列.优先级高于GetPKColumnName方法
func (entity *EntityStruct) GetPkSequence() string {
	return ""
}

//-------------------------------------------------------------------------//

//IBaseEntity 的基础实现,所有的实体类都匿名注入.这样就类似实现继承了,如果接口增加方法,调整这个默认实现即可
type EntityMap struct {
	//表名
	tableName string
	//主键列名
	pkColumnName string
	//数据库字段
	DBFieldMap map[string]interface{}
	//自定义的kv
	TransientMap map[string]interface{}
}

//初始化Map,必须传入表名称
func NewEntityMap(tbName string) EntityMap {
	entityMap := EntityMap{}
	entityMap.DBFieldMap = map[string]interface{}{}
	entityMap.TransientMap = map[string]interface{}{}
	entityMap.tableName = tbName
	entityMap.pkColumnName = defaultPkName
	return entityMap
}

//获取表名称
func (entity *EntityMap) GetTableName() string {
	return entity.tableName
}

func (entity *EntityMap) SetPKColumnName(pkName string) {
	entity.pkColumnName = pkName
}

//获取数据库表的主键字段名称.因为要兼容Map,只能是数据库的字段名称.对应的struct 属性field,使用cacheStructPKFieldNameMap缓存
func (entity *EntityMap) GetPKColumnName() string {
	return entity.pkColumnName
}

//针对Map类型,记录数据库字段
func (entity *EntityMap) GetDBFieldMap() map[string]interface{} {
	return entity.DBFieldMap
}

//设置数据库字段
func (entity *EntityMap) Set(key string, value interface{}) map[string]interface{} {
	entity.DBFieldMap[key] = value
	return entity.DBFieldMap
}

//设置非数据库字段
func (entity *EntityMap) Put(key string, value interface{}) map[string]interface{} {
	entity.TransientMap[key] = value
	return entity.TransientMap
}
