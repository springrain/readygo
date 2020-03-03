package permservice

import (
	"errors"
	"readygo/cache"
	"readygo/permission/permstruct"
	"readygo/zorm"
)

const (
	userOrgRoleMenuCacheKey string = "userOrgRoleMenuCacheKey"
	qxCacheKey              string = "qxCacheKey"
)

//FindRoleByUserId 根据用户Id查询用户的角色
func FindRoleByUserId(session *zorm.Session, userId string) ([]permstruct.RoleStruct, error) {
	if len(userId) < 1 {
		return nil, errors.New("参数userId不能为空")
	}

	cacheKey := "FindRoleByUserId_" + userId

	roles := make([]permstruct.RoleStruct, 0)

	//从缓存中取数据
	errFromCache := cache.GetFromCache(userOrgRoleMenuCacheKey, cacheKey, &roles)
	if errFromCache != nil {
		return nil, errFromCache
	}
	if len(roles) > 0 { //缓存中有数据
		return roles, nil
	}
	//按照 r.privateOrg,r.sortno desc 先处理强制部门权限的角色
	finder := zorm.NewFinder()
	finder.Append("SELECT r.* from ").Append(permstruct.RoleStructTableName).Append(" r,")
	finder.Append(permstruct.UserRoleStructTableName).Append("  re where re.userId=? and re.roleId=r.id and r.active=1 order by r.privateOrg,r.sortno desc", userId)

	//查询列表
	errQueryList := zorm.QueryStructList(session, finder, &roles, nil)
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

//FindMenuByRoleId 根据角色Id,查询这个角色有权限的菜单
func FindMenuByRoleId(session *zorm.Session, roleId string) ([]permstruct.MenuStruct, error) {

	if len(roleId) < 1 {
		return nil, errors.New("roleId的值不能为空")
	}
	cacheKey := "FindMenuByRoleId_" + roleId
	menus := make([]permstruct.MenuStruct, 0)

	//从缓存中取数据
	errFromCache := cache.GetFromCache(qxCacheKey, cacheKey, &menus)
	if errFromCache != nil {
		return nil, errFromCache
	}
	if len(menus) > 0 { //缓存中有数据
		return menus, nil
	}
	//查询角色有权限的菜单
	finder := zorm.NewFinder()
	finder.Append("SELECT m.* from ").Append(permstruct.MenuStructTableName).Append(" m,")
	finder.Append(permstruct.RoleMenuStructTableName).Append("  re where re.roleId=? and re.menuId=m.id and m.active=1 order by m.sortno desc ", roleId)

	//查询列表
	errQueryList := zorm.QueryStructList(session, finder, &menus, nil)
	if errQueryList != nil {
		return nil, errQueryList
	}

	//放入缓存
	errPutCache := cache.PutToCache(qxCacheKey, cacheKey, &menus)
	if errPutCache != nil {
		return nil, errPutCache
	}

	return menus, nil

}
