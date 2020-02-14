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
	//获取主键名称,需要兼容Map,所以不放到tag里了
	GetPkName() string
	//兼容主键序列.如果有值,优先级最高
	GetPkSequence() string
	//Struct对象类型和Map类型.两者都是Struct类型,一个是对象载体需要反射,一个是Map载体,不需要反射
	GetEntityType() EntityType

	//针对Map类型,记录数据库字段
	GetDBFieldMap() map[string]interface{}
	//针对Map类型,记录非数据库字段
	GetTransientMap() map[string]interface{}
}

//IBaseEntity 的基础实现,所有的实体类都匿名注入.这样就类似实现继承了,如果接口增加方法,调整这个默认实现即可
type EntityStruct struct {
}

//获取表名称
func (entity *EntityStruct) GetTableName() string {
	return ""
}

//获取主键名称,需要兼容Map,所以不放到tag里了
func (entity *EntityStruct) GetPkName() string {
	return "Id"
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

//针对Map类型,记录非数据库字段
func (entity *EntityStruct) GetTransientMap() map[string]interface{} {
	return nil
}

//-------------------------------------------------------------------------//

//IBaseEntity 的基础实现,所有的实体类都匿名注入.这样就类似实现继承了,如果接口增加方法,调整这个默认实现即可
type EntityMap struct {
	DBFieldMap   map[string]interface{}
	TransientMap map[string]interface{}
}

//获取表名称
func (entity *EntityMap) GetTableName() string {
	return ""
}

//获取主键名称,需要兼容Map,所以不放到tag里了
func (entity *EntityMap) GetPkName() string {
	return "Id"
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

//针对Map类型,记录非数据库字段
func (entity *EntityMap) GetTransientMap() map[string]interface{} {
	return entity.TransientMap
}
