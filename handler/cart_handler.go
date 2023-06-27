package handler

import (
	"encoding/json"
	"errors"
	"final-project-backend/domain"
	"final-project-backend/dto"
	"final-project-backend/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetCartItems(c *gin.Context) {
	user, ok := c.Get("user")

	if !ok {
		util.ResponseErrorJSON(c, domain.ErrUnauthorized.Error(), "UNAUTHORIZED", 401)
		return
	}

	userID := user.(dto.UserResponse).ID

	cartItems, err := h.cartUsecase.GetCartItems(userID)
	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInternalServer.Error(), "INTERNAL_SERVER_ERROR", http.StatusInternalServerError)
		return
	}

	util.ResponseSuccesJSON(c, cartItems, http.StatusOK)
}

func (h *Handler) AddToCart(c *gin.Context) {
	var item dto.CartItemRequest
	if err := c.ShouldBindJSON(&item); err != nil {
		util.ResponseErrorJSON(c, domain.ErrInvalidBody.Error(), "INVALID_BODY_REQUEST", http.StatusBadRequest)
		return
	}

	user, ok := c.Get("user")
	if !ok {
		util.ResponseErrorJSON(c, domain.ErrUnauthorized.Error(), "UNAUTHORIZED", 401)
		return
	}

	userID := user.(dto.UserResponse).ID

	menuOptionsString, err := json.Marshal(item.MenuOptions)
	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInvalidBody.Error(), "INVALID_BODY_REQUEST", http.StatusBadRequest)
		return
	}

	menu, err := h.menuUsecase.GetMenuById(item.MenuID)
	if menu == nil {
		util.ResponseErrorJSON(c, domain.ErrMenuNotFound.Error(), "MENU_NOT_FOUND", http.StatusNotFound)
		return
	}

	cartItemData := dto.CartItemData{
		UserID:      userID,
		MenuID:      item.MenuID,
		Quantity:    item.Quantity,
		MenuOptions: string(menuOptionsString),
	}

	cartRes, err := h.cartUsecase.AddToCart(cartItemData)
	if err != nil {
		util.ResponseErrorJSON(c, err.Error(), "INVALID_BODY_REQUEST", http.StatusBadRequest)
		return
	}

	util.ResponseSuccesJSON(c, cartRes, http.StatusCreated)
}

func (h *Handler) UpdateCartItem(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		util.ResponseErrorJSON(c, domain.ErrUnauthorized.Error(), "UNAUTHORIZED", 401)
		return
	}

	userID := user.(dto.UserResponse).ID

	var cartItemUpdateRequest dto.CartItemUpdateRequest
	if err := c.ShouldBindJSON(&cartItemUpdateRequest); err != nil {
		util.ResponseErrorJSON(c, domain.ErrInvalidBody.Error(), "INVALID_BODY_REQUEST", http.StatusBadRequest)
		return
	}

	cartItemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInvalidBody.Error(), "INVALID_BODY_REQUEST", http.StatusBadRequest)
		return
	}

	cartItem, err := h.cartUsecase.GetCartItemById(uint(cartItemId))
	if errors.Is(err, domain.ErrCartItemNotFound) {
		util.ResponseErrorJSON(c, domain.ErrCartItemNotFound.Error(), "ITEM_NOT_FOUND", http.StatusNotFound)
		return
	}

	if cartItem.UserID != userID {
		util.ResponseErrorJSON(c, domain.ErrUnauthorized.Error(), "UNAUTHORIZED", 401)
		return
	}
	menuOptionsString, err := json.Marshal(cartItemUpdateRequest.MenuOptions)
	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInvalidBody.Error(), "INVALID_BODY_REQUEST", http.StatusBadRequest)
		return
	}

	cartItem.Quantity = cartItemUpdateRequest.Quantity
	cartItem.MenuOptions = string(menuOptionsString)

	cartItem, err = h.cartUsecase.UpdateCartItem(*cartItem)
	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrUpdateCart.Error(), "UPDATE_CART_FAILED", http.StatusInternalServerError)
		return
	}

	util.ResponseSuccesJSON(c, cartItem, http.StatusOK)
}

func (h *Handler) DeleteCartItem(c *gin.Context) {
	id := c.Param("id")
	idInt64, err := strconv.ParseInt(id, 10, 64)
	idUint := uint(idInt64)

	err = h.cartUsecase.DeleteCartItem(idUint)
	if errors.Is(err, domain.ErrCartItemNotFound) {
		util.ResponseErrorJSON(c, domain.ErrCartItemNotFound.Error(), "ITEM_NOT_FOUND", http.StatusNotFound)
		return
	}

	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInternalServer.Error(), "INTERNAL_SERVER_ERROR", http.StatusInternalServerError)
		return
	}

	util.ResponseSuccesJSON(c, nil, http.StatusNoContent)
}

func (h *Handler) EmptyCart(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		util.ResponseErrorJSON(c, domain.ErrUnauthorized.Error(), "UNAUTHORIZED", 401)
		return
	}

	userID := user.(dto.UserResponse).ID

	err := h.cartUsecase.EmptyCart(int(userID))
	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInternalServer.Error(), "INTERNAL_SERVER_ERROR", http.StatusInternalServerError)
		return
	}

	util.ResponseSuccesJSON(c, nil, http.StatusNoContent)
}
