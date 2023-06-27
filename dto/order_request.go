package dto

import "final-project-backend/entity"

type OrderRequest struct {
	CouponID           *uint                 `json:"coupon_id,omitempty"`
	PaymentOptionID    uint                  `json:"payment_option_id" binding:"required"`
	TotalPrice         int                   `json:"total_price"`
	DeliveryAddress    string                `json:"delivery_address" binding:"required"`
	OrderDetailRequest []*OrderDetailRequest `json:"order_detail_request" binding:"required"`
}

type OrderDetailRequest struct {
	MenuID      uint                `json:"menu_id" binding:"required"`
	Quantity    int                 `json:"quantity" binding:"required"`
	MenuOptions []entity.MenuOption `json:"menu_options" binding:"required"`
}

type OrderTotalRerquest struct {
	Date string `json:"date" binding:"required"`
}
