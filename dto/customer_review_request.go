package dto

type CustomerReviewRequest struct {
	OrderDetailID uint   `json:"order_detail_id" binding:"required"`
	Review        string `json:"review" binding:"required"`
	Rating        int    `json:"rating" binding:"required,min=1,max=5"`
	UserID        uint   `json:"user_id"`
}
