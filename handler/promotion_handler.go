package handler

import (
	"errors"
	"final-project-backend/domain"
	"final-project-backend/dto"
	"final-project-backend/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func (h *Handler) GetPromotions(c *gin.Context) {
	promotions, err := h.promotionUsecase.GetPromotions()
	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInternalServer.Error(), "INTERNAL_SERVER_ERROR", http.StatusInternalServerError)
		return
	}

	util.ResponseSuccesJSON(c, promotions, http.StatusOK)
}

func (h *Handler) GetPromotionById(c *gin.Context) {
	id := c.Param("id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInvalidBody.Error(), "INVALID_BODY_REQUEST", http.StatusBadRequest)
		return
	}

	promotion, err := h.promotionUsecase.GetPromotionById(uint(intId))
	if promotion == nil {
		util.ResponseErrorJSON(c, domain.ErrResourceNotFound.Error(), "NOT_FOUND", http.StatusNotFound)
		return
	}
	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInternalServer.Error(), "INTERNAL_SERVER_ERROR", http.StatusInternalServerError)
	}

	util.ResponseSuccesJSON(c, promotion, http.StatusOK)
}

func (h *Handler) CreatePromotion(c *gin.Context) {
	var input dto.PromotionFormRequest
	err := c.MustBindWith(&input, binding.FormMultipart)
	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInvalidBody.Error(), "INVALID_BODY_REQUEST", http.StatusBadRequest)
		return
	}

	promotion, err := h.promotionUsecase.CreatePromotion(input)

	if errors.Is(err, domain.ErrMenuNotFound) {
		util.ResponseErrorJSON(c, domain.ErrMenuNotFound.Error(), "MENU_NOT_FOUND", http.StatusBadRequest)
		return
	}
	if errors.Is(err, domain.ErrPaymentOptionNotFound) {
		util.ResponseErrorJSON(c, domain.ErrPaymentOptionNotFound.Error(), "PAYMENT_OPTION_NOT_FOUND", http.StatusBadRequest)
		return
	}

	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInternalServer.Error(), "INSERT_MENU_FAILED", http.StatusInternalServerError)
		return
	}

	util.ResponseSuccesJSON(c, promotion, http.StatusOK)
}

func (h *Handler) UpdatePromotion(c *gin.Context) {
	id := c.Param("id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInvalidBody.Error(), "INVALID_BODY_REQUEST", http.StatusBadRequest)
		return
	}

	var input dto.PromotionFormRequest
	err = c.MustBindWith(&input, binding.FormMultipart)
	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInvalidBody.Error(), "INVALID_BODY_REQUEST", http.StatusBadRequest)
		return
	}

	promotion, err := h.promotionUsecase.UpdatePromotion(input, uint(intId))
	if errors.Is(err, domain.ErrMenuNotFound) {
		util.ResponseErrorJSON(c, domain.ErrMenuNotFound.Error(), "MENU_NOT_FOUND", http.StatusBadRequest)
		return
	}
	if errors.Is(err, domain.ErrPaymentOptionNotFound) {
		util.ResponseErrorJSON(c, domain.ErrPaymentOptionNotFound.Error(), "PAYMENT_OPTION_NOT_FOUND", http.StatusBadRequest)
		return
	}

	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInternalServer.Error(), "UPDATE_MENU_FAILED", http.StatusInternalServerError)
		return
	}

	util.ResponseSuccesJSON(c, promotion, http.StatusOK)
}

func (h *Handler) DeletePromotion(c *gin.Context) {
	id := c.Param("id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInvalidBody.Error(), "INVALID_BODY_REQUEST", http.StatusBadRequest)
		return
	}

	err = h.promotionUsecase.DeletePromotion(uint(intId))
	if errors.Is(err, domain.ErrPromotionNotFound) {
		util.ResponseErrorJSON(c, domain.ErrPromotionNotFound.Error(), "PROMOTION_NOT_FOUND", http.StatusNotFound)
		return
	}
	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInternalServer.Error(), "DELETE_MENU_FAILED", http.StatusInternalServerError)
		return
	}

	util.ResponseSuccesJSON(c, nil, http.StatusOK)
}

func (h *Handler) CreatePromotionOrder(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		util.ResponseErrorJSON(c, domain.ErrUnauthorized.Error(), "UNAUTHORIZED", http.StatusUnauthorized)
		return
	}

	var input dto.PromotionOrderRequest
	err := c.ShouldBindJSON(&input)
	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInvalidBody.Error(), "INVALID_BODY_REQUEST", http.StatusBadRequest)
		return
	}

	promotionId := c.Param("id")
	promotionIdInt, err := strconv.Atoi(promotionId)
	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInvalidBody.Error(), "INVALID_BODY_REQUEST", http.StatusBadRequest)
		return
	}

	input.PromotionID = uint(promotionIdInt)

	input.UserID = user.(dto.UserResponse).ID

	promotionOrder, err := h.promotionUsecase.CreatePromotionOrder(input)

	if errors.Is(err, domain.ErrPromotionNotFound) {
		util.ResponseErrorJSON(c, domain.ErrPromotionNotFound.Error(), "PROMOTION_NOT_FOUND", http.StatusBadRequest)
		return
	}

	if errors.Is(err, domain.ErrPaymentOptionNotFound) {
		util.ResponseErrorJSON(c, domain.ErrPaymentOptionNotFound.Error(), "PAYMENT_OPTION_NOT_FOUND", http.StatusBadRequest)
		return
	}

	if errors.Is(err, domain.ErrCouponNotFound) {
		util.ResponseErrorJSON(c, domain.ErrCouponNotFound.Error(), "COUPON_NOT_FOUND", http.StatusBadRequest)
		return
	}

	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInternalServer.Error(), "INSERT_MENU_FAILED", http.StatusInternalServerError)
		return
	}

	util.ResponseSuccesJSON(c, promotionOrder, http.StatusOK)
}
