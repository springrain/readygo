package permservice

import (
	"errors"
	"readygo/orm"
)

//保存用户角色中间表,session参数是为了保证在其他事务内,可以为nil
func SaveUserRoleStruct(session *orm.Session, userRoleStruct *permstruct.UserRoleStruct) error {

	if session != nil { //如果在其他的事务内
		err := orm.SaveStruct(session, userRoleStruct)
		if err != nil {
			return err
		}
		return nil
	}

	//不再其他事务内,新开事务
	_, err := orm.Transaction(func(session *orm.Session) (interface{}, error) {
		//事务下的业务代码

		err := orm.SaveStruct(session, userRoleStruct)
		if err != nil {
			return nil, err
		}
		return nil, nil

	})
	if err != nil {
		return err
	}
	return nil
}

//更新用户角色中间表,session参数是为了保证在其他事务内,可以为nil
func UpdateUserRoleStruct(session *orm.Session, userRoleStruct *permstruct.UserRoleStruct) error {

	if session != nil { //如果在其他的事务内
		err := orm.UpdateStruct(session, userRoleStruct)
		if err != nil {
			return err
		}
		return nil
	}

	//不再其他事务内,新开事务
	_, err := orm.Transaction(func(session *orm.Session) (interface{}, error) {
		//业务代码开始

		err := orm.UpdateStruct(session, userRoleStruct)
		if err != nil {
			return nil, err
		}

		return nil, nil

		//业务代码结束

	})
	if err != nil {
		return err
	}
	return nil
}

//删除用户角色中间表,session参数是为了保证在其他事务内,可以为nil
func DeleteUserRoleStruct(session *orm.Session, userRoleStruct *permstruct.UserRoleStruct) error {

	if session != nil { //如果在其他的事务内
		err := orm.UpdateStruct(session, userRoleStruct)
		if err != nil {
			return err
		}
		return nil
	}

	//不再其他事务内,新开事务
	_, err := orm.Transaction(func(session *orm.Session) (interface{}, error) {
		//业务代码开始

		err := orm.DeleteStruct(session, userRoleStruct)
		if err != nil {
			return nil, err
		}

		return nil, nil

		//业务代码结束

	})

	if err != nil {
		return err
	}
	return nil
}

//根据Id查询用户角色中间表信息,session参数是为了保证在其他事务内,可以为nil
func FindUserRoleStructById(session *orm.Session, id string) (*permstruct.UserRoleStruct, error) {
	//如果Id为空
	if len(id) < 1 {
		return nil, errors.New("id为空")
	}

	//根据Id查询
	finder := orm.NewSelectFinder(" WHERE id=?", id)
	userRoleStruct := permstruct.UserRoleStruct{}
	err := orm.QueryStruct(session, finder, &userRoleStruct)
	if err != nil {
		return nil, err
	}
	return &userRoleStruct, nil

}

//根据Finder查询用户角色中间表列表,session参数是为了保证在其他事务内,可以为nil
func FindUserRoleStructList(session *orm.Session, finder *orm.Finder, page *orm.Page) ([]permstruct.UserRoleStruct, error) {
	userRoleStructList := make([]permstruct.UserRoleStruct, 0)
	err := orm.QueryStructList(session, finder, &userRoleStructList, page)
	if err != nil {
		return nil, err
	}
	return userRoleStructList, nil
}
