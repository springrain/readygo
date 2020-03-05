package permservice

import (
	"errors"
	"fmt"
	"readygo/cache"
	"readygo/logger"
	"readygo/permission/permstruct"
	"readygo/zorm"
)

//SaveMenuStruct 保存菜单
//如果入参dbConnection为nil,使用defaultDao开启事务并最后提交.如果入参dbConnection没有事务,调用dbConnection.begin()开启事务并最后提交.如果入参dbConnection有事务,只使用不提交,有开启方提交事务.但是如果遇到错误或者异常,虽然不是事务的开启方,也会回滚事务,让事务尽早回滚
func SaveMenuStruct(dbConnection *zorm.DBConnection, menuStruct *permstruct.MenuStruct) error {

	//menuStruct对象指针不能为空
	if menuStruct == nil {
		return errors.New("menuStruct对象指针不能为空")
	}

	//匿名函数return的error如果不为nil,事务就会回滚
	_, errSaveMenuStruct := zorm.Transaction(dbConnection, func(dbConnection *zorm.DBConnection) (interface{}, error) {

		//事务下的业务代码开始

		//赋值ID主键
		if len(menuStruct.Id) < 1 {
			menuStruct.Id = zorm.GenerateStringID()
		}

		//获取新的comcode
		comcode, errComcode := findMenuNewComcode(dbConnection, menuStruct.Id, menuStruct.Pid)
		if errComcode != nil {
			return nil, errComcode
		}
		menuStruct.Comcode = comcode

		//保存menu
		errSaveMenuStruct := zorm.SaveStruct(dbConnection, menuStruct)

		if errSaveMenuStruct != nil {
			return nil, errSaveMenuStruct
		}

		return nil, nil
		//事务下的业务代码结束

	})

	//记录错误
	if errSaveMenuStruct != nil {
		errSaveMenuStruct := fmt.Errorf("permservice.SaveMenuStruct错误:%w", errSaveMenuStruct)
		logger.Error(errSaveMenuStruct)
		return errSaveMenuStruct
	}

	// 清除缓存
	cache.EvictKey(qxCacheKey, "findAllMenuTree")

	return nil
}

//UpdateMenuStruct 更新菜单
//如果入参dbConnection为nil,使用defaultDao开启事务并最后提交.如果入参dbConnection没有事务,调用dbConnection.begin()开启事务并最后提交.如果入参dbConnection有事务,只使用不提交,有开启方提交事务.但是如果遇到错误或者异常,虽然不是事务的开启方,也会回滚事务,让事务尽早回滚
func UpdateMenuStruct(dbConnection *zorm.DBConnection, menuStruct *permstruct.MenuStruct) error {

	//menuStruct对象指针或主键Id不能为空
	if menuStruct == nil || len(menuStruct.Id) < 1 {
		return errors.New("menuStruct对象指针或主键Id不能为空")
	}

	//匿名函数return的error如果不为nil,事务就会回滚
	_, errUpdateMenuStruct := zorm.Transaction(dbConnection, func(dbConnection *zorm.DBConnection) (interface{}, error) {

		//事务下的业务代码开始

		oldMenu, errById := FindMenuStructById(dbConnection, menuStruct.Id)
		if errById != nil {
			return nil, errById
		}
		if oldMenu == nil {
			return nil, errors.New("数据库不存在要更新的对象")
		}

		oldComcode := oldMenu.Comcode
		newComcode, errComcode := findMenuNewComcode(dbConnection, menuStruct.Id, menuStruct.Pid)
		if errComcode != nil {
			return nil, errComcode
		}
		menuStruct.Comcode = newComcode
		errUpdateMenuStruct := zorm.UpdateStruct(dbConnection, menuStruct)

		if errUpdateMenuStruct != nil {
			return nil, errUpdateMenuStruct
		}

		if newComcode == oldComcode { // 编码没有改变
			return nil, nil
		}

		// 编码改变,级联更新所有的子菜单
		childrenFinder := zorm.NewSelectFinder(permstruct.MenuStructTableName, "id")
		childrenFinder.Append(" WHERE comcode like ? ", oldComcode+"%")

		childrenIds := make([]string, 0)
		errChildrenIds := zorm.QueryStructList(dbConnection, childrenFinder, &childrenIds, nil)
		if errChildrenIds != nil {
			return nil, errChildrenIds
		}

		//没有子菜单
		if len(childrenIds) < 1 {
			return nil, nil
		}

		for _, menuId := range childrenIds {

			if menuId == menuStruct.Id {
				continue
			}

			updateComcode, errComcode := findMenuNewComcode(dbConnection, menuId, menuStruct.Id)
			if errComcode != nil {
				return nil, errComcode
			}
			// 清理缓存
			cache.EvictKey(qxCacheKey, "findMenuById_"+menuId)

			//更新 comCode
			comcodeFinder := zorm.NewUpdateFinder(permstruct.MenuStructTableName).Append(" comcode=? WHERE id=? ", updateComcode, menuId)
			errComcodeFinder := zorm.UpdateFinder(dbConnection, comcodeFinder)
			if errComcodeFinder != nil {
				return nil, errComcodeFinder
			}
		}

		// 清除缓存
		cache.EvictKey(qxCacheKey, "findMenuById_"+menuStruct.Id)
		cache.EvictKey(qxCacheKey, "findAllMenuTree")

		return nil, nil
		//事务下的业务代码结束

	})

	//记录错误
	if errUpdateMenuStruct != nil {
		errUpdateMenuStruct := fmt.Errorf("permservice.UpdateMenuStruct错误:%w", errUpdateMenuStruct)
		logger.Error(errUpdateMenuStruct)
		return errUpdateMenuStruct
	}

	return nil
}

//DeleteMenuStructById 根据Id删除菜单
//如果入参dbConnection为nil,使用defaultDao开启事务并最后提交.如果入参dbConnection没有事务,调用dbConnection.begin()开启事务并最后提交.如果入参dbConnection有事务,只使用不提交,有开启方提交事务.但是如果遇到错误或者异常,虽然不是事务的开启方,也会回滚事务,让事务尽早回滚
func DeleteMenuStructById(dbConnection *zorm.DBConnection, id string) error {

	//id不能为空
	if len(id) < 1 {
		return errors.New("id不能为空")
	}

	//匿名函数return的error如果不为nil,事务就会回滚
	_, errDeleteMenuStruct := zorm.Transaction(dbConnection, func(dbConnection *zorm.DBConnection) (interface{}, error) {

		//事务下的业务代码开始
		finder := zorm.NewDeleteFinder(permstruct.MenuStructTableName).Append(" WHERE id=?", id)
		errDeleteMenuStruct := zorm.UpdateFinder(dbConnection, finder)

		if errDeleteMenuStruct != nil {
			return nil, errDeleteMenuStruct
		}

		return nil, nil
		//事务下的业务代码结束

	})

	//记录错误
	if errDeleteMenuStruct != nil {
		errDeleteMenuStruct := fmt.Errorf("permservice.DeleteMenuStruct错误:%w", errDeleteMenuStruct)
		logger.Error(errDeleteMenuStruct)
		return errDeleteMenuStruct
	}

	return nil
}

//FindMenuStructById 根据Id查询菜单信息
//dbConnection如果为nil,则会使用默认的datasource进行无事务查询
func FindMenuStructById(dbConnection *zorm.DBConnection, id string) (*permstruct.MenuStruct, error) {
	//如果Id为空
	if len(id) < 1 {
		return nil, errors.New("id为空")
	}

	//根据Id查询
	finder := zorm.NewSelectFinder(permstruct.MenuStructTableName).Append(" WHERE id=?", id)
	menuStruct := permstruct.MenuStruct{}
	errFindMenuStructById := zorm.QueryStruct(dbConnection, finder, &menuStruct)

	//记录错误
	if errFindMenuStructById != nil {
		errFindMenuStructById := fmt.Errorf("permservice.FindMenuStructById错误:%w", errFindMenuStructById)
		logger.Error(errFindMenuStructById)
		return nil, errFindMenuStructById
	}

	return &menuStruct, nil

}

//FindMenuStructList 根据Finder查询菜单列表
//dbConnection如果为nil,则会使用默认的datasource进行无事务查询
func FindMenuStructList(dbConnection *zorm.DBConnection, finder *zorm.Finder, page *zorm.Page) ([]permstruct.MenuStruct, error) {

	//finder不能为空
	if finder == nil {
		return nil, errors.New("finder不能为空")
	}

	menuStructList := make([]permstruct.MenuStruct, 0)
	errFindMenuStructList := zorm.QueryStructList(dbConnection, finder, &menuStructList, page)

	//记录错误
	if errFindMenuStructList != nil {
		errFindMenuStructList := fmt.Errorf("permservice.FindMenuStructList错误:%w", errFindMenuStructList)
		logger.Error(errFindMenuStructList)
		return nil, errFindMenuStructList
	}

	return menuStructList, nil
}

//FindMenuByPid 根据pid查询所有的子菜单
func FindMenuByPid(dbConnection *zorm.DBConnection, pid string, page *zorm.Page) ([]string, error) {

	f_select := zorm.NewSelectFinder(permstruct.AliPayconfigStructTableName, "id").Append(" WHERE active=1 ")

	if len(pid) > 0 { // pid不是根节点
		menu, errById := FindMenuStructById(dbConnection, pid)
		if errById != nil {
			return nil, errById
		}

		if menu.Comcode == "" { //没有编码,错误数据
			return nil, errors.New("Comcode为空,错误数据,pid:" + pid)
		}
		f_select.Append(" and comcode like ? ", menu.Comcode+"%")
	}

	f_select.Append(" order by sortno desc ")
	menuIds := make([]string, 0)
	errQueryList := zorm.QueryStructList(dbConnection, f_select, &menuIds, page)
	if errQueryList != nil {
		return menuIds, errQueryList
	}

	return menuIds, nil
}

//FindAllMenuTree 查询所有的菜单树形结构
func FindAllMenuTree(dbConnection *zorm.DBConnection) ([]permstruct.MenuStruct, error) {
	cacheKey := "FindAllMenuTree"
	menus := make([]permstruct.MenuStruct, 0)

	//从缓存中取数据
	errFromCache := cache.GetFromCache(qxCacheKey, cacheKey, &menus)
	if errFromCache != nil {
		return nil, errFromCache
	}
	if len(menus) > 0 { //缓存中有数据
		return menus, nil
	}

	finder := zorm.NewSelectFinder(permstruct.MenuStructTableName).Append(" WHERE active=1 order by sortno desc ")

	errQueryList := zorm.QueryStructList(dbConnection, finder, menus, nil)
	if errQueryList != nil {
		return nil, errQueryList
	}

	//菜单变成树形结构
	menus = menuList2Tree(menus)

	//放入缓存
	errPutCache := cache.PutToCache(qxCacheKey, cacheKey, menus)
	if errPutCache != nil {
		return nil, errPutCache
	}

	return menus, nil
}

// 将平行的List,变成树形结构
func menuList2Tree(menuList []permstruct.MenuStruct) []permstruct.MenuStruct {

	if len(menuList) < 1 {
		return menuList
	}
	// 先把数据放到map里,方便取值
	menuMap := make(map[string]permstruct.MenuStruct)

	//map赋值
	for _, menu := range menuList {
		menuMap[menu.Id] = menu
	}
	// 循环遍历menuList
	list := make([]permstruct.MenuStruct, 0)
	for _, menu := range menuList {
		pid := menu.Pid
		parent, pidOk := menuMap[pid]
		// 没有父节点
		if !pidOk {
			list = append(list, menu)
			continue
		}

		//如果有父节点
		children := parent.Children
		if children == nil {
			children = make([]permstruct.MenuStruct, 0)
			parent.Children = children
		}
		children = append(children, menu)
	}

	return list
}

/**
  @Override
  public void wrapVueMenu(List<Menu> listMenu, List<Map<String, Object>> listMap) {
      if (CollectionUtils.isEmpty(listMenu)) {
          return;
      }
      for (Menu menu : listMenu) {
          Map<String, Object> map = new HashMap<>();
          listMap.add(map);
          map.put("path", menu.getPath());
          map.put("redirect", menu.getRedirect());
          map.put("component", menu.getComponent());
          map.put("name", menu.getName());
          map.put("menuType", menu.getMenuType());
          map.put("roleId", menu.getRoleId());
          map.put("menuId", menu.getId());
          map.put("pageurl",menu.getPageurl());

          // meta
          Map<String, Object> meta = new HashMap<>();
          map.put("meta", meta);
          meta.put("title", menu.getTitle());
          meta.put("permission", menu.getPermission());
          meta.put("keepAlive", menu.getKeepAlive());
          //meta.put("target",menu.getTarget());


          List<Menu> listChildren = menu.getChildren();
          if (CollectionUtils.isNotEmpty(listChildren)) {
              // children
              List<Map<String, Object>> children = new ArrayList<>();
              map.put("children", children);
              wrapVueMenu(listChildren, children);
          }

      }

  }
*/

// findMenuNewComcode 根据id和pid生成菜单的Comcode
func findMenuNewComcode(dbConnection *zorm.DBConnection, id string, pid string) (string, error) {

	//id不能为空
	if len(id) < 1 {
		return "", errors.New("id不能为空")
	}

	//没有上级
	if len(pid) < 1 {
		return "," + id + ",", nil
	}

	comcode := ""
	finder := zorm.NewSelectFinder(permstruct.MenuStructTableName, "comcode").Append(" WHERE id=? ", pid)
	errComcode := zorm.QueryStruct(dbConnection, finder, &comcode)
	if errComcode != nil {
		return "", errComcode
	}

	//没有上级
	if len(comcode) < 1 {
		return "," + id + ",", nil
	}

	comcode = comcode + id + ","

	return comcode, nil
}
