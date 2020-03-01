package permservice

import (
	"errors"
	"readygo/orm"
)

//保存微信号需要的配置信息,session参数是为了保证在其他事务内,可以为nil
func SaveWxCpconfigStruct(session *orm.Session, wxCpconfigStruct *permstruct.WxCpconfigStruct) error {

	if session != nil { //如果在其他的事务内
		err := orm.SaveStruct(session, wxCpconfigStruct)
		if err != nil {
			return err
		}
		return nil
	}

	//不再其他事务内,新开事务
	_, err := orm.Transaction(func(session *orm.Session) (interface{}, error) {
		//事务下的业务代码

		err := orm.SaveStruct(session, wxCpconfigStruct)
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

//更新微信号需要的配置信息,session参数是为了保证在其他事务内,可以为nil
func UpdateWxCpconfigStruct(session *orm.Session, wxCpconfigStruct *permstruct.WxCpconfigStruct) error {

	if session != nil { //如果在其他的事务内
		err := orm.UpdateStruct(session, wxCpconfigStruct)
		if err != nil {
			return err
		}
		return nil
	}

	//不再其他事务内,新开事务
	_, err := orm.Transaction(func(session *orm.Session) (interface{}, error) {
		//业务代码开始

		err := orm.UpdateStruct(session, wxCpconfigStruct)
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

//删除微信号需要的配置信息,session参数是为了保证在其他事务内,可以为nil
func DeleteWxCpconfigStruct(session *orm.Session, wxCpconfigStruct *permstruct.WxCpconfigStruct) error {

	if session != nil { //如果在其他的事务内
		err := orm.UpdateStruct(session, wxCpconfigStruct)
		if err != nil {
			return err
		}
		return nil
	}

	//不再其他事务内,新开事务
	_, err := orm.Transaction(func(session *orm.Session) (interface{}, error) {
		//业务代码开始

		err := orm.DeleteStruct(session, wxCpconfigStruct)
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

//根据Id查询微信号需要的配置信息信息,session参数是为了保证在其他事务内,可以为nil
func FindWxCpconfigStructById(session *orm.Session, id string) (*permstruct.WxCpconfigStruct, error) {
	//如果Id为空
	if len(id) < 1 {
		return nil, errors.New("id为空")
	}

	//根据Id查询
	finder := orm.NewSelectFinder(" WHERE id=?", id)
	wxCpconfigStruct := permstruct.WxCpconfigStruct{}
	err := orm.QueryStruct(session, finder, &wxCpconfigStruct)
	if err != nil {
		return nil, err
	}
	return &wxCpconfigStruct, nil

}

//根据Finder查询微信号需要的配置信息列表,session参数是为了保证在其他事务内,可以为nil
func FindWxCpconfigStructList(session *orm.Session, finder *orm.Finder, page *orm.Page) ([]permstruct.WxCpconfigStruct, error) {
	wxCpconfigStructList := make([]permstruct.WxCpconfigStruct, 0)
	err := orm.QueryStructList(session, finder, &wxCpconfigStructList, page)
	if err != nil {
		return nil, err
	}
	return wxCpconfigStructList, nil
}
