package permservice

import (
	"errors"
	"fmt"
	"readygo/logger"
	"readygo/permission/permstruct"
	"readygo/zorm"
)

//SaveDicDataStruct 保存公共字典
//如果入参session为nil,使用defaultDao开启事务并最后提交.如果入参session没有事务,调用session.begin()开启事务并最后提交.如果入参session有事务,只使用不提交,有开启方提交事务.但是如果遇到错误或者异常,虽然不是事务的开启方,也会回滚事务,让事务尽早回滚
func SaveDicDataStruct(session *zorm.Session, dicDataStruct *permstruct.DicDataStruct) error {

	//匿名函数return的error如果不为nil,事务就会回滚
	_, errSaveDicDataStruct := zorm.Transaction(session, func(session *zorm.Session) (interface{}, error) {

		//事务下的业务代码开始
		errSaveDicDataStruct := zorm.SaveStruct(session, dicDataStruct)

		if errSaveDicDataStruct != nil {
			return nil, errSaveDicDataStruct
		}

		return nil, nil
		//事务下的业务代码结束

	})

	//记录错误
	if errSaveDicDataStruct != nil {
		errSaveDicDataStruct := fmt.Errorf("permservice.SaveDicDataStruct错误:%w", errSaveDicDataStruct)
		logger.Error(errSaveDicDataStruct)
		return errSaveDicDataStruct
	}

	return nil
}

//UpdateDicDataStruct 更新公共字典
//如果入参session为nil,使用defaultDao开启事务并最后提交.如果入参session没有事务,调用session.begin()开启事务并最后提交.如果入参session有事务,只使用不提交,有开启方提交事务.但是如果遇到错误或者异常,虽然不是事务的开启方,也会回滚事务,让事务尽早回滚
func UpdateDicDataStruct(session *zorm.Session, dicDataStruct *permstruct.DicDataStruct) error {

	//匿名函数return的error如果不为nil,事务就会回滚
	_, errUpdateDicDataStruct := zorm.Transaction(session, func(session *zorm.Session) (interface{}, error) {

		//事务下的业务代码开始
		errUpdateDicDataStruct := zorm.UpdateStruct(session, dicDataStruct)

		if errUpdateDicDataStruct != nil {
			return nil, errUpdateDicDataStruct
		}

		return nil, nil
		//事务下的业务代码结束

	})

	//记录错误
	if errUpdateDicDataStruct != nil {
		errUpdateDicDataStruct := fmt.Errorf("permservice.UpdateDicDataStruct错误:%w", errUpdateDicDataStruct)
		logger.Error(errUpdateDicDataStruct)
		return errUpdateDicDataStruct
	}

	return nil
}

//DeleteDicDataStruct 删除公共字典
//如果入参session为nil,使用defaultDao开启事务并最后提交.如果入参session没有事务,调用session.begin()开启事务并最后提交.如果入参session有事务,只使用不提交,有开启方提交事务.但是如果遇到错误或者异常,虽然不是事务的开启方,也会回滚事务,让事务尽早回滚
func DeleteDicDataStruct(session *zorm.Session, dicDataStruct *permstruct.DicDataStruct) error {

	//匿名函数return的error如果不为nil,事务就会回滚
	_, errDeleteDicDataStruct := zorm.Transaction(session, func(session *zorm.Session) (interface{}, error) {

		//事务下的业务代码开始
		errDeleteDicDataStruct := zorm.DeleteStruct(session, dicDataStruct)

		if errDeleteDicDataStruct != nil {
			return nil, errDeleteDicDataStruct
		}

		return nil, nil
		//事务下的业务代码结束

	})

	//记录错误
	if errDeleteDicDataStruct != nil {
		errDeleteDicDataStruct := fmt.Errorf("permservice.DeleteDicDataStruct错误:%w", errDeleteDicDataStruct)
		logger.Error(errDeleteDicDataStruct)
		return errDeleteDicDataStruct
	}

	return nil
}

//FindDicDataStructById 根据Id查询公共字典信息
//session如果为nil,则会使用默认的datasource进行无事务查询
func FindDicDataStructById(session *zorm.Session, id string) (*permstruct.DicDataStruct, error) {
	//如果Id为空
	if len(id) < 1 {
		return nil, errors.New("id为空")
	}

	//根据Id查询
	finder := zorm.NewSelectFinder(" WHERE id=?", id)
	dicDataStruct := permstruct.DicDataStruct{}
	errFindDicDataStructById := zorm.QueryStruct(session, finder, &dicDataStruct)

	//记录错误
	if errFindDicDataStructById != nil {
		errFindDicDataStructById := fmt.Errorf("permservice.FindDicDataStructById错误:%w", errFindDicDataStructById)
		logger.Error(errFindDicDataStructById)
		return nil, errFindDicDataStructById
	}

	return &dicDataStruct, nil

}

//FindDicDataStructList 根据Finder查询公共字典列表
//session如果为nil,则会使用默认的datasource进行无事务查询
func FindDicDataStructList(session *zorm.Session, finder *zorm.Finder, page *zorm.Page) ([]permstruct.DicDataStruct, error) {
	dicDataStructList := make([]permstruct.DicDataStruct, 0)
	errFindDicDataStructList := zorm.QueryStructList(session, finder, &dicDataStructList, page)

	//记录错误
	if errFindDicDataStructList != nil {
		errFindDicDataStructList := fmt.Errorf("permservice.FindDicDataStructList错误:%w", errFindDicDataStructList)
		logger.Error(errFindDicDataStructList)
		return nil, errFindDicDataStructList
	}

	return dicDataStructList, nil
}