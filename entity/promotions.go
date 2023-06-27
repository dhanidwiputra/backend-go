package entity

import (
	"time"

	"gorm.io/gorm"
)

type Promotion struct {
	gorm.Model
	Name                string               `json:"name"`
	Description         string               `json:"description"`
	Price               int                  `json:"price"`
	ExpiredDate         time.Time            `json:"expired_date"`
	PictureUrl          string               `json:"picture_url"`
	PicturePublicID     string               `json:"picture_public_id"`
	PaymentRequirements []PaymentRequirement `json:"payment_requirements"`
	PromotionDetails    []PromotionDetail    `json:"promotion_details"`
}
