package shop

import "goshop/org/goshop/frame/orm"

type User struct {
	*orm.EntityStruct
	Id      string
	Account string
	age     int
}

func (user *User) GetTableName() string {
	return "t_user"
}
