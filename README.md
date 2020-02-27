# readygo

#### 介绍
golang开发脚手架

#### 软件架构
基于gin和自研ORM  
[自带代码生成器](https://gitee.com/chunanyong/readygo/tree/master/codeGenerator)  
使用orm.Finder作为sql载体,所有的sql语句最终都是通过finder执行.  
支持事务传播  


#### 例子
具体可以参照 [UserStructService.go](https://gitee.com/chunanyong/readygo/tree/master/permission/permservice)

1.  增
    ```
    var user permstruct.UserStruct
    err := orm.SaveStruct(nil, &user)
    ```
2.  删
    ```
    err := orm.DeleteStruct(nil,&user)
    ```
  
3.  改
    ```
    err := orm.UpdateStruct(nil,&user)
    //finder更新
    err := orm.UpdateFinder(nil,finder)
    ```
4.  查
    ```
	finder := orm.NewSelectFinder(permstruct.UserStructTableName)
	page := orm.NewPage()
	var users = make([]permstruct.UserStruct, 0)
	err := orm.QueryStructList(nil, finder, &users, &page)
    ```
5.  事务传播
    ```
    //匿名函数return的error如果不为nil,事务就会回滚
	_, errSaveUserStruct := orm.Transaction(session, func(session *orm.Session) (interface{}, error) {

		//事务下的业务代码开始
		errSaveUserStruct := orm.SaveStruct(session, userStruct)

		if errSaveUserStruct != nil {
			return nil, errSaveUserStruct
		}

		return nil, nil
		//事务下的业务代码结束

	})
    ```

