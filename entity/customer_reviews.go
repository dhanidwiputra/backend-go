package entity

import "gorm.io/gorm"

type CustomerReview struct {
	gorm.Model
	OrderDetailID uint   `json:"order_detail_id"`
	UserID        uint   `json:"user_id"`
	Review        string `json:"review"`
	Rating        int    `json:"rating"`
}
