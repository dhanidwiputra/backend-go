package entity

import "gorm.io/gorm"

type CategoriesMenu struct {
	gorm.Model
	CategoryID uint `json:"category_id"`
	MenuID     uint `json:"menu_id"`
}
