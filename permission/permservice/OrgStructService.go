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
		comcode, errComcode := findOrgNewComcode(dbConnection, orgStruct.Id, orgStruct.Pid)
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

	//匿名函数return的error如果不为nil,事务就会回滚
	_, errUpdateOrgStruct := zorm.Transaction(dbConnection, func(dbConnection *zorm.DBConnection) (interface{}, error) {

		//事务下的业务代码开始
		errUpdateOrgStruct := zorm.UpdateStruct(dbConnection, orgStruct)

		if errUpdateOrgStruct != nil {
			return nil, errUpdateOrgStruct
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

	return nil
}

//DeleteOrgStructById 根据Id删除部门
//如果入参dbConnection为nil,使用defaultDao开启事务并最后提交.如果入参dbConnection没有事务,调用dbConnection.begin()开启事务并最后提交.如果入参dbConnection有事务,只使用不提交,有开启方提交事务.但是如果遇到错误或者异常,虽然不是事务的开启方,也会回滚事务,让事务尽早回滚
func DeleteOrgStructById(dbConnection *zorm.DBConnection, id string) error {

	//id不能为空
	if len(id) < 1 {
		return errors.New("id不能为空")
	}

	//匿名函数return的error如果不为nil,事务就会回滚
	_, errDeleteOrgStruct := zorm.Transaction(dbConnection, func(dbConnection *zorm.DBConnection) (interface{}, error) {

		//事务下的业务代码开始
		finder := zorm.NewDeleteFinder(permstruct.OrgStructTableName).Append(" WHERE id=?", id)
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
	errCacheOrg := cache.GetFromCache(qxCacheKey, cacheKey, &orgStruct)
	if errCacheOrg != nil {
		return nil, errCacheOrg
	}

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

	cache.PutToCache(qxCacheKey, cacheKey, orgStruct)

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

// findOrgNewComcode 根据id和pid生成部门的Comcode
func findOrgNewComcode(dbConnection *zorm.DBConnection, id string, pid string) (string, error) {

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
