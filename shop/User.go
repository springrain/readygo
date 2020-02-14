package shop

import "goshop/org/goshop/frame/orm"

type User struct {
	orm.EntityStruct
	Id      string
	Account string
	age     int
}

type User2 struct {
	User
	Id      string
	Account string
	age     string
	sex     string
}

func (user *User) GetTableName() string {
	return "t_user"
}
