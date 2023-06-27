package entity

import (
	"gorm.io/gorm"
)

type OrderDetail struct {
	gorm.Model
	OrderID     uint   `json:"order_id"`
	MenuID      uint   `json:"menu_id"`
	Quantity    int    `json:"quantity"`
	MenuOptions string `json:"menu_options"`
	IsReviewed  bool   `json:"is_reviewed"`
	Menu        Menu
}

type OrderDetailResp struct {
	OrderID     uint         `json:"order_id"`
	MenuID      uint         `json:"menu_id"`
	Quantity    int          `json:"quantity"`
	MenuOptions []MenuOption `json:"menu_options"`
}
