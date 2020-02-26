package permservice

import (
	"errors"
	"readygo/orm"
	"readygo/permission/permstruct"
)

//SaveUserStruct 保存用户
//如果入参session为nil或者没事务,则会使用本机的开启,并提交.如果session有事务,则只使用,不提交,有开启方提交事务.但是如果遇到错误或者异常,虽然不是事务的开启方,也会回滚事务,让事务尽早回滚
func SaveUserStruct(session *orm.Session, userStruct *permstruct.UserStruct) error {

	//匿名函数return的error如果不为nil,事务就会回滚
	_, err := orm.Transaction(session, func(session *orm.Session) (interface{}, error) {

		//事务下的业务代码开始

		err := orm.SaveStruct(session, userStruct)
		if err != nil {
			return nil, err
		}
		return nil, nil

		//事务下的业务代码结束

	})
	if err != nil {
		return err
	}
	return nil
}

//UpdateUserStruct 更新用户
//如果入参session为nil或者没事务,则会使用本机的开启,并提交.如果session有事务,则只使用,不提交,有开启方提交事务.但是如果遇到错误或者异常,虽然不是事务的开启方,也会回滚事务,让事务尽早回滚
func UpdateUserStruct(session *orm.Session, userStruct *permstruct.UserStruct) error {

	//匿名函数return的error如果不为nil,事务就会回滚
	_, err := orm.Transaction(session, func(session *orm.Session) (interface{}, error) {

		//事务下的业务代码开始

		err := orm.UpdateStruct(session, userStruct)
		if err != nil {
			return nil, err
		}

		return nil, nil

		//事务下的业务代码结束

	})
	if err != nil {
		return err
	}
	return nil
}

//DeleteUserStruct 删除用户
//如果入参session为nil或者没事务,则会使用本机的开启,并提交.如果session有事务,则只使用,不提交,有开启方提交事务.但是如果遇到错误或者异常,虽然不是事务的开启方,也会回滚事务,让事务尽早回滚
func DeleteUserStruct(session *orm.Session, userStruct *permstruct.UserStruct) error {

	//匿名函数return的error如果不为nil,事务就会回滚
	_, err := orm.Transaction(session, func(session *orm.Session) (interface{}, error) {

		//事务下的业务代码开始

		err := orm.DeleteStruct(session, userStruct)
		if err != nil {
			return nil, err
		}

		return nil, nil

		//事务下的业务代码结束

	})

	if err != nil {
		return err
	}
	return nil
}

//FindUserStructById 根据Id查询用户信息
//session如果为nil,则会使用默认的datasource进行无事务查询
func FindUserStructById(session *orm.Session, id string) (*permstruct.UserStruct, error) {
	//如果Id为空
	if len(id) < 1 {
		return nil, errors.New("id为空")
	}

	//根据Id查询
	finder := orm.NewSelectFinder(" WHERE id=?", id)
	userStruct := permstruct.UserStruct{}
	err := orm.QueryStruct(session, finder, &userStruct)
	if err != nil {
		return nil, err
	}
	return &userStruct, nil

}

//FindUserStructList 根据Finder查询用户列表
//session如果为nil,则会使用默认的datasource进行无事务查询
func FindUserStructList(session *orm.Session, finder *orm.Finder, page *orm.Page) ([]permstruct.UserStruct, error) {
	userStructList := make([]permstruct.UserStruct, 0)
	err := orm.QueryStructList(session, finder, &userStructList, page)
	if err != nil {
		return nil, err
	}
	return userStructList, nil
}
