package orm

type EntityType int

//Struct对象类型和Map类型.两者都是Struct类型,一个是对象载体需要反射,一个是Map载体,不需要反射
const (
	EntityType_struct EntityType = 0
	EntityType_Map    EntityType = 1
)

//Entity实体类接口,所有实体类必须实现,否则baseDao无法执行.baseDao函数形参只有Finder和IBaseEntity
type IBaseEntity interface {
	//获取表名称
	GetTableName() string
	//获取数据库表的主键字段名称.因为要兼容Map,只能是数据库的字段名称.对应的struct 属性field,使用cacheStructPKFieldNameMap缓存
	GetPKColumnName() string
	//兼容主键序列.如果有值,优先级最高
	GetPkSequence() string
	//Struct对象类型和Map类型.两者都是Struct类型,一个是对象载体需要反射,一个是Map载体,不需要反射
	GetEntityType() EntityType

	//针对Map类型,记录数据库字段
	GetDBFieldMap() map[string]interface{}
}

//IBaseEntity 的基础实现,所有的实体类都匿名注入.这样就类似实现继承了,如果接口增加方法,调整这个默认实现即可
type EntityStruct struct {
}

//获取表名称
func (entity *EntityStruct) GetTableName() string {
	return ""
}

//获取数据库表的主键字段名称.因为要兼容Map,只能是数据库的字段名称.对应的struct 属性field,使用cacheStructPKFieldNameMap缓存
func (entity *EntityStruct) GetPKColumnName() string {
	return "id"
}

//兼容主键序列.如果有值,优先级最高
func (entity *EntityStruct) GetPkSequence() string {
	return ""
}

//Struct对象类型和Map类型.两者都是Struct类型,一个是对象载体需要反射,一个是Map载体,不需要反射
func (entity *EntityStruct) GetEntityType() EntityType {
	return EntityType_struct
}

//针对Map类型,记录数据库字段
func (entity *EntityStruct) GetDBFieldMap() map[string]interface{} {
	return nil
}

//-------------------------------------------------------------------------//

//IBaseEntity 的基础实现,所有的实体类都匿名注入.这样就类似实现继承了,如果接口增加方法,调整这个默认实现即可
type EntityMap struct {
	//记得init初始化
	//数据库字段
	DBFieldMap map[string]interface{}
	//自定义的kv
	TransientMap map[string]interface{}
}

//获取表名称
func (entity *EntityMap) GetTableName() string {
	return ""
}

//获取数据库表的主键字段名称.因为要兼容Map,只能是数据库的字段名称.对应的struct 属性field,使用cacheStructPKFieldNameMap缓存
func (entity *EntityMap) GetPKColumnName() string {
	return "id"
}

//兼容主键序列.如果有值,优先级最高
func (entity *EntityMap) GetPkSequence() string {
	return ""
}

//Struct对象类型和Map类型.两者都是Struct类型,一个是对象载体需要反射,一个是Map载体,不需要反射
func (entity *EntityMap) GetEntityType() EntityType {
	return EntityType_Map
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
