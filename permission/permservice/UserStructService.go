package permservice

import (
	"errors"
	"readygo/orm"
)

//保存用户,session参数是为了保证在其他事务内,可以为nil
func SaveUserStruct(session *orm.Session, userStruct *permstruct.UserStruct) error {

	if session != nil { //如果在其他的事务内
		err := orm.GetDefaultDao().SaveStruct(session, userStruct)
		if err != nil {
			return err
		}
		return nil
	}

	//不再其他事务内,新开事务
	_, err := orm.GetDefaultDao().Transaction(func(session *orm.Session) (interface{}, error) {
		//事务下的业务代码

		err := orm.GetDefaultDao().SaveStruct(session, userStruct)
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

//更新用户,session参数是为了保证在其他事务内,可以为nil
func UpdateUserStruct(session *orm.Session, userStruct *permstruct.UserStruct) error {

	if session != nil { //如果在其他的事务内
		err := orm.GetDefaultDao().UpdateStruct(session, userStruct)
		if err != nil {
			return err
		}
		return nil
	}

	//不再其他事务内,新开事务
	_, err := orm.GetDefaultDao().Transaction(func(session *orm.Session) (interface{}, error) {
		//业务代码开始

		err := orm.GetDefaultDao().UpdateStruct(session, userStruct)
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

//删除用户,session参数是为了保证在其他事务内,可以为nil
func DeleteUserStruct(session *orm.Session, userStruct *permstruct.UserStruct) error {

	if session != nil { //如果在其他的事务内
		err := orm.GetDefaultDao().UpdateStruct(session, userStruct)
		if err != nil {
			return err
		}
		return nil
	}

	//不再其他事务内,新开事务
	_, err := orm.GetDefaultDao().Transaction(func(session *orm.Session) (interface{}, error) {
		//业务代码开始

		err := orm.GetDefaultDao().DeleteStruct(session, userStruct)
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

//根据Id查询用户信息,session参数是为了保证在其他事务内,可以为nil
func FindUserStructById(session *orm.Session, id string) (*permstruct.UserStruct, error) {
	//如果Id为空
	if len(id) < 1 {
		return nil, errors.New("id为空")
	}

	//根据Id查询
	finder := orm.NewSelectFinder(" WHERE id=?", id)
	userStruct := permstruct.UserStruct{}
	err := orm.GetDefaultDao().QueryStruct(session, finder, &userStruct)
	if err != nil {
		return nil, err
	}
	return &userStruct, nil

}

//根据Finder查询用户列表,session参数是为了保证在其他事务内,可以为nil
func FindUserStructList(session *orm.Session, finder *orm.Finder, page *orm.Page) ([]permstruct.UserStruct, error) {
	userStructList := make([]permstruct.UserStruct, 0)
	err := orm.GetDefaultDao().QueryStructList(session, finder, &userStructList, page)
	if err != nil {
		return nil, err
	}
	return userStructList, nil
}
