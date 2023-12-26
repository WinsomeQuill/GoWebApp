package models

type Cart struct {
	UserID uint `gorm:"foreignkey:UserID; not null"`
	ItemID uint `gorm:"foreignkey:ItemID; not null"`
	Count  uint `gorm:"not null"`
}
