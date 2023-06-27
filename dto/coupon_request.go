package dto

type CouponRequest struct {
	IssuerID     uint   `json:"admin_id"`
	Description  string `json:"description" binding:"required"`
	Discount     int    `json:"discount" binding:"required"`
	Availability bool   `json:"availability"`
}
