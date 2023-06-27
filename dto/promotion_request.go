package dto

import (
	"mime/multipart"
	"time"
)

type PromotionFormRequest struct {
	Name             string               `form:"name" binding:"required"`
	Description      string               `form:"description" binding:"required"`
	Price            int                  `form:"price" binding:"required"`
	Picture          multipart.FileHeader `form:"picture"`
	ExpiredDate      time.Time            `form:"expired_date" binding:"required"`
	PaymentOptionIDs []uint               `form:"payment_option_ids" binding:"required"`
	MenuIDs          []uint               `form:"menu_ids" binding:"required"`
}
