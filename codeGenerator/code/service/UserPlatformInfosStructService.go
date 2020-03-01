package permservice

import (
	"errors"
	"readygo/orm"
)

//保存用户平台信息表,session参数是为了保证在其他事务内,可以为nil
func SaveUserPlatformInfosStruct(session *orm.Session, userPlatformInfosStruct *permstruct.UserPlatformInfosStruct) error {

	if session != nil { //如果在其他的事务内
		err := orm.SaveStruct(session, userPlatformInfosStruct)
		if err != nil {
			return err
		}
		return nil
	}

	//不再其他事务内,新开事务
	_, err := orm.Transaction(func(session *orm.Session) (interface{}, error) {
		//事务下的业务代码

		err := orm.SaveStruct(session, userPlatformInfosStruct)
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

//更新用户平台信息表,session参数是为了保证在其他事务内,可以为nil
func UpdateUserPlatformInfosStruct(session *orm.Session, userPlatformInfosStruct *permstruct.UserPlatformInfosStruct) error {

	if session != nil { //如果在其他的事务内
		err := orm.UpdateStruct(session, userPlatformInfosStruct)
		if err != nil {
			return err
		}
		return nil
	}

	//不再其他事务内,新开事务
	_, err := orm.Transaction(func(session *orm.Session) (interface{}, error) {
		//业务代码开始

		err := orm.UpdateStruct(session, userPlatformInfosStruct)
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

//删除用户平台信息表,session参数是为了保证在其他事务内,可以为nil
func DeleteUserPlatformInfosStruct(session *orm.Session, userPlatformInfosStruct *permstruct.UserPlatformInfosStruct) error {

	if session != nil { //如果在其他的事务内
		err := orm.UpdateStruct(session, userPlatformInfosStruct)
		if err != nil {
			return err
		}
		return nil
	}

	//不再其他事务内,新开事务
	_, err := orm.Transaction(func(session *orm.Session) (interface{}, error) {
		//业务代码开始

		err := orm.DeleteStruct(session, userPlatformInfosStruct)
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

//根据Id查询用户平台信息表信息,session参数是为了保证在其他事务内,可以为nil
func FindUserPlatformInfosStructById(session *orm.Session, id string) (*permstruct.UserPlatformInfosStruct, error) {
	//如果Id为空
	if len(id) < 1 {
		return nil, errors.New("id为空")
	}

	//根据Id查询
	finder := orm.NewSelectFinder(" WHERE id=?", id)
	userPlatformInfosStruct := permstruct.UserPlatformInfosStruct{}
	err := orm.QueryStruct(session, finder, &userPlatformInfosStruct)
	if err != nil {
		return nil, err
	}
	return &userPlatformInfosStruct, nil

}

//根据Finder查询用户平台信息表列表,session参数是为了保证在其他事务内,可以为nil
func FindUserPlatformInfosStructList(session *orm.Session, finder *orm.Finder, page *orm.Page) ([]permstruct.UserPlatformInfosStruct, error) {
	userPlatformInfosStructList := make([]permstruct.UserPlatformInfosStruct, 0)
	err := orm.QueryStructList(session, finder, &userPlatformInfosStructList, page)
	if err != nil {
		return nil, err
	}
	return userPlatformInfosStructList, nil
}
