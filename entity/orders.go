package entity

import (
	"time"

	gorm "gorm.io/gorm"
)

type Order struct {
	gorm.Model
	OrderDate       time.Time `json:"order_date"`
	CouponID        *uint     `json:"coupon_id,omitempty"`
	PaymentOptionID uint      `json:"payment_option_id"`
	OrderedMenus    string    `json:"ordered_menus"`
	OrderDetails    []OrderDetail
	TotalPrice      int `json:"total_price"`
	Delivery        Delivery
	UserID          uint `json:"user_id"`
}
