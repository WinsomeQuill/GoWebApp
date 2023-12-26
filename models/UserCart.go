package models

import (
	"GoWebApp/models/Dao"
	"GoWebApp/models/Dto"
)

type UserCart struct {
	User  Dto.UserDto
	Items []Dao.ItemDao
}
