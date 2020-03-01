package permservice

import (
	"errors"
	"readygo/orm"
)

//保存支付宝的配置信息,session参数是为了保证在其他事务内,可以为nil
func SaveAliPayconfigStruct(session *orm.Session, aliPayconfigStruct *permstruct.AliPayconfigStruct) error {

	if session != nil { //如果在其他的事务内
		err := orm.SaveStruct(session, aliPayconfigStruct)
		if err != nil {
			return err
		}
		return nil
	}

	//不再其他事务内,新开事务
	_, err := orm.Transaction(func(session *orm.Session) (interface{}, error) {
		//事务下的业务代码

		err := orm.SaveStruct(session, aliPayconfigStruct)
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

//更新支付宝的配置信息,session参数是为了保证在其他事务内,可以为nil
func UpdateAliPayconfigStruct(session *orm.Session, aliPayconfigStruct *permstruct.AliPayconfigStruct) error {

	if session != nil { //如果在其他的事务内
		err := orm.UpdateStruct(session, aliPayconfigStruct)
		if err != nil {
			return err
		}
		return nil
	}

	//不再其他事务内,新开事务
	_, err := orm.Transaction(func(session *orm.Session) (interface{}, error) {
		//业务代码开始

		err := orm.UpdateStruct(session, aliPayconfigStruct)
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

//删除支付宝的配置信息,session参数是为了保证在其他事务内,可以为nil
func DeleteAliPayconfigStruct(session *orm.Session, aliPayconfigStruct *permstruct.AliPayconfigStruct) error {

	if session != nil { //如果在其他的事务内
		err := orm.UpdateStruct(session, aliPayconfigStruct)
		if err != nil {
			return err
		}
		return nil
	}

	//不再其他事务内,新开事务
	_, err := orm.Transaction(func(session *orm.Session) (interface{}, error) {
		//业务代码开始

		err := orm.DeleteStruct(session, aliPayconfigStruct)
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

//根据Id查询支付宝的配置信息信息,session参数是为了保证在其他事务内,可以为nil
func FindAliPayconfigStructById(session *orm.Session, id string) (*permstruct.AliPayconfigStruct, error) {
	//如果Id为空
	if len(id) < 1 {
		return nil, errors.New("id为空")
	}

	//根据Id查询
	finder := orm.NewSelectFinder(" WHERE id=?", id)
	aliPayconfigStruct := permstruct.AliPayconfigStruct{}
	err := orm.QueryStruct(session, finder, &aliPayconfigStruct)
	if err != nil {
		return nil, err
	}
	return &aliPayconfigStruct, nil

}

//根据Finder查询支付宝的配置信息列表,session参数是为了保证在其他事务内,可以为nil
func FindAliPayconfigStructList(session *orm.Session, finder *orm.Finder, page *orm.Page) ([]permstruct.AliPayconfigStruct, error) {
	aliPayconfigStructList := make([]permstruct.AliPayconfigStruct, 0)
	err := orm.QueryStructList(session, finder, &aliPayconfigStructList, page)
	if err != nil {
		return nil, err
	}
	return aliPayconfigStructList, nil
}
