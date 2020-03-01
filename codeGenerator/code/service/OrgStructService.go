package permservice

import (
	"errors"
	"readygo/orm"
	"readygo/permission/permstruct"
)

//保存部门,session参数是为了保证在其他事务内,可以为nil
func SaveOrgStruct(session *orm.Session, orgStruct *permstruct.OrgStruct) error {

	if session != nil { //如果在其他的事务内
		err := orm.SaveStruct(session, orgStruct)
		if err != nil {
			return err
		}
		return nil
	}

	//不再其他事务内,新开事务
	_, err := orm.Transaction(func(session *orm.Session) (interface{}, error) {
		//事务下的业务代码

		err := orm.SaveStruct(session, orgStruct)
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

//更新部门,session参数是为了保证在其他事务内,可以为nil
func UpdateOrgStruct(session *orm.Session, orgStruct *permstruct.OrgStruct) error {

	if session != nil { //如果在其他的事务内
		err := orm.UpdateStruct(session, orgStruct)
		if err != nil {
			return err
		}
		return nil
	}

	//不再其他事务内,新开事务
	_, err := orm.Transaction(func(session *orm.Session) (interface{}, error) {
		//业务代码开始

		err := orm.UpdateStruct(session, orgStruct)
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

//删除部门,session参数是为了保证在其他事务内,可以为nil
func DeleteOrgStruct(session *orm.Session, orgStruct *permstruct.OrgStruct) error {

	if session != nil { //如果在其他的事务内
		err := orm.UpdateStruct(session, orgStruct)
		if err != nil {
			return err
		}
		return nil
	}

	//不再其他事务内,新开事务
	_, err := orm.Transaction(func(session *orm.Session) (interface{}, error) {
		//业务代码开始

		err := orm.DeleteStruct(session, orgStruct)
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

//根据Id查询部门信息,session参数是为了保证在其他事务内,可以为nil
func FindOrgStructById(session *orm.Session, id string) (*permstruct.OrgStruct, error) {
	//如果Id为空
	if len(id) < 1 {
		return nil, errors.New("id为空")
	}

	//根据Id查询
	finder := orm.NewSelectFinder(" WHERE id=?", id)
	orgStruct := permstruct.OrgStruct{}
	err := orm.QueryStruct(session, finder, &orgStruct)
	if err != nil {
		return nil, err
	}
	return &orgStruct, nil

}

//根据Finder查询部门列表,session参数是为了保证在其他事务内,可以为nil
func FindOrgStructList(session *orm.Session, finder *orm.Finder, page *orm.Page) ([]permstruct.OrgStruct, error) {
	orgStructList := make([]permstruct.OrgStruct, 0)
	err := orm.QueryStructList(session, finder, &orgStructList, page)
	if err != nil {
		return nil, err
	}
	return orgStructList, nil
}
