package permservice

import (
	"errors"
	"readygo/orm"
)

//保存菜单,session参数是为了保证在其他事务内,可以为nil
func SaveMenuStruct(session *orm.Session, menuStruct *permstruct.MenuStruct) error {

	if session != nil { //如果在其他的事务内
		err := orm.SaveStruct(session, menuStruct)
		if err != nil {
			return err
		}
		return nil
	}

	//不再其他事务内,新开事务
	_, err := orm.Transaction(func(session *orm.Session) (interface{}, error) {
		//事务下的业务代码

		err := orm.SaveStruct(session, menuStruct)
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

//更新菜单,session参数是为了保证在其他事务内,可以为nil
func UpdateMenuStruct(session *orm.Session, menuStruct *permstruct.MenuStruct) error {

	if session != nil { //如果在其他的事务内
		err := orm.UpdateStruct(session, menuStruct)
		if err != nil {
			return err
		}
		return nil
	}

	//不再其他事务内,新开事务
	_, err := orm.Transaction(func(session *orm.Session) (interface{}, error) {
		//业务代码开始

		err := orm.UpdateStruct(session, menuStruct)
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

//删除菜单,session参数是为了保证在其他事务内,可以为nil
func DeleteMenuStruct(session *orm.Session, menuStruct *permstruct.MenuStruct) error {

	if session != nil { //如果在其他的事务内
		err := orm.UpdateStruct(session, menuStruct)
		if err != nil {
			return err
		}
		return nil
	}

	//不再其他事务内,新开事务
	_, err := orm.Transaction(func(session *orm.Session) (interface{}, error) {
		//业务代码开始

		err := orm.DeleteStruct(session, menuStruct)
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

//根据Id查询菜单信息,session参数是为了保证在其他事务内,可以为nil
func FindMenuStructById(session *orm.Session, id string) (*permstruct.MenuStruct, error) {
	//如果Id为空
	if len(id) < 1 {
		return nil, errors.New("id为空")
	}

	//根据Id查询
	finder := orm.NewSelectFinder(" WHERE id=?", id)
	menuStruct := permstruct.MenuStruct{}
	err := orm.QueryStruct(session, finder, &menuStruct)
	if err != nil {
		return nil, err
	}
	return &menuStruct, nil

}

//根据Finder查询菜单列表,session参数是为了保证在其他事务内,可以为nil
func FindMenuStructList(session *orm.Session, finder *orm.Finder, page *orm.Page) ([]permstruct.MenuStruct, error) {
	menuStructList := make([]permstruct.MenuStruct, 0)
	err := orm.QueryStructList(session, finder, &menuStructList, page)
	if err != nil {
		return nil, err
	}
	return menuStructList, nil
}
