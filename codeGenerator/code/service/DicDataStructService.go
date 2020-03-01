package permservice

import (
	"errors"
	"readygo/orm"
)

//保存公共字典,session参数是为了保证在其他事务内,可以为nil
func SaveDicDataStruct(session *orm.Session, dicDataStruct *permstruct.DicDataStruct) error {

	if session != nil { //如果在其他的事务内
		err := orm.SaveStruct(session, dicDataStruct)
		if err != nil {
			return err
		}
		return nil
	}

	//不再其他事务内,新开事务
	_, err := orm.Transaction(func(session *orm.Session) (interface{}, error) {
		//事务下的业务代码

		err := orm.SaveStruct(session, dicDataStruct)
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

//更新公共字典,session参数是为了保证在其他事务内,可以为nil
func UpdateDicDataStruct(session *orm.Session, dicDataStruct *permstruct.DicDataStruct) error {

	if session != nil { //如果在其他的事务内
		err := orm.UpdateStruct(session, dicDataStruct)
		if err != nil {
			return err
		}
		return nil
	}

	//不再其他事务内,新开事务
	_, err := orm.Transaction(func(session *orm.Session) (interface{}, error) {
		//业务代码开始

		err := orm.UpdateStruct(session, dicDataStruct)
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

//删除公共字典,session参数是为了保证在其他事务内,可以为nil
func DeleteDicDataStruct(session *orm.Session, dicDataStruct *permstruct.DicDataStruct) error {

	if session != nil { //如果在其他的事务内
		err := orm.UpdateStruct(session, dicDataStruct)
		if err != nil {
			return err
		}
		return nil
	}

	//不再其他事务内,新开事务
	_, err := orm.Transaction(func(session *orm.Session) (interface{}, error) {
		//业务代码开始

		err := orm.DeleteStruct(session, dicDataStruct)
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

//根据Id查询公共字典信息,session参数是为了保证在其他事务内,可以为nil
func FindDicDataStructById(session *orm.Session, id string) (*permstruct.DicDataStruct, error) {
	//如果Id为空
	if len(id) < 1 {
		return nil, errors.New("id为空")
	}

	//根据Id查询
	finder := orm.NewSelectFinder(" WHERE id=?", id)
	dicDataStruct := permstruct.DicDataStruct{}
	err := orm.QueryStruct(session, finder, &dicDataStruct)
	if err != nil {
		return nil, err
	}
	return &dicDataStruct, nil

}

//根据Finder查询公共字典列表,session参数是为了保证在其他事务内,可以为nil
func FindDicDataStructList(session *orm.Session, finder *orm.Finder, page *orm.Page) ([]permstruct.DicDataStruct, error) {
	dicDataStructList := make([]permstruct.DicDataStruct, 0)
	err := orm.QueryStructList(session, finder, &dicDataStructList, page)
	if err != nil {
		return nil, err
	}
	return dicDataStructList, nil
}
