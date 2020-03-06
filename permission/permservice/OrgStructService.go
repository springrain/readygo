package permservice

import (
	"errors"
	"fmt"
	"readygo/cache"
	"readygo/logger"
	"readygo/permission/permstruct"
	"readygo/zorm"
)

//SaveOrgStruct 保存部门
//如果入参dbConnection为nil,使用defaultDao开启事务并最后提交.如果入参dbConnection没有事务,调用dbConnection.begin()开启事务并最后提交.如果入参dbConnection有事务,只使用不提交,有开启方提交事务.但是如果遇到错误或者异常,虽然不是事务的开启方,也会回滚事务,让事务尽早回滚
func SaveOrgStruct(dbConnection *zorm.DBConnection, orgStruct *permstruct.OrgStruct) error {

	// orgStruct对象指针不能为空
	if orgStruct == nil {
		return errors.New("orgStruct对象指针不能为空")
	}
	//匿名函数return的error如果不为nil,事务就会回滚
	_, errSaveOrgStruct := zorm.Transaction(dbConnection, func(dbConnection *zorm.DBConnection) (interface{}, error) {

		//事务下的业务代码开始

		//赋值主键Id
		if len(orgStruct.Id) < 1 {
			orgStruct.Id = zorm.GenerateStringID()
		}

		//获取新的comcode
		comcode, errComcode := newOrgComcode(dbConnection, orgStruct.Id, orgStruct.Pid)
		if errComcode != nil {
			return nil, errComcode
		}
		orgStruct.Comcode = comcode
		orgStruct.Active = 1

		errSaveOrgStruct := zorm.SaveStruct(dbConnection, orgStruct)

		if errSaveOrgStruct != nil {
			return nil, errSaveOrgStruct
		}

		return nil, nil
		//事务下的业务代码结束

	})

	//记录错误
	if errSaveOrgStruct != nil {
		errSaveOrgStruct := fmt.Errorf("permservice.SaveOrgStruct错误:%w", errSaveOrgStruct)
		logger.Error(errSaveOrgStruct)
		return errSaveOrgStruct
	}

	return nil
}

//UpdateOrgStruct 更新部门
//如果入参dbConnection为nil,使用defaultDao开启事务并最后提交.如果入参dbConnection没有事务,调用dbConnection.begin()开启事务并最后提交.如果入参dbConnection有事务,只使用不提交,有开启方提交事务.但是如果遇到错误或者异常,虽然不是事务的开启方,也会回滚事务,让事务尽早回滚
func UpdateOrgStruct(dbConnection *zorm.DBConnection, orgStruct *permstruct.OrgStruct) error {

	// orgStruct对象指针或主键Id不能为空
	if orgStruct == nil || len(orgStruct.Id) < 1 {
		return errors.New("orgStruct对象指针或主键Id不能为空")
	}

	oldOrg, errById := FindOrgStructById(dbConnection, orgStruct.Id)
	if errById != nil {
		return errById
	}
	if oldOrg == nil {
		return errors.New("数据库不存在要更新的对象")
	}

	oldComcode := oldOrg.Comcode
	newComcode, errComcode := newOrgComcode(dbConnection, orgStruct.Id, orgStruct.Pid)
	if errComcode != nil {
		return errComcode
	}

	childrenIds, errChildrenIds := FindOrgIdByPid(dbConnection, orgStruct.Id)
	if errChildrenIds != nil {
		return errChildrenIds
	}

	//匿名函数return的error如果不为nil,事务就会回滚
	_, errUpdateOrgStruct := zorm.Transaction(dbConnection, func(dbConnection *zorm.DBConnection) (interface{}, error) {

		orgStruct.Comcode = newComcode
		errUpdateOrgStruct := zorm.UpdateStruct(dbConnection, orgStruct)

		if errUpdateOrgStruct != nil {
			return nil, errUpdateOrgStruct
		}

		if newComcode == oldComcode { // 编码没有改变
			return nil, nil
		}

		// 编码改变,级联更新所有的子部门
		//没有子部门
		if len(childrenIds) < 1 {
			return nil, nil
		}
		for _, orgId := range childrenIds {

			if orgId == orgStruct.Id {
				continue
			}

			updateComcode, errComcode := newOrgComcode(dbConnection, orgId, orgStruct.Id)
			if errComcode != nil {
				return nil, errComcode
			}

			//更新 comCode
			comcodeFinder := zorm.NewUpdateFinder(permstruct.OrgStructTableName).Append(" comcode=? WHERE id=? ", updateComcode, orgId)
			errComcodeFinder := zorm.UpdateFinder(dbConnection, comcodeFinder)
			if errComcodeFinder != nil {
				return nil, errComcodeFinder
			}

		}

		return nil, nil
		//事务下的业务代码结束

	})

	//记录错误
	if errUpdateOrgStruct != nil {
		errUpdateOrgStruct := fmt.Errorf("permservice.UpdateOrgStruct错误:%w", errUpdateOrgStruct)
		logger.Error(errUpdateOrgStruct)
		return errUpdateOrgStruct
	}
	// 清除缓存
	for _, orgId := range childrenIds {
		go cache.EvictKey(baseInfoCacheKey, "FindOrgStructById_"+orgId)
	}
	//go cache.EvictKey(baseInfoCacheKey, "FindOrgStructById_"+orgStruct.Id)
	return nil
}

//DeleteOrgStructById 根据Id删除部门
//如果入参dbConnection为nil,使用defaultDao开启事务并最后提交.如果入参dbConnection没有事务,调用dbConnection.begin()开启事务并最后提交.如果入参dbConnection有事务,只使用不提交,有开启方提交事务.但是如果遇到错误或者异常,虽然不是事务的开启方,也会回滚事务,让事务尽早回滚
func DeleteOrgStructById(dbConnection *zorm.DBConnection, id string) error {

	//id不能为空
	if len(id) < 1 {
		return errors.New("id不能为空")
	}

	org, errById := FindOrgStructById(dbConnection, id)
	if errById != nil {
		return errById
	}
	if org == nil {
		return errors.New("数据库不存在要删除的对象")
	}

	orgIds, errByPid := FindOrgIdByPid(dbConnection, id)
	if errByPid != nil {
		return errByPid
	}

	//匿名函数return的error如果不为nil,事务就会回滚
	_, errDeleteOrgStruct := zorm.Transaction(dbConnection, func(dbConnection *zorm.DBConnection) (interface{}, error) {

		//事务下的业务代码开始

		finder := zorm.NewUpdateFinder(permstruct.OrgStructTableName).Append("  active=0  WHERE comcode like ? ", org.Comcode+"%")
		errDeleteOrgStruct := zorm.UpdateFinder(dbConnection, finder)

		if errDeleteOrgStruct != nil {
			return nil, errDeleteOrgStruct
		}

		return nil, nil
		//事务下的业务代码结束

	})

	//记录错误
	if errDeleteOrgStruct != nil {
		errDeleteOrgStruct := fmt.Errorf("permservice.DeleteOrgStruct错误:%w", errDeleteOrgStruct)
		logger.Error(errDeleteOrgStruct)
		return errDeleteOrgStruct
	}
	//清理缓存
	for _, orgId := range orgIds {
		//清理缓存
		go cache.EvictKey(baseInfoCacheKey, "FindOrgStructById_"+orgId)
	}
	//go cache.EvictKey(baseInfoCacheKey, "FindOrgStructById_"+id)
	go cache.ClearCache(qxCacheKey)
	return nil
}

//FindOrgStructById 根据Id查询部门信息
//dbConnection如果为nil,则会使用默认的datasource进行无事务查询
func FindOrgStructById(dbConnection *zorm.DBConnection, id string) (*permstruct.OrgStruct, error) {
	//id不能为空
	if len(id) < 1 {
		return nil, errors.New("id不能为空")
	}
	orgStruct := permstruct.OrgStruct{}
	cacheKey := "FindOrgStructById_" + id
	cache.GetFromCache(baseInfoCacheKey, cacheKey, &orgStruct)
	if len(orgStruct.Id) > 0 { //缓存中有值
		return &orgStruct, nil
	}

	//根据Id查询
	finder := zorm.NewSelectFinder(permstruct.OrgStructTableName).Append(" WHERE id=?", id)

	errFindOrgStructById := zorm.QueryStruct(dbConnection, finder, &orgStruct)

	//记录错误
	if errFindOrgStructById != nil {
		errFindOrgStructById := fmt.Errorf("permservice.FindOrgStructById错误:%w", errFindOrgStructById)
		logger.Error(errFindOrgStructById)
		return nil, errFindOrgStructById
	}

	//放入缓存
	cache.PutToCache(baseInfoCacheKey, cacheKey, orgStruct)

	return &orgStruct, nil

}

//FindOrgStructList 根据Finder查询部门列表
//dbConnection如果为nil,则会使用默认的datasource进行无事务查询
func FindOrgStructList(dbConnection *zorm.DBConnection, finder *zorm.Finder, page *zorm.Page) ([]permstruct.OrgStruct, error) {

	//finder不能为空
	if finder == nil {
		return nil, errors.New("finder不能为空")
	}

	orgStructList := make([]permstruct.OrgStruct, 0)
	errFindOrgStructList := zorm.QueryStructList(dbConnection, finder, &orgStructList, page)

	//记录错误
	if errFindOrgStructList != nil {
		errFindOrgStructList := fmt.Errorf("permservice.FindOrgStructList错误:%w", errFindOrgStructList)
		logger.Error(errFindOrgStructList)
		return nil, errFindOrgStructList
	}

	return orgStructList, nil
}

//FindOrgTreeByPid 根据pid查询组织树形的组织结构
func FindOrgTreeByPid(dbConnection *zorm.DBConnection, pid string) ([]permstruct.OrgStruct, error) {

	finder := zorm.NewSelectFinder(permstruct.OrgStructTableName).Append("WHERE active=1 ")
	if len(pid) > 0 { //不是根目录
		org, errById := FindOrgStructById(dbConnection, pid)
		if errById != nil {
			return nil, errById
		}
		if org == nil {
			return nil, errors.New("数据库不存在对象,id:" + pid)
		}

		finder.Append(" and comcode like ? ", org.Comcode)

	}
	finder.Append(" order by sortno asc ")

	orgs := make([]permstruct.OrgStruct, 0)
	errQueryList := zorm.QueryStructList(dbConnection, finder, &orgs, nil)
	if errQueryList != nil {
		return nil, errQueryList
	}

	//菜单变成树形结构
	orgs = orgList2Tree(orgs)

	return orgs, nil

}

//FindOrgIdByPid 根据pid查询子部门的Id
func FindOrgIdByPid(dbConnection *zorm.DBConnection, pid string) ([]string, error) {

	finder := zorm.NewSelectFinder(permstruct.OrgStructTableName, "id").Append("WHERE active=1 ")
	if len(pid) > 0 { //不是根目录
		org, errById := FindOrgStructById(dbConnection, pid)
		if errById != nil {
			return nil, errById
		}
		if org == nil {
			return nil, errors.New("数据库不存在对象,id:" + pid)
		}

		finder.Append(" and comcode like ? ", org.Comcode)

	}
	finder.Append(" order by sortno asc ")

	orgIds := make([]string, 0)
	errQueryList := zorm.QueryStructList(dbConnection, finder, &orgIds, nil)
	if errQueryList != nil {
		return nil, errQueryList
	}
	return orgIds, nil

}

// UpdateOrgManagerUserId 更新部门主管
func UpdateOrgManagerUserId(dbConnection *zorm.DBConnection, orgId string, managerUserId string) error {

	if len(orgId) < 1 || len(managerUserId) < 1 {
		return errors.New("orgId或者managerUserId不能为空")
	}
	//匿名函数return的error如果不为nil,事务就会回滚
	_, errUpdateOrgManagerUserId := zorm.Transaction(dbConnection, func(dbConnection *zorm.DBConnection) (interface{}, error) {
		finder := zorm.NewDeleteFinder(permstruct.UserOrgStructTableName).Append(" WHERE orgId=? and managerType=2 ", orgId)
		errUpdateFinder := zorm.UpdateFinder(dbConnection, finder)
		if errUpdateFinder != nil {
			return nil, errUpdateFinder
		}
		userOrg := permstruct.UserOrgStruct{}
		userOrg.Id = zorm.GenerateStringID()
		userOrg.OrgId = orgId
		userOrg.UserId = managerUserId
		userOrg.ManagerType = 2

		errSave := zorm.SaveStruct(dbConnection, &userOrg)
		if errSave != nil {
			return nil, errSave
		}
		return nil, nil
	})

	//记录错误
	if errUpdateOrgManagerUserId != nil {
		errUpdateOrgStruct := fmt.Errorf("permservice.DeleteOrgStruct错误:%w", errUpdateOrgManagerUserId)
		logger.Error(errUpdateOrgStruct)
		return errUpdateOrgStruct
	}

	return nil

}

// 将平行的List,变成树形结构
func orgList2Tree(orgList []permstruct.OrgStruct) []permstruct.OrgStruct {

	if len(orgList) < 1 {
		return orgList
	}
	// 先把数据放到map里,方便取值
	orgMap := make(map[string]permstruct.OrgStruct)

	//map赋值
	for _, org := range orgList {
		orgMap[org.Id] = org
	}
	// 循环遍历OrgList
	list := make([]permstruct.OrgStruct, 0)
	for _, org := range orgList {
		pid := org.Pid
		parent, pidOk := orgMap[pid]
		// 没有父节点
		if !pidOk {
			list = append(list, org)
			continue
		}

		//如果有父节点
		children := parent.Children
		if children == nil {
			children = make([]permstruct.OrgStruct, 0)
			parent.Children = children
		}
		children = append(children, org)
	}

	return list
}

// newOrgComcode 根据id和pid生成部门的Comcode
func newOrgComcode(dbConnection *zorm.DBConnection, id string, pid string) (string, error) {

	//id不能为空
	if len(id) < 1 {
		return "", errors.New("id不能为空")
	}

	//没有上级
	if len(pid) < 1 {
		return "," + id + ",", nil
	}

	comcode := ""
	finder := zorm.NewSelectFinder(permstruct.OrgStructTableName, "comcode").Append(" WHERE id=? ", pid)
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
