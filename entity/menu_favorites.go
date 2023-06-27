package entity

import "gorm.io/gorm"

type MenuFavorite struct {
	gorm.Model
	MenuID uint `json:"menu_id"`
	UserID uint `json:"user_id"`
}
