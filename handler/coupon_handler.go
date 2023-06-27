package handler

import (
	"final-project-backend/domain"
	"final-project-backend/dto"
	"final-project-backend/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateCoupon(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		util.ResponseErrorJSON(c, domain.ErrUnauthorized.Error(), "UNAUTHORIZED", 401)
		return
	}

	userId := user.(dto.UserResponse).ID

	couponRequestBody := dto.CouponRequest{}
	err := c.BindJSON(&couponRequestBody)
	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInvalidBody.Error(), "INVALID_BODY_REQUEST", 400)
		return
	}

	couponRequestBody.IssuerID = userId

	coupon, err := h.couponUsecase.CreateCoupon(couponRequestBody)
	if err != nil {
		util.ResponseErrorJSON(c, err.Error(), "INTERNAL_SERVER_ERROR", http.StatusInternalServerError)
		return
	}

	util.ResponseSuccesJSON(c, coupon, http.StatusOK)
}

func (h *Handler) GetCoupons(c *gin.Context) {
	coupons, err := h.couponUsecase.GetCoupons()
	if err != nil {
		util.ResponseErrorJSON(c, err.Error(), "INTERNAL_SERVER_ERROR", http.StatusInternalServerError)
		return
	}

	util.ResponseSuccesJSON(c, coupons, http.StatusOK)
}

func (h *Handler) GetUserCoupons(c *gin.Context) {
	user := c.MustGet("user").(dto.UserResponse)

	userCoupons, err := h.couponUsecase.GetCouponsByUserId(user.ID)
	if err != nil {
		util.ResponseErrorJSON(c, err.Error(), "INTERNAL_SERVER_ERROR", http.StatusInternalServerError)
		return
	}

	util.ResponseSuccesJSON(c, userCoupons, http.StatusOK)
}

func (h *Handler) UpdateCoupon(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInvalidBody.Error(), "INVALID_BODY_REQUEST", 400)
		return
	}

	couponRequestBody := dto.CouponRequest{}
	err = c.BindJSON(&couponRequestBody)
	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInvalidBody.Error(), "INVALID_BODY_REQUEST", 400)
		return
	}

	user := c.MustGet("user").(dto.UserResponse)

	coupon, err := h.couponUsecase.GetCouponById(uint(idInt))
	if coupon == nil {
		util.ResponseErrorJSON(c, domain.ErrCouponNotFound.Error(), "COUPON_NOT_FOUND", 404)
		return
	}

	if coupon.IssuerID != user.ID {
		util.ResponseErrorJSON(c, domain.ErrUnauthorized.Error(), "UNAUTHORIZED", 401)
		return
	}

	coupon.Discount = couponRequestBody.Discount
	coupon.Description = couponRequestBody.Description
	coupon.Availability = couponRequestBody.Availability

	couponRes, err := h.couponUsecase.UpdateCoupon(*coupon)
	if err != nil {
		util.ResponseErrorJSON(c, err.Error(), "INTERNAL_SERVER_ERROR", http.StatusInternalServerError)
		return
	}

	util.ResponseSuccesJSON(c, couponRes, http.StatusOK)
}

func (h *Handler) DeleteCouponById(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInvalidBody.Error(), "INVALID_BODY_REQUEST", 400)
		return
	}

	coupon, err := h.couponUsecase.GetCouponById(uint(idInt))
	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrCouponNotFound.Error(), "COUPON_NOT_FOUND", 404)
		return
	}

	err = h.couponUsecase.DeleteCoupon(*coupon)
	if err != nil {
		util.ResponseErrorJSON(c, err.Error(), "INTERNAL_SERVER_ERROR", http.StatusInternalServerError)
		return
	}

	util.ResponseSuccesJSON(c, nil, http.StatusNoContent)
}
