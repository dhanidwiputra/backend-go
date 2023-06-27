package entity

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name string `json:"name"`
	Menu []Menu `gorm:"many2many:categories_menus;" json:"menu,omitempty"`
}
