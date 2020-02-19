package main

import (
	"goshop/org/springrain/orm"
	"time"
)


type Model struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}




type User struct {
	orm.EntityStruct
	Model
	UserName       string
	PasswordDigest string
	Nickname       string
	Status         string
	Avatar         string
 	Languages []*Language
}

type Language struct {
	orm.EntityStruct
	Model
	Name string
	Users []*User
}

type UserLanguages struct {
	orm.EntityStruct
	Model
	UserId    uint
	LanguageId uint
}




type Product struct {
	orm.EntityStruct
	Model
	Name string
	Amount float32

}

type OrderDetail struct {

	orm.EntityStruct
	Model
	Product Product
	ProductId uint
	Discount float32
	LastAmount float32
	OrderId uint
}

type Order struct {
	orm.EntityStruct
	Model
	Total float32
	Coupon float32
	OrderDetails []OrderDetail
}
