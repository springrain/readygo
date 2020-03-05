package permservice

import (
	"errors"
	"fmt"
	"readygo/logger"
	"readygo/permission/permstruct"
	"readygo/zorm"
)

//SaveRoleStruct 保存角色
//如果入参dbConnection为nil,使用defaultDao开启事务并最后提交.如果入参dbConnection没有事务,调用dbConnection.begin()开启事务并最后提交.如果入参dbConnection有事务,只使用不提交,有开启方提交事务.但是如果遇到错误或者异常,虽然不是事务的开启方,也会回滚事务,让事务尽早回滚
func SaveRoleStruct(dbConnection *zorm.DBConnection, roleStruct *permstruct.RoleStruct) error {

	// roleStruct对象指针不能为空
	if roleStruct == nil {
		return errors.New("roleStruct对象指针不能为空")
	}
	//匿名函数return的error如果不为nil,事务就会回滚
	_, errSaveRoleStruct := zorm.Transaction(dbConnection, func(dbConnection *zorm.DBConnection) (interface{}, error) {

		//事务下的业务代码开始

		//赋值主键Id
		if len(roleStruct.Id) < 1 {
			roleStruct.Id = zorm.GenerateStringID()
		}

		errSaveRoleStruct := zorm.SaveStruct(dbConnection, roleStruct)

		if errSaveRoleStruct != nil {
			return nil, errSaveRoleStruct
		}

		return nil, nil
		//事务下的业务代码结束

	})

	//记录错误
	if errSaveRoleStruct != nil {
		errSaveRoleStruct := fmt.Errorf("permservice.SaveRoleStruct错误:%w", errSaveRoleStruct)
		logger.Error(errSaveRoleStruct)
		return errSaveRoleStruct
	}

	return nil
}

//UpdateRoleStruct 更新角色
//如果入参dbConnection为nil,使用defaultDao开启事务并最后提交.如果入参dbConnection没有事务,调用dbConnection.begin()开启事务并最后提交.如果入参dbConnection有事务,只使用不提交,有开启方提交事务.但是如果遇到错误或者异常,虽然不是事务的开启方,也会回滚事务,让事务尽早回滚
func UpdateRoleStruct(dbConnection *zorm.DBConnection, roleStruct *permstruct.RoleStruct) error {

	// roleStruct对象指针或主键Id不能为空
	if roleStruct == nil || len(roleStruct.Id) < 1 {
		return errors.New("roleStruct对象指针或主键Id不能为空")
	}

	//匿名函数return的error如果不为nil,事务就会回滚
	_, errUpdateRoleStruct := zorm.Transaction(dbConnection, func(dbConnection *zorm.DBConnection) (interface{}, error) {

		//事务下的业务代码开始
		errUpdateRoleStruct := zorm.UpdateStruct(dbConnection, roleStruct)

		if errUpdateRoleStruct != nil {
			return nil, errUpdateRoleStruct
		}

		return nil, nil
		//事务下的业务代码结束

	})

	//记录错误
	if errUpdateRoleStruct != nil {
		errUpdateRoleStruct := fmt.Errorf("permservice.UpdateRoleStruct错误:%w", errUpdateRoleStruct)
		logger.Error(errUpdateRoleStruct)
		return errUpdateRoleStruct
	}

	return nil
}

//DeleteRoleStructById 根据Id删除角色
//如果入参dbConnection为nil,使用defaultDao开启事务并最后提交.如果入参dbConnection没有事务,调用dbConnection.begin()开启事务并最后提交.如果入参dbConnection有事务,只使用不提交,有开启方提交事务.但是如果遇到错误或者异常,虽然不是事务的开启方,也会回滚事务,让事务尽早回滚
func DeleteRoleStructById(dbConnection *zorm.DBConnection, id string) error {

	//id不能为空
	if len(id) < 1 {
		return errors.New("id不能为空")
	}

	//匿名函数return的error如果不为nil,事务就会回滚
	_, errDeleteRoleStruct := zorm.Transaction(dbConnection, func(dbConnection *zorm.DBConnection) (interface{}, error) {

		//事务下的业务代码开始
		finder := zorm.NewDeleteFinder(permstruct.RoleStructTableName).Append(" WHERE id=?", id)
		errDeleteRoleStruct := zorm.UpdateFinder(dbConnection, finder)

		if errDeleteRoleStruct != nil {
			return nil, errDeleteRoleStruct
		}

		return nil, nil
		//事务下的业务代码结束

	})

	//记录错误
	if errDeleteRoleStruct != nil {
		errDeleteRoleStruct := fmt.Errorf("permservice.DeleteRoleStruct错误:%w", errDeleteRoleStruct)
		logger.Error(errDeleteRoleStruct)
		return errDeleteRoleStruct
	}

	return nil
}

//FindRoleStructById 根据Id查询角色信息
//dbConnection如果为nil,则会使用默认的datasource进行无事务查询
func FindRoleStructById(dbConnection *zorm.DBConnection, id string) (*permstruct.RoleStruct, error) {
	//id不能为空
	if len(id) < 1 {
		return nil, errors.New("id不能为空")
	}

	//根据Id查询
	finder := zorm.NewSelectFinder(permstruct.RoleStructTableName).Append(" WHERE id=?", id)
	roleStruct := permstruct.RoleStruct{}
	errFindRoleStructById := zorm.QueryStruct(dbConnection, finder, &roleStruct)

	//记录错误
	if errFindRoleStructById != nil {
		errFindRoleStructById := fmt.Errorf("permservice.FindRoleStructById错误:%w", errFindRoleStructById)
		logger.Error(errFindRoleStructById)
		return nil, errFindRoleStructById
	}

	return &roleStruct, nil

}

//FindRoleStructList 根据Finder查询角色列表
//dbConnection如果为nil,则会使用默认的datasource进行无事务查询
func FindRoleStructList(dbConnection *zorm.DBConnection, finder *zorm.Finder, page *zorm.Page) ([]permstruct.RoleStruct, error) {

	//finder不能为空
	if finder == nil {
		return nil, errors.New("finder不能为空")
	}

	roleStructList := make([]permstruct.RoleStruct, 0)
	errFindRoleStructList := zorm.QueryStructList(dbConnection, finder, &roleStructList, page)

	//记录错误
	if errFindRoleStructList != nil {
		errFindRoleStructList := fmt.Errorf("permservice.FindRoleStructList错误:%w", errFindRoleStructList)
		logger.Error(errFindRoleStructList)
		return nil, errFindRoleStructList
	}

	return roleStructList, nil
}
