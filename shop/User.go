package shop

import "goshop/org/springrain/orm"

type User struct {
	orm.EntityStruct
	Id      string
	Account string
}

type User2 struct {
	orm.EntityStruct
	Id      string `column:"id"`
	Account string `column:"account"`
	Age     int    `column:"age"`
	sex     string
}

func (user *User) GetTableName() string {
	return "t_user"
}
