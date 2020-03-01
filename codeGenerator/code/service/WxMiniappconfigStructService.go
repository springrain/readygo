package permservice

import (
	"errors"
	"readygo/orm"
)

//保存小程序配置表,session参数是为了保证在其他事务内,可以为nil
func SaveWxMiniappconfigStruct(session *orm.Session, wxMiniappconfigStruct *permstruct.WxMiniappconfigStruct) error {

	if session != nil { //如果在其他的事务内
		err := orm.SaveStruct(session, wxMiniappconfigStruct)
		if err != nil {
			return err
		}
		return nil
	}

	//不再其他事务内,新开事务
	_, err := orm.Transaction(func(session *orm.Session) (interface{}, error) {
		//事务下的业务代码

		err := orm.SaveStruct(session, wxMiniappconfigStruct)
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

//更新小程序配置表,session参数是为了保证在其他事务内,可以为nil
func UpdateWxMiniappconfigStruct(session *orm.Session, wxMiniappconfigStruct *permstruct.WxMiniappconfigStruct) error {

	if session != nil { //如果在其他的事务内
		err := orm.UpdateStruct(session, wxMiniappconfigStruct)
		if err != nil {
			return err
		}
		return nil
	}

	//不再其他事务内,新开事务
	_, err := orm.Transaction(func(session *orm.Session) (interface{}, error) {
		//业务代码开始

		err := orm.UpdateStruct(session, wxMiniappconfigStruct)
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

//删除小程序配置表,session参数是为了保证在其他事务内,可以为nil
func DeleteWxMiniappconfigStruct(session *orm.Session, wxMiniappconfigStruct *permstruct.WxMiniappconfigStruct) error {

	if session != nil { //如果在其他的事务内
		err := orm.UpdateStruct(session, wxMiniappconfigStruct)
		if err != nil {
			return err
		}
		return nil
	}

	//不再其他事务内,新开事务
	_, err := orm.Transaction(func(session *orm.Session) (interface{}, error) {
		//业务代码开始

		err := orm.DeleteStruct(session, wxMiniappconfigStruct)
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

//根据Id查询小程序配置表信息,session参数是为了保证在其他事务内,可以为nil
func FindWxMiniappconfigStructById(session *orm.Session, id string) (*permstruct.WxMiniappconfigStruct, error) {
	//如果Id为空
	if len(id) < 1 {
		return nil, errors.New("id为空")
	}

	//根据Id查询
	finder := orm.NewSelectFinder(" WHERE id=?", id)
	wxMiniappconfigStruct := permstruct.WxMiniappconfigStruct{}
	err := orm.QueryStruct(session, finder, &wxMiniappconfigStruct)
	if err != nil {
		return nil, err
	}
	return &wxMiniappconfigStruct, nil

}

//根据Finder查询小程序配置表列表,session参数是为了保证在其他事务内,可以为nil
func FindWxMiniappconfigStructList(session *orm.Session, finder *orm.Finder, page *orm.Page) ([]permstruct.WxMiniappconfigStruct, error) {
	wxMiniappconfigStructList := make([]permstruct.WxMiniappconfigStruct, 0)
	err := orm.QueryStructList(session, finder, &wxMiniappconfigStructList, page)
	if err != nil {
		return nil, err
	}
	return wxMiniappconfigStructList, nil
}
