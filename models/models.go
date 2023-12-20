package models

import (
	uuid2 "github.com/google/uuid"
	"gorm.io/gorm"
)

type Item struct {
	gorm.Model
	Name  string `gorm:"not null"`
	Price uint   `gorm:"not null"`
	Count uint   `gorm:"not null"`
}

type User struct {
	gorm.Model
	Name  string `gorm:"not null"`
	Email string `gorm:"not null"`
}

type Cart struct {
	UserID uint `gorm:"foreignkey:UserID; not null"`
	ItemID uint `gorm:"foreignkey:ItemID; not null"`
	Count  uint `gorm:"not null"`
}

type Order struct {
	UUID          uuid2.UUID `gorm:"not null"`
	UserID        uint       `gorm:"foreignkey:UserID; not null"`
	OrderStatusID uint       `gorm:"foreignkey:OrderStatusID; not null"`
	ItemID        uint       `gorm:"foreignkey:ItemID; not null"`
}

type OrderStatus struct {
	Name string `gorm:"not null"`
}

type ItemDao struct {
	Name         string
	Count        uint
	PricePerUnit uint
	PriceTotal   uint
}

type UserCart struct {
	User  UserDto
	Items []ItemDao
}

type UserOrder struct {
	User   UserDto
	Items  []ItemDao
	Status OrderStatus
}
