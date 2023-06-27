package entity

import "gorm.io/gorm"

type PaymentOption struct {
	gorm.Model
	Name        string `json:"name"`
	Description string `json:"description"`
}
