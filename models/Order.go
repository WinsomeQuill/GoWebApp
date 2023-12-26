package models

import uuid2 "github.com/google/uuid"

type Order struct {
	UUID          uuid2.UUID `gorm:"not null"`
	UserID        uint       `gorm:"foreignkey:UserID; not null"`
	OrderStatusID uint       `gorm:"foreignkey:OrderStatusID; not null"`
	ItemID        uint       `gorm:"foreignkey:ItemID; not null"`
}
