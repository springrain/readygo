package permservice

import (
	"errors"
	"readygo/cache"
	"readygo/orm"
	"readygo/permission/permstruct"
)

const userOrgRoleMenuCacheKey string = "userOrgRoleMenuCacheKey"

//FindRoleByUserId 根据用户Id查询用户的角色
func FindRoleByUserId(session *orm.Session, userId string) ([]permstruct.RoleStruct, error) {
	if len(userId) < 1 {
		return nil, errors.New("参数userId不能为空")
	}

	cacheKey := "findRoleByUserId_" + userId

	roles := make([]permstruct.RoleStruct, 0)

	errFromCache := cache.GetFromCache(userOrgRoleMenuCacheKey, cacheKey, &roles)
	if errFromCache != nil {
		return nil, errFromCache
	}
	if len(roles) > 0 { //缓存中有数据
		return roles, nil
	}
	//按照 r.privateOrg,r.sortno desc 先处理强制部门权限的角色
	finder := orm.NewFinder()
	finder.Append("SELECT r.* from ").Append(permstruct.RoleStructTableName).Append(" r,")
	finder.Append(permstruct.UserRoleStructTableName).Append("  re where re.userId=? and re.roleId=r.id and r.active=1 order by r.privateOrg,r.sortno desc", userId)

	//查询列表
	errQueryList := orm.QueryStructList(session, finder, roles, nil)
	if errQueryList != nil {
		return nil, errQueryList
	}

	//放入缓存
	errPutCache := cache.PutToCache(userOrgRoleMenuCacheKey, cacheKey, &roles)
	if errPutCache != nil {
		return nil, errPutCache
	}

	return roles, nil

}
