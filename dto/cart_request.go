package dto

import "final-project-backend/entity"

type CartItemRequest struct {
	MenuID      uint                `json:"menu_id"`
	Quantity    int                 `json:"quantity"`
	UserID      uint                `json:"user_id"`
	MenuOptions []entity.MenuOption `json:"menu_options" binding:"required"`
}

type CartItemData struct {
	MenuID      uint   `json:"menu_id" `
	Quantity    int    `json:"quantity" binding:"required"`
	UserID      uint   `json:"user_id"`
	MenuOptions string `json:"menu_options" binding:"required"`
}

type CartItemUpdateRequest struct {
	Quantity    int                 `json:"quantity"`
	MenuOptions []entity.MenuOption `json:"menu_options" binding:"required"`
}
