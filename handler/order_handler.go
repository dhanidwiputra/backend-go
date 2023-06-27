package handler

import (
	"encoding/json"
	"errors"
	"final-project-backend/domain"
	"final-project-backend/dto"
	"final-project-backend/entity"
	"final-project-backend/util"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetAllOrders(c *gin.Context) {
	user, ok := c.MustGet("user").(dto.UserResponse)
	if !ok {
		util.ResponseErrorJSON(c, domain.ErrUnauthorized.Error(), "UNAUTHORIZED", http.StatusUnauthorized)
		return
	}
	query := dto.Query{}
	err := c.ShouldBindQuery(&query)
	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInvalidParams.Error(), "BAD_REQUEST", http.StatusBadRequest)
		return
	}

	orders, err := h.orderUsecase.GetAllOrders(user, query)
	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInvalidRequest.Error(), "BAD_REQUEST", http.StatusBadRequest)
		return
	}

	util.ResponseSuccesJSON(c, orders, http.StatusOK)
}

func (h *Handler) GetAllPaymentOptions(c *gin.Context) {
	paymentOptions, err := h.orderUsecase.GetAllPaymentOptions()
	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInternalServer.Error(), "INTERNAL_SERVER_ERROR", http.StatusInternalServerError)
		return
	}

	util.ResponseSuccesJSON(c, paymentOptions, http.StatusOK)
}

func (h *Handler) CreateOrder(c *gin.Context) {
	user := c.MustGet("user").(dto.UserResponse)
	var orderRequest dto.OrderRequest
	var err error

	if err = c.ShouldBindJSON(&orderRequest); err != nil {
		util.ResponseErrorJSON(c, domain.ErrInvalidBody.Error(), "BAD_REQUEST", http.StatusBadRequest)
		return
	}

	order := entity.Order{
		CouponID:        orderRequest.CouponID,
		PaymentOptionID: orderRequest.PaymentOptionID,
	}

	totalPrice := 0
	var orderedMenusArr []string
	for _, orderDetailRequest := range orderRequest.OrderDetailRequest {
		// menuOptionStr, _ := json.Marshal(orderDetailRequest.MenuOptions)
		// err := h.menuUsecase.IsValidMenuOptions(orderDetailRequest.MenuID, string(menuOptionStr))
		// if err != nil {
		// 	util.ResponseErrorJSON(c, domain.ErrMalformedRequest.Error(), "BAD_REQUEST", http.StatusBadRequest)
		// 	return
		// }

		menuPrice := 0
		menu, _ := h.menuUsecase.GetMenuById(orderDetailRequest.MenuID)
		if menu == nil {
			util.ResponseErrorJSON(c, domain.ErrMenuNotFound.Error(), "MENU_NOT_FOUND", http.StatusNotFound)
			return
		}
		menuPrice += menu.Price * orderDetailRequest.Quantity
		orderedMenusArr = append(orderedMenusArr, menu.Name)
		for _, menuOption := range orderDetailRequest.MenuOptions {
			for _, optionList := range menuOption.MenuOptionLists {
				if optionList.Checked {
					menuPrice += optionList.Price * orderDetailRequest.Quantity
				}
			}
		}
		totalPrice += menuPrice
	}

	orderedMenusArr = util.UniqueString(orderedMenusArr)
	orderedMenus := strings.Join(orderedMenusArr, ",")

	order.TotalPrice = totalPrice
	order.OrderedMenus = orderedMenus

	var coupon *entity.Coupon
	if orderRequest.CouponID != nil {
		coupon, err = h.couponUsecase.GetCouponById(*orderRequest.CouponID)
		if coupon == nil {
			util.ResponseErrorJSON(c, domain.ErrCouponNotFound.Error(), "COUPON_NOT_FOUND", http.StatusNotFound)
			return
		}

		userCoupon, _ := h.couponUsecase.GetUserCouponByFK(*orderRequest.CouponID, user.ID)
		if userCoupon == nil {
			util.ResponseErrorJSON(c, domain.ErrUserCouponNotFound.Error(), "USER_COUPON_NOT_FOUND", http.StatusNotFound)
			return
		}
		userCoupon.Stock = userCoupon.Stock - 1
		h.couponUsecase.UpdateUserCoupon(*userCoupon)
		order.TotalPrice = order.TotalPrice - coupon.Discount
	}

	order.UserID = user.ID
	order.OrderDate = time.Now()

	if order.TotalPrice < 0 {
		order.TotalPrice = 0
	}

	var orderDetails []entity.OrderDetail
	for _, orderDetailRequest := range orderRequest.OrderDetailRequest {
		menuOptionStr, _ := json.Marshal(orderDetailRequest.MenuOptions)
		orderDetail := entity.OrderDetail{
			MenuID:      orderDetailRequest.MenuID,
			Quantity:    orderDetailRequest.Quantity,
			MenuOptions: string(menuOptionStr),
		}
		orderDetails = append(orderDetails, orderDetail)
	}

	deliveryData := entity.Delivery{
		Address: orderRequest.DeliveryAddress,
		Status:  "pending",
	}

	orderRes, err := h.orderUsecase.CreateOrderProcess(order, orderDetails, deliveryData)
	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInternalServer.Error(), "INTERNAL_SERVER_ERROR", http.StatusInternalServerError)
		return
	}

	err = h.cartUsecase.EmptyCart(int(user.ID))
	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInternalServer.Error(), "INTERNAL_SERVER_ERROR", http.StatusInternalServerError)
		return
	}

	util.ResponseSuccesJSON(c, orderRes, http.StatusCreated)
}

func (h *Handler) GetOrderById(c *gin.Context) {
	id := c.Param("id")
	uintId64, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrMalformedRequest.Error(), "BAD_REQUEST", http.StatusBadRequest)
		return
	}

	order, err := h.orderUsecase.GetOrderByID(uint(uintId64))
	if errors.Is(err, domain.ErrOrderNotFound) {
		util.ResponseErrorJSON(c, domain.ErrOrderNotFound.Error(), "ORDER_NOT_FOUND", http.StatusNotFound)
		return
	}
	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInternalServer.Error(), "INTERNAL_SERVER_ERROR", http.StatusInternalServerError)
		return
	}

	util.ResponseSuccesJSON(c, order, http.StatusOK)
}

func (h *Handler) CreateCustomerReview(c *gin.Context) {
	user, ok := c.MustGet("user").(dto.UserResponse)
	if !ok {
		util.ResponseErrorJSON(c, domain.ErrUnauthorized.Error(), "UNAUTHORIZED", http.StatusUnauthorized)
		return
	}

	var customerReviewRequest dto.CustomerReviewRequest
	var err error

	if err = c.ShouldBindJSON(&customerReviewRequest); err != nil {
		util.ResponseErrorJSON(c, domain.ErrInvalidBody.Error(), "BAD_REQUEST", http.StatusBadRequest)
		return
	}

	customerReviewRequest.UserID = user.ID

	customerReviewRes, err := h.orderUsecase.CreateCustomerReview(customerReviewRequest)
	if errors.Is(err, domain.ErrDuplicateReview) {
		util.ResponseErrorJSON(c, domain.ErrDuplicateReview.Error(), "DUPLICATE_REVIEW", http.StatusConflict)
		return
	}
	if errors.Is(err, domain.ErrUnauthorized) {
		util.ResponseErrorJSON(c, domain.ErrUnauthorized.Error(), "UNAUTHORIZED", http.StatusUnauthorized)
		return
	}
	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInternalServer.Error(), "INTERNAL_SERVER_ERROR", http.StatusInternalServerError)
		return
	}

	util.ResponseSuccesJSON(c, customerReviewRes, http.StatusCreated)
}

func (h *Handler) GetCustomerReviewsByMenuId(c *gin.Context) {
	id := c.Param("id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrMalformedRequest.Error(), "BAD_REQUEST", http.StatusBadRequest)
		return
	}

	menu, err := h.menuUsecase.GetMenuById(uint(intId))
	if errors.Is(err, domain.ErrMenuNotFound) {
		util.ResponseErrorJSON(c, domain.ErrMenuNotFound.Error(), "MENU_NOT_FOUND", http.StatusNotFound)
		return
	}

	customerReviews, err := h.orderUsecase.GetCustomerReviewsByMenuId(menu.ID)
	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInternalServer.Error(), "INTERNAL_SERVER_ERROR", http.StatusInternalServerError)
		return
	}

	util.ResponseSuccesJSON(c, customerReviews, http.StatusOK)
}

func (h *Handler) GetTransactionTotalByDate(c *gin.Context) {
	days := c.Query("d")
	intDays, err := strconv.Atoi(days)
	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrMalformedRequest.Error(), "BAD_REQUEST", http.StatusBadRequest)
		return
	}

	dateBefore := time.Now().AddDate(0, 0, -intDays)

	transactionTotal, err := h.orderUsecase.GetTransactionTotalByDate(dateBefore)
	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInternalServer.Error(), "INTERNAL_SERVER_ERROR", http.StatusInternalServerError)
		return
	}

	util.ResponseSuccesJSON(c, transactionTotal, http.StatusOK)
}
