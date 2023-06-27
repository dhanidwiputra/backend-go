package dto

type PromotionOrderRequest struct {
	PromotionID        uint                  `json:"promotion_id"`
	CouponID           *uint                 `json:"coupon_id,omitempty"`
	PaymentOptionID    uint                  `json:"payment_option_id" binding:"required"`
	TotalPrice         int                   `json:"total_price"`
	DeliveryAddress    string                `json:"delivery_address" binding:"required"`
	OrderDetailRequest []*OrderDetailRequest `json:"order_detail_request" binding:"required"`
	UserID             uint                  `json:"user_id"`
}
