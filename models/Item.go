package models

import "gorm.io/gorm"

type Item struct {
	gorm.Model
	Name  string  `gorm:"not null"`
	Price float32 `gorm:"not null"`
	Count int64   `gorm:"not null"`
}
