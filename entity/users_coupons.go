package entity

import "gorm.io/gorm"

type UsersCoupon struct {
	gorm.Model
	UserID   uint `json:"user_id"`
	CouponID uint `json:"coupon_id"`
	Stock    int  `json:"stock"`
	Coupon   Coupon
}
