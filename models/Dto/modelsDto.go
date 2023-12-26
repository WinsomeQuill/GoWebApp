package Dto

import uuid2 "github.com/google/uuid"

type ItemToCartUserDto struct {
	User UserDto `json:"user"`
	Item ItemDto `json:"item"`
}

type UserDto struct {
	Name  string `json:"userName"`
	Email string `json:"email"`
}

type ItemDto struct {
	Name  string `json:"itemName"`
	Count int64  `json:"count"`
}

type UpdateOrderStatusDto struct {
	Uuid       uuid2.UUID `json:"uuid"`
	StatusName string     `json:"statusName"`
}
