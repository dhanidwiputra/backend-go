package entity

import "gorm.io/gorm"

type PaymentRequirement struct {
	gorm.Model
	PromotionID     uint
	PaymentOptionID uint
	PaymentOption   PaymentOption
}
