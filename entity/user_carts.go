package entity

import "gorm.io/gorm"

type UserCartItem struct {
	gorm.Model
	MenuID      uint   `json:"menu_id"`
	UserID      uint   `json:"user_id"`
	Quantity    int    `json:"quantity"`
	MenuOptions string `json:"menu_options"`
	Menu        Menu   `json:"menu"`
}
