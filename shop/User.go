package shop

import "goshop/org/springrain/orm"

type User struct {
	orm.EntityStruct
	Id      string
	Account string
	age     int
}

type User2 struct {
	orm.EntityStruct
	ID      string `column:"id"`
	ACCOUNT string `column:"account"`
	Age     string
	sex     string
}

func (user *User) GetTableName() string {
	return "t_user"
}
