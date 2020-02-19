package shop

import (
	"database/sql"
	"goshop/org/springrain/orm"
)

type User struct {
	orm.EntityStruct
	Id      string
	Account string
}

type User2 struct {
	orm.EntityStruct
	Id      string        `column:"id"`
	Account string        `column:"account"`
	Age     sql.NullInt32 `column:"age"`
	sex     string
}

func (user *User) GetTableName() string {
	return "t_user"
}
