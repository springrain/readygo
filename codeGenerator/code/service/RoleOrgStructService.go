package permservice

import (
	"errors"
	"readygo/orm"
)

//保存角色部门中间表,session参数是为了保证在其他事务内,可以为nil
func SaveRoleOrgStruct(session *orm.Session, roleOrgStruct *permstruct.RoleOrgStruct) error {

	if session != nil { //如果在其他的事务内
		err := orm.SaveStruct(session, roleOrgStruct)
		if err != nil {
			return err
		}
		return nil
	}

	//不再其他事务内,新开事务
	_, err := orm.Transaction(func(session *orm.Session) (interface{}, error) {
		//事务下的业务代码

		err := orm.SaveStruct(session, roleOrgStruct)
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

//更新角色部门中间表,session参数是为了保证在其他事务内,可以为nil
func UpdateRoleOrgStruct(session *orm.Session, roleOrgStruct *permstruct.RoleOrgStruct) error {

	if session != nil { //如果在其他的事务内
		err := orm.UpdateStruct(session, roleOrgStruct)
		if err != nil {
			return err
		}
		return nil
	}

	//不再其他事务内,新开事务
	_, err := orm.Transaction(func(session *orm.Session) (interface{}, error) {
		//业务代码开始

		err := orm.UpdateStruct(session, roleOrgStruct)
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

//删除角色部门中间表,session参数是为了保证在其他事务内,可以为nil
func DeleteRoleOrgStruct(session *orm.Session, roleOrgStruct *permstruct.RoleOrgStruct) error {

	if session != nil { //如果在其他的事务内
		err := orm.UpdateStruct(session, roleOrgStruct)
		if err != nil {
			return err
		}
		return nil
	}

	//不再其他事务内,新开事务
	_, err := orm.Transaction(func(session *orm.Session) (interface{}, error) {
		//业务代码开始

		err := orm.DeleteStruct(session, roleOrgStruct)
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

//根据Id查询角色部门中间表信息,session参数是为了保证在其他事务内,可以为nil
func FindRoleOrgStructById(session *orm.Session, id string) (*permstruct.RoleOrgStruct, error) {
	//如果Id为空
	if len(id) < 1 {
		return nil, errors.New("id为空")
	}

	//根据Id查询
	finder := orm.NewSelectFinder(" WHERE id=?", id)
	roleOrgStruct := permstruct.RoleOrgStruct{}
	err := orm.QueryStruct(session, finder, &roleOrgStruct)
	if err != nil {
		return nil, err
	}
	return &roleOrgStruct, nil

}

//根据Finder查询角色部门中间表列表,session参数是为了保证在其他事务内,可以为nil
func FindRoleOrgStructList(session *orm.Session, finder *orm.Finder, page *orm.Page) ([]permstruct.RoleOrgStruct, error) {
	roleOrgStructList := make([]permstruct.RoleOrgStruct, 0)
	err := orm.QueryStructList(session, finder, &roleOrgStructList, page)
	if err != nil {
		return nil, err
	}
	return roleOrgStructList, nil
}
