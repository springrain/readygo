package permservice

import (
	"context"
	"errors"
	"fmt"
	"readygo/cache"
	"readygo/permission/permstruct"

	"gitee.com/chunanyong/logger"

	"gitee.com/chunanyong/zorm"
)

//SaveMenuStruct 保存菜单
//如果入参ctx中没有dbConnection,使用defaultDao开启事务并最后提交
//如果入参ctx有dbConnection且没有事务,调用dbConnection.begin()开启事务并最后提交
//如果入参ctx有dbConnection且有事务,只使用不提交,有开启方提交事务
//但是如果遇到错误或者异常,虽然不是事务的开启方,也会回滚事务,让事务尽早回滚
func SaveMenuStruct(ctx context.Context, menuStruct *permstruct.MenuStruct) error {

	//menuStruct对象指针不能为空
	if menuStruct == nil {
		return errors.New("menuStruct对象指针不能为空")
	}

	//匿名函数return的error如果不为nil,事务就会回滚
	_, errSaveMenuStruct := zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {

		//事务下的业务代码开始

		//赋值ID主键
		if len(menuStruct.Id) < 1 {
			menuStruct.Id = zorm.FuncGenerateStringID()
		}

		//获取新的comcode
		comcode, errComcode := newMenuComcode(ctx, menuStruct.Id, menuStruct.Pid)
		if errComcode != nil {
			return nil, errComcode
		}
		menuStruct.Comcode = comcode

		//保存menu
		_, errSaveMenuStruct := zorm.Insert(ctx, menuStruct)

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

	// 清理缓存
	go cache.EvictKey(ctx, baseInfoCacheKey, "findAllMenuTree")

	return nil
}

//UpdateMenuStruct 更新菜单
//如果入参ctx中没有dbConnection,使用defaultDao开启事务并最后提交
//如果入参ctx有dbConnection且没有事务,调用dbConnection.begin()开启事务并最后提交
//如果入参ctx有dbConnection且有事务,只使用不提交,有开启方提交事务
//但是如果遇到错误或者异常,虽然不是事务的开启方,也会回滚事务,让事务尽早回滚
func UpdateMenuStruct(ctx context.Context, menuStruct *permstruct.MenuStruct) error {

	//menuStruct对象指针或主键Id不能为空
	if menuStruct == nil || len(menuStruct.Id) < 1 {
		return errors.New("menuStruct对象指针或主键Id不能为空")
	}

	oldMenu, errById := FindMenuStructById(ctx, menuStruct.Id)
	if errById != nil {
		return errById
	}
	if oldMenu == nil {
		return errors.New("数据库不存在要更新的对象")
	}

	oldComcode := oldMenu.Comcode
	newComcode, errComcode := newMenuComcode(ctx, menuStruct.Id, menuStruct.Pid)

	if errComcode != nil {
		return errComcode
	}

	// 编码改变,级联更新所有的子菜单
	childrenIds, errChildrenIds := FindMenuIdByPid(ctx, menuStruct.Id, nil)
	if errChildrenIds != nil {
		return errChildrenIds
	}

	//匿名函数return的error如果不为nil,事务就会回滚
	_, errUpdateMenuStruct := zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {

		//事务下的业务代码开始

		menuStruct.Comcode = newComcode
		_, errUpdateMenuStruct := zorm.Update(ctx, menuStruct)

		if errUpdateMenuStruct != nil {
			return nil, errUpdateMenuStruct
		}

		if newComcode == oldComcode { // 编码没有改变
			return nil, nil
		}

		// 编码改变,级联更新所有的子菜单

		//没有子菜单
		if len(childrenIds) < 1 {
			return nil, nil
		}

		for _, menuId := range childrenIds {

			if menuId == menuStruct.Id {
				continue
			}

			updateComcode, errComcode := newMenuComcode(ctx, menuId, menuStruct.Id)
			if errComcode != nil {
				return nil, errComcode
			}

			//更新 comCode
			comcodeFinder := zorm.NewUpdateFinder(permstruct.MenuStructTableName).Append(" comcode=? WHERE id=? ", updateComcode, menuId)
			_, errComcodeFinder := zorm.UpdateFinder(ctx, comcodeFinder)
			if errComcodeFinder != nil {
				return nil, errComcodeFinder
			}

		}

		return nil, nil
		//事务下的业务代码结束

	})

	//记录错误
	if errUpdateMenuStruct != nil {
		errUpdateMenuStruct := fmt.Errorf("permservice.UpdateMenuStruct错误:%w", errUpdateMenuStruct)
		logger.Error(errUpdateMenuStruct)
		return errUpdateMenuStruct
	}

	// 清理缓存
	for _, menuId := range childrenIds {
		go cache.EvictKey(ctx, baseInfoCacheKey, "FindMenuStructById_"+menuId)
	}
	//go cache.EvictKey(baseInfoCacheKey, "FindMenuStructById_"+menuStruct.Id)
	go cache.EvictKey(ctx, baseInfoCacheKey, "FindAllMenuTree")

	return nil
}

//DeleteMenuStructById 根据Id删除菜单
//如果入参ctx中没有dbConnection,使用defaultDao开启事务并最后提交
//如果入参ctx有dbConnection且没有事务,调用dbConnection.begin()开启事务并最后提交
//如果入参ctx有dbConnection且有事务,只使用不提交,有开启方提交事务
//但是如果遇到错误或者异常,虽然不是事务的开启方,也会回滚事务,让事务尽早回滚
func DeleteMenuStructById(ctx context.Context, id string) error {

	//id不能为空
	if len(id) < 1 {
		return errors.New("id不能为空")
	}

	//匿名函数return的error如果不为nil,事务就会回滚
	_, errDeleteMenuStruct := zorm.Transaction(ctx, func(ctx context.Context) (interface{}, error) {

		//事务下的业务代码开始

		menuIds, errMenuIds := FindMenuIdByPid(ctx, id, nil)
		if errMenuIds != nil {
			return nil, errMenuIds
		}
		if len(menuIds) < 1 {
			return nil, errors.New("数据库中不存在,id:" + id)
		}

		//删除中间表
		f_delete_re := zorm.NewDeleteFinder(permstruct.RoleMenuStructTableName).Append(" WHERE menuId in (?)", menuIds)
		_, errDeleteRE := zorm.UpdateFinder(ctx, f_delete_re)
		if errDeleteRE != nil {
			return nil, errDeleteRE
		}

		f_delete := zorm.NewDeleteFinder(permstruct.MenuStructTableName).Append(" WHERE id in (?)", menuIds)
		_, errDelete := zorm.UpdateFinder(ctx, f_delete)
		if errDelete != nil {
			return nil, errDelete
		}

		//清理缓存
		for _, menuId := range menuIds {
			go cache.EvictKey(ctx, baseInfoCacheKey, "FindMenuStructById_"+menuId)
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

	// 清理缓存
	go cache.ClearCache(ctx, qxCacheKey)
	go cache.EvictKey(ctx, baseInfoCacheKey, "FindAllMenuTree")

	return nil
}

//FindMenuStructById 根据Id查询菜单信息
//ctx中如果没有dbConnection,则会使用默认的datasource进行无事务查询
func FindMenuStructById(ctx context.Context, id string) (*permstruct.MenuStruct, error) {
	//如果Id为空
	if len(id) < 1 {
		return nil, errors.New("id为空")
	}

	menuStruct := permstruct.MenuStruct{}
	cacheKey := "FindMenuStructById_" + id
	cache.GetFromCache(ctx, baseInfoCacheKey, cacheKey, &menuStruct)
	if len(menuStruct.Id) > 0 { //如果缓存中存在
		return &menuStruct, nil
	}

	//根据Id查询
	finder := zorm.NewSelectFinder(permstruct.MenuStructTableName).Append(" WHERE id=?", id)

	errFindMenuStructById := zorm.QueryRow(ctx, finder, &menuStruct)

	//记录错误
	if errFindMenuStructById != nil {
		errFindMenuStructById := fmt.Errorf("permservice.FindMenuStructById错误:%w", errFindMenuStructById)
		logger.Error(errFindMenuStructById)
		return nil, errFindMenuStructById
	}
	//放入缓存
	cache.PutToCache(ctx, baseInfoCacheKey, cacheKey, menuStruct)
	return &menuStruct, nil

}

//FindMenuStructList 根据Finder查询菜单列表
//ctx中如果没有dbConnection,则会使用默认的datasource进行无事务查询
func FindMenuStructList(ctx context.Context, finder *zorm.Finder, page *zorm.Page) ([]permstruct.MenuStruct, error) {

	//finder不能为空
	if finder == nil {
		return nil, errors.New("finder不能为空")
	}

	menuStructList := make([]permstruct.MenuStruct, 0)
	errFindMenuStructList := zorm.Query(ctx, finder, &menuStructList, page)

	//记录错误
	if errFindMenuStructList != nil {
		errFindMenuStructList := fmt.Errorf("permservice.FindMenuStructList错误:%w", errFindMenuStructList)
		logger.Error(errFindMenuStructList)
		return nil, errFindMenuStructList
	}

	return menuStructList, nil
}

//FindMenuIdByPid 根据pid查询所有的子菜单
func FindMenuIdByPid(ctx context.Context, pid string, page *zorm.Page) ([]string, error) {

	f_select := zorm.NewSelectFinder(permstruct.MenuStructTableName, "id").Append(" WHERE active=1 ")

	if len(pid) > 0 { // pid不是根节点
		menu, errById := FindMenuStructById(ctx, pid)
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
	errQueryList := zorm.Query(ctx, f_select, &menuIds, page)
	if errQueryList != nil {
		return menuIds, errQueryList
	}

	return menuIds, nil
}

//FindAllMenuTree 查询所有的菜单树形结构
func FindAllMenuTree(ctx context.Context) ([]permstruct.MenuStruct, error) {
	cacheKey := "FindAllMenuTree"
	menus := make([]permstruct.MenuStruct, 0)

	//从缓存中取数据
	errFromCache := cache.GetFromCache(ctx, baseInfoCacheKey, cacheKey, &menus)
	if errFromCache != nil {
		return nil, errFromCache
	}
	if len(menus) > 0 { //缓存中有数据
		return menus, nil
	}

	finder := zorm.NewSelectFinder(permstruct.MenuStructTableName).Append(" WHERE active=1 order by sortno desc ")

	errQueryList := zorm.Query(ctx, finder, &menus, nil)
	if errQueryList != nil {
		return nil, errQueryList
	}

	//菜单变成树形结构
	menus = menuList2Tree(menus)

	//放入缓存
	errPutCache := cache.PutToCache(ctx, baseInfoCacheKey, cacheKey, menus)
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

// newMenuComcode 根据id和pid生成菜单的Comcode
func newMenuComcode(ctx context.Context, id string, pid string) (string, error) {

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
	errComcode := zorm.QueryRow(ctx, finder, &comcode)
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
