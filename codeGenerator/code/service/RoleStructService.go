package permservice

import (
	"errors"
	"readygo/orm"
)

//保存角色,session参数是为了保证在其他事务内,可以为nil
func SaveRoleStruct(session *orm.Session, roleStruct *permstruct.RoleStruct) error {

	if session != nil { //如果在其他的事务内
		err := orm.SaveStruct(session, roleStruct)
		if err != nil {
			return err
		}
		return nil
	}

	//不再其他事务内,新开事务
	_, err := orm.Transaction(func(session *orm.Session) (interface{}, error) {
		//事务下的业务代码

		err := orm.SaveStruct(session, roleStruct)
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

//更新角色,session参数是为了保证在其他事务内,可以为nil
func UpdateRoleStruct(session *orm.Session, roleStruct *permstruct.RoleStruct) error {

	if session != nil { //如果在其他的事务内
		err := orm.UpdateStruct(session, roleStruct)
		if err != nil {
			return err
		}
		return nil
	}

	//不再其他事务内,新开事务
	_, err := orm.Transaction(func(session *orm.Session) (interface{}, error) {
		//业务代码开始

		err := orm.UpdateStruct(session, roleStruct)
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

//删除角色,session参数是为了保证在其他事务内,可以为nil
func DeleteRoleStruct(session *orm.Session, roleStruct *permstruct.RoleStruct) error {

	if session != nil { //如果在其他的事务内
		err := orm.UpdateStruct(session, roleStruct)
		if err != nil {
			return err
		}
		return nil
	}

	//不再其他事务内,新开事务
	_, err := orm.Transaction(func(session *orm.Session) (interface{}, error) {
		//业务代码开始

		err := orm.DeleteStruct(session, roleStruct)
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

//根据Id查询角色信息,session参数是为了保证在其他事务内,可以为nil
func FindRoleStructById(session *orm.Session, id string) (*permstruct.RoleStruct, error) {
	//如果Id为空
	if len(id) < 1 {
		return nil, errors.New("id为空")
	}

	//根据Id查询
	finder := orm.NewSelectFinder(" WHERE id=?", id)
	roleStruct := permstruct.RoleStruct{}
	err := orm.QueryStruct(session, finder, &roleStruct)
	if err != nil {
		return nil, err
	}
	return &roleStruct, nil

}

//根据Finder查询角色列表,session参数是为了保证在其他事务内,可以为nil
func FindRoleStructList(session *orm.Session, finder *orm.Finder, page *orm.Page) ([]permstruct.RoleStruct, error) {
	roleStructList := make([]permstruct.RoleStruct, 0)
	err := orm.QueryStructList(session, finder, &roleStructList, page)
	if err != nil {
		return nil, err
	}
	return roleStructList, nil
}
