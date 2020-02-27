# readygo

#### 介绍
golang开发脚手架

#### 软件架构
基于gin和自研ORM


#### 例子

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
    ```
4.  查
    ```
	finder := orm.NewSelectFinder(permstruct.UserStructTableName)
	page := orm.NewPage()
	var users = make([]permstruct.UserStruct, 0)
	err := orm.QueryStructList(nil, finder, &users, &page)
    ```

