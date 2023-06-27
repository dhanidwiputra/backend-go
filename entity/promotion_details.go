package entity

import "gorm.io/gorm"

type PromotionDetail struct {
	gorm.Model
	PromotionID uint
	MenuID      uint
	Menu        Menu
}
