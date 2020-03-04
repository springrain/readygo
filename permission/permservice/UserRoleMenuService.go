package permservice

import (
	"errors"
	"readygo/cache"
	"readygo/permission/permstruct"
	"readygo/zorm"
)

const (
	//	userOrgRoleMenuCacheKey string = "userOrgRoleMenuCacheKey"
	qxCacheKey string = "qxCacheKey"
)

//FindRoleByUserId 根据用户Id查询用户的角色
func FindRoleByUserId(dbConnection *zorm.DBConnection, userId string, page *zorm.Page) ([]permstruct.RoleStruct, error) {
	if len(userId) < 1 {
		return nil, errors.New("参数userId不能为空")
	}

	cacheKey := "FindRoleByUserId_" + userId

	roles := make([]permstruct.RoleStruct, 0)

	//从缓存中取数据
	errFromCache := cache.GetFromCache(qxCacheKey, cacheKey, &roles)
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
	errQueryList := zorm.QueryStructList(dbConnection, finder, &roles, page)
	if errQueryList != nil {
		return nil, errQueryList
	}

	//放入缓存
	errPutCache := cache.PutToCache(qxCacheKey, cacheKey, roles)
	if errPutCache != nil {
		return nil, errPutCache
	}

	return roles, nil

}

//FindMenuByRoleId 根据角色Id,查询这个角色有权限的菜单
func FindMenuByRoleId(dbConnection *zorm.DBConnection, roleId string, page *zorm.Page) ([]permstruct.MenuStruct, error) {

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
	errQueryList := zorm.QueryStructList(dbConnection, finder, &menus, page)
	if errQueryList != nil {
		return nil, errQueryList
	}

	//放入缓存
	errPutCache := cache.PutToCache(qxCacheKey, cacheKey, menus)
	if errPutCache != nil {
		return nil, errPutCache
	}

	return menus, nil

}

//FindMenuByUserId 查询用户有权限的菜单
func FindMenuByUserId(dbConnection *zorm.DBConnection, userId string) ([]permstruct.MenuStruct, error) {

	if len(userId) < 1 {
		return nil, errors.New("roleId的值不能为空")
	}
	cacheKey := "FindMenuByUserId_" + userId
	menus := make([]permstruct.MenuStruct, 0)

	//从缓存中取数据
	errFromCache := cache.GetFromCache(qxCacheKey, cacheKey, &menus)
	if errFromCache != nil {
		return nil, errFromCache
	}
	if len(menus) > 0 { //缓存中有数据
		return menus, nil
	}

	//查询用户所有的角色
	roles, errFindRoleByUserId := FindRoleByUserId(dbConnection, userId, nil)
	if errFindRoleByUserId != nil {
		return nil, errFindRoleByUserId
	}

	//用户没有角色
	if len(roles) < 1 {
		return nil, nil
	}

	//去重map
	menusMap := make(map[string]permstruct.MenuStruct)

	//循环所有的角色
	for _, role := range roles {
		menusByRoleId, errFindMenuByRoleId := FindMenuByRoleId(dbConnection, role.Id, nil)
		//出现错误
		if errFindMenuByRoleId != nil {
			return nil, errFindMenuByRoleId
		}

		//没有菜单
		if len(menusByRoleId) < 1 {
			continue
		}

		//遍历角色的菜单
		for _, menu := range menusByRoleId {
			//menus是否已经包含这个menu,如果包含continue
			_, mok := menusMap[menu.Id]
			if mok {
				continue
			}
			menusMap[menu.Id] = menu
			//设置roleId
			menu.RoleId = role.Id
			//添加菜单
			menus = append(menus, menu)
		}

	}

	//放入缓存
	errPutCache := cache.PutToCache(qxCacheKey, cacheKey, menus)
	if errPutCache != nil {
		return nil, errPutCache
	}

	return menus, nil

}

//UpdateUserRoles 更新用户的角色信息
func UpdateUserRoles(dbConnection *zorm.DBConnection, userId string, roleIds []string) error {

	if len(userId) < 1 {
		return errors.New("userId不能为空")
	}
	//查询用户的现有的角色,清理缓存
	f_select_old := zorm.NewSelectFinder(permstruct.UserRoleStructTableName, "roleId").Append(" WHERE userId=? ", userId)
	listOld := make([]string, 0)
	errQueryList := zorm.QueryStructList(dbConnection, f_select_old, listOld, nil)
	if errQueryList != nil {
		return errQueryList
	}

	//清理老角色缓存
	for _, roleId := range listOld {
		cache.EvictKey(qxCacheKey, "FindUserByRoleId_"+roleId)
	}

	//开启事务,批量保存
	_, errTransaction := zorm.Transaction(dbConnection, func(dbConnection *zorm.DBConnection) (interface{}, error) {

		//删除用户现有的角色
		f_del := zorm.NewDeleteFinder(permstruct.UserRoleStructTableName).Append(" WHERE userId=? ", userId)
		errUpdateFinder := zorm.UpdateFinder(dbConnection, f_del)
		if errUpdateFinder != nil {
			return nil, errUpdateFinder
		}

		if len(roleIds) < 1 {
			return nil, nil
		}
		for _, roleId := range roleIds {
			ur := permstruct.UserRoleStruct{}
			//ur.Id = ""
			ur.UserId = userId
			ur.RoleId = roleId
			errSaveStruct := zorm.SaveStruct(dbConnection, &ur)
			if errSaveStruct != nil {
				return nil, errSaveStruct
			}
		}

		return nil, nil
	})

	//清理用户的缓存
	cache.EvictKey(qxCacheKey, "FindRoleByUserId_"+userId)
	cache.EvictKey(qxCacheKey, "FindMenuByUserId_"+userId)

	//清理新角色缓存
	for _, roleId := range roleIds {
		//清理缓存
		cache.EvictKey(qxCacheKey, "FindUserByRoleId_"+roleId)
	}

	return errTransaction

}
