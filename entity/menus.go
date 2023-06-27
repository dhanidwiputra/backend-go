package entity

import (
	"gorm.io/gorm"
)

type Menu struct {
	gorm.Model
	Name            string     `json:"name"`
	Description     string     `json:"description"`
	Price           int        `json:"price"`
	PictureUrl      string     `json:"picture_url"`
	PicturePublicId string     `json:"picture_public_id"`
	MenuOptions     string     `json:"menu_options"`
	AvgRating       float64    `json:"avg_rating"`
	UserRatingCount int        `json:"user_rating_count"`
	Categories      []Category `gorm:"many2many:categories_menus;" json:"categories"`
}
