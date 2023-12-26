package models

import (
	"GoWebApp/models/Dao"
	"GoWebApp/models/Dto"
)

type UserOrder struct {
	User   Dto.UserDto
	Items  []Dao.ItemDao
	Status OrderStatus
}
