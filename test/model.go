package main

import (
	"readygo/orm"
	"time"
)

type Model struct {
	Id        int        `column:"id"`
	CreatedAt time.Time  `column:"created_at"`
	UpdatedAt time.Time  `column:"updated_at"`
	DeletedAt *time.Time `column:"deleted_at"`
}

type Toy struct {
	orm.EntityStruct
	Id        int    `column:"id"`
	Name      string `column:"name"`
	OwnerID   int    `column:"owner_id"`
	OwnerType string `column:"owner_type"`
}

func (t Toy) GetTableName() string {
	return "toys"
}

func (t Toy) GetPKColumnName() string {
	return "id"
}

type User struct {
	orm.EntityStruct
	Model
	UserName       string
	PasswordDigest string
	Nickname       string
	Status         string
	Avatar         string
	Languages      []*Language
}

func (u User) GetTableName() string {
	return "user"
}

type Language struct {
	orm.EntityStruct
	Model
	Name  string `column:"name"`
	Users []*User
}

func (l Language) GetTableName() string {
	return "languages"
}

type UserLanguages struct {
	orm.EntityStruct
	Model
	UserId     uint
	LanguageId uint
}

type Product struct {
	orm.EntityStruct
	Model
	Name   string
	Amount float32
}

type OrderDetail struct {
	orm.EntityStruct
	Model
	Product    Product
	ProductId  uint
	Discount   float32
	LastAmount float32
	OrderId    uint
}

type Order struct {
	orm.EntityStruct
	Model
	Total        float32
	Coupon       float32
	OrderDetails []OrderDetail
}
