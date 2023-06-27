package handler

import (
	"final-project-backend/usecase"
)

type Handler struct {
	authUsecase      usecase.AuthUsecase
	userUsecase      usecase.UserUsecase
	couponUsecase    usecase.CouponUsecase
	menuUsecase      usecase.MenuUsecase
	mediaUsecase     usecase.MediaUsecase
	cartUsecase      usecase.CartUsecase
	orderUsecase     usecase.OrderUsecase
	deliveryUsecase  usecase.DeliveryUsecase
	gameUsecase      usecase.GameUsecase
	promotionUsecase usecase.PromotionUsecase
}

type HandlerConfig struct {
	AuthUsecase      usecase.AuthUsecase
	UserUsecase      usecase.UserUsecase
	CouponUsecase    usecase.CouponUsecase
	MenuUsecase      usecase.MenuUsecase
	MediaUsecase     usecase.MediaUsecase
	CartUsecase      usecase.CartUsecase
	OrderUsecase     usecase.OrderUsecase
	DeliveryUsecase  usecase.DeliveryUsecase
	GameUsecase      usecase.GameUsecase
	PromotionUsecase usecase.PromotionUsecase
}

func New(c HandlerConfig) *Handler {
	return &Handler{
		authUsecase:      c.AuthUsecase,
		userUsecase:      c.UserUsecase,
		couponUsecase:    c.CouponUsecase,
		menuUsecase:      c.MenuUsecase,
		mediaUsecase:     c.MediaUsecase,
		cartUsecase:      c.CartUsecase,
		orderUsecase:     c.OrderUsecase,
		deliveryUsecase:  c.DeliveryUsecase,
		gameUsecase:      c.GameUsecase,
		promotionUsecase: c.PromotionUsecase,
	}
}
