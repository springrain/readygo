package permservice

import (
	"errors"
	"readygo/orm"
)

//保存角色菜单中间表,session参数是为了保证在其他事务内,可以为nil
func SaveRoleMenuStruct(session *orm.Session, roleMenuStruct *permstruct.RoleMenuStruct) error {

	if session != nil { //如果在其他的事务内
		err := orm.SaveStruct(session, roleMenuStruct)
		if err != nil {
			return err
		}
		return nil
	}

	//不再其他事务内,新开事务
	_, err := orm.Transaction(func(session *orm.Session) (interface{}, error) {
		//事务下的业务代码

		err := orm.SaveStruct(session, roleMenuStruct)
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

//更新角色菜单中间表,session参数是为了保证在其他事务内,可以为nil
func UpdateRoleMenuStruct(session *orm.Session, roleMenuStruct *permstruct.RoleMenuStruct) error {

	if session != nil { //如果在其他的事务内
		err := orm.UpdateStruct(session, roleMenuStruct)
		if err != nil {
			return err
		}
		return nil
	}

	//不再其他事务内,新开事务
	_, err := orm.Transaction(func(session *orm.Session) (interface{}, error) {
		//业务代码开始

		err := orm.UpdateStruct(session, roleMenuStruct)
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

//删除角色菜单中间表,session参数是为了保证在其他事务内,可以为nil
func DeleteRoleMenuStruct(session *orm.Session, roleMenuStruct *permstruct.RoleMenuStruct) error {

	if session != nil { //如果在其他的事务内
		err := orm.UpdateStruct(session, roleMenuStruct)
		if err != nil {
			return err
		}
		return nil
	}

	//不再其他事务内,新开事务
	_, err := orm.Transaction(func(session *orm.Session) (interface{}, error) {
		//业务代码开始

		err := orm.DeleteStruct(session, roleMenuStruct)
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

//根据Id查询角色菜单中间表信息,session参数是为了保证在其他事务内,可以为nil
func FindRoleMenuStructById(session *orm.Session, id string) (*permstruct.RoleMenuStruct, error) {
	//如果Id为空
	if len(id) < 1 {
		return nil, errors.New("id为空")
	}

	//根据Id查询
	finder := orm.NewSelectFinder(" WHERE id=?", id)
	roleMenuStruct := permstruct.RoleMenuStruct{}
	err := orm.QueryStruct(session, finder, &roleMenuStruct)
	if err != nil {
		return nil, err
	}
	return &roleMenuStruct, nil

}

//根据Finder查询角色菜单中间表列表,session参数是为了保证在其他事务内,可以为nil
func FindRoleMenuStructList(session *orm.Session, finder *orm.Finder, page *orm.Page) ([]permstruct.RoleMenuStruct, error) {
	roleMenuStructList := make([]permstruct.RoleMenuStruct, 0)
	err := orm.QueryStructList(session, finder, &roleMenuStructList, page)
	if err != nil {
		return nil, err
	}
	return roleMenuStructList, nil
}
