package permservice

import (
	"errors"
	"fmt"
	"readygo/logger"
	"readygo/permission/permstruct"
	"readygo/zorm"
)

//SaveOrgStruct 保存部门
//如果入参session为nil,使用defaultDao开启事务并最后提交.如果入参session没有事务,调用session.begin()开启事务并最后提交.如果入参session有事务,只使用不提交,有开启方提交事务.但是如果遇到错误或者异常,虽然不是事务的开启方,也会回滚事务,让事务尽早回滚
func SaveOrgStruct(session *zorm.Session, orgStruct *permstruct.OrgStruct) error {

	//匿名函数return的error如果不为nil,事务就会回滚
	_, errSaveOrgStruct := zorm.Transaction(session, func(session *zorm.Session) (interface{}, error) {

		//事务下的业务代码开始
		errSaveOrgStruct := zorm.SaveStruct(session, orgStruct)

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
//如果入参session为nil,使用defaultDao开启事务并最后提交.如果入参session没有事务,调用session.begin()开启事务并最后提交.如果入参session有事务,只使用不提交,有开启方提交事务.但是如果遇到错误或者异常,虽然不是事务的开启方,也会回滚事务,让事务尽早回滚
func UpdateOrgStruct(session *zorm.Session, orgStruct *permstruct.OrgStruct) error {

	//匿名函数return的error如果不为nil,事务就会回滚
	_, errUpdateOrgStruct := zorm.Transaction(session, func(session *zorm.Session) (interface{}, error) {

		//事务下的业务代码开始
		errUpdateOrgStruct := zorm.UpdateStruct(session, orgStruct)

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

//DeleteOrgStruct 删除部门
//如果入参session为nil,使用defaultDao开启事务并最后提交.如果入参session没有事务,调用session.begin()开启事务并最后提交.如果入参session有事务,只使用不提交,有开启方提交事务.但是如果遇到错误或者异常,虽然不是事务的开启方,也会回滚事务,让事务尽早回滚
func DeleteOrgStruct(session *zorm.Session, orgStruct *permstruct.OrgStruct) error {

	//匿名函数return的error如果不为nil,事务就会回滚
	_, errDeleteOrgStruct := zorm.Transaction(session, func(session *zorm.Session) (interface{}, error) {

		//事务下的业务代码开始
		errDeleteOrgStruct := zorm.DeleteStruct(session, orgStruct)

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
//session如果为nil,则会使用默认的datasource进行无事务查询
func FindOrgStructById(session *zorm.Session, id string) (*permstruct.OrgStruct, error) {
	//如果Id为空
	if len(id) < 1 {
		return nil, errors.New("id为空")
	}

	//根据Id查询
	finder := zorm.NewSelectFinder(" WHERE id=?", id)
	orgStruct := permstruct.OrgStruct{}
	errFindOrgStructById := zorm.QueryStruct(session, finder, &orgStruct)

	//记录错误
	if errFindOrgStructById != nil {
		errFindOrgStructById := fmt.Errorf("permservice.FindOrgStructById错误:%w", errFindOrgStructById)
		logger.Error(errFindOrgStructById)
		return nil, errFindOrgStructById
	}

	return &orgStruct, nil

}

//FindOrgStructList 根据Finder查询部门列表
//session如果为nil,则会使用默认的datasource进行无事务查询
func FindOrgStructList(session *zorm.Session, finder *zorm.Finder, page *zorm.Page) ([]permstruct.OrgStruct, error) {
	orgStructList := make([]permstruct.OrgStruct, 0)
	errFindOrgStructList := zorm.QueryStructList(session, finder, &orgStructList, page)

	//记录错误
	if errFindOrgStructList != nil {
		errFindOrgStructList := fmt.Errorf("permservice.FindOrgStructList错误:%w", errFindOrgStructList)
		logger.Error(errFindOrgStructList)
		return nil, errFindOrgStructList
	}

	return orgStructList, nil
}