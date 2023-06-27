package entity

import "gorm.io/gorm"

type Game struct {
	gorm.Model
	UserID     uint   `json:"user_id"`
	Score      *int   `json:"score"`
	CouponID   *uint  `json:"coupon_id"`
	Answer     string `json:"answer"`
	Difficulty string `json:"difficulty"`
}
