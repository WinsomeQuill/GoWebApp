package models

type OrderStatus struct {
	ID   int64  `gorm:"primaryKey"`
	Name string `gorm:"not null"`
}
