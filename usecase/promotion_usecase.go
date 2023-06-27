package usecase

import (
	"encoding/json"
	"final-project-backend/domain"
	"final-project-backend/dto"
	"final-project-backend/entity"
	"final-project-backend/repository"
	"time"
)

type PromotionUsecase interface {
	GetPromotions() ([]entity.Promotion, error)
	GetPromotionById(uint) (*entity.Promotion, error)
	CreatePromotion(dto.PromotionFormRequest) (*entity.Promotion, error)
	UpdatePromotion(dto.PromotionFormRequest, uint) (*entity.Promotion, error)
	DeletePromotion(uint) error
	CreatePromotionOrder(dto.PromotionOrderRequest) (*entity.Order, error)
}

type promotionUsecaseImpl struct {
	promotionRepo     repository.PromotionRepository
	paymentOptionRepo repository.PaymentOptionRepository
	mediaUsecase      MediaUsecase
	menuUsecase       MenuUsecase
	cartUsecase       CartUsecase
	couponUsecase     CouponUsecase
	orderUsecase      OrderUsecase
}

type PromotionUsecaseConfig struct {
	PromotionRepo     repository.PromotionRepository
	PaymentOptionRepo repository.PaymentOptionRepository
	MediaUsecase      MediaUsecase
	MenuUsecase       MenuUsecase
	CartUsecase       CartUsecase
	CouponUsecase     CouponUsecase
	OrderUsecase      OrderUsecase
}

func NewPromotionUsecase(c PromotionUsecaseConfig) PromotionUsecase {
	return &promotionUsecaseImpl{
		promotionRepo:     c.PromotionRepo,
		paymentOptionRepo: c.PaymentOptionRepo,
		mediaUsecase:      c.MediaUsecase,
		menuUsecase:       c.MenuUsecase,
		cartUsecase:       c.CartUsecase,
		couponUsecase:     c.CouponUsecase,
		orderUsecase:      c.OrderUsecase,
	}
}

func (u *promotionUsecaseImpl) GetPromotions() ([]entity.Promotion, error) {
	promotions, err := u.promotionRepo.GetPromotions()

	if err != nil {
		return nil, err
	}

	return promotions, nil
}

func (u *promotionUsecaseImpl) GetPromotionById(id uint) (*entity.Promotion, error) {
	promotion, err := u.promotionRepo.GetPromotionById(id)

	if err != nil {
		return nil, err
	}

	return promotion, nil
}

func (u *promotionUsecaseImpl) CreatePromotion(input dto.PromotionFormRequest) (*entity.Promotion, error) {
	var promotionDetails []entity.PromotionDetail

	for _, menuIDs := range input.MenuIDs {
		menu, _ := u.menuUsecase.GetMenuById(menuIDs)
		if menu == nil {
			return nil, domain.ErrMenuNotFound
		}
		promotionDetails = append(promotionDetails, entity.PromotionDetail{
			MenuID: menu.ID,
		})
	}

	var paymentOptions []entity.PaymentRequirement

	for _, paymentIDs := range input.PaymentOptionIDs {
		paymentOption, _ := u.paymentOptionRepo.GetPaymentOptionById(paymentIDs)
		paymentOptions = append(paymentOptions, entity.PaymentRequirement{
			PaymentOptionID: paymentOption.ID,
		})

		if paymentOption == nil {
			return nil, domain.ErrPaymentOptionNotFound
		}
	}

	var promotion = entity.Promotion{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		ExpiredDate: input.ExpiredDate,
	}
	var picUrl, publicId string

	if input.Picture.Size != 0 {
		picUrl, publicId, _ = u.mediaUsecase.FileUpload(input.Picture)
		if picUrl == "" {
			return nil, domain.ErrUploadImage
		}
		promotion.PictureUrl = picUrl
		promotion.PicturePublicID = publicId
	}

	promotion.PaymentRequirements = paymentOptions
	promotion.PromotionDetails = promotionDetails

	promotionRes, err := u.promotionRepo.CreatePromotion(promotion)
	if err != nil {
		return nil, err
	}

	return promotionRes, nil
}

func (u *promotionUsecaseImpl) UpdatePromotion(input dto.PromotionFormRequest, id uint) (*entity.Promotion, error) {
	var promotionDetails []entity.PromotionDetail

	for _, menuIDs := range input.MenuIDs {
		menu, _ := u.menuUsecase.GetMenuById(menuIDs)
		if menu == nil {
			return nil, domain.ErrMenuNotFound
		}
		promotionDetails = append(promotionDetails, entity.PromotionDetail{
			MenuID:      menu.ID,
			PromotionID: id,
		})
	}
	var paymentOptions []entity.PaymentRequirement

	for _, paymentIDs := range input.PaymentOptionIDs {
		paymentOption, _ := u.paymentOptionRepo.GetPaymentOptionById(paymentIDs)
		paymentOptions = append(paymentOptions, entity.PaymentRequirement{
			PaymentOptionID: paymentOption.ID,
			PromotionID:     id,
		})

		if paymentOption == nil {
			return nil, domain.ErrPaymentOptionNotFound
		}
	}

	promotion, err := u.promotionRepo.GetPromotionById(id)
	if promotion == nil {
		return nil, domain.ErrPromotionNotFound
	}

	var picUrl, publicId string

	if input.Picture.Size != 0 {
		picUrl, publicId, _ = u.mediaUsecase.FileUpload(input.Picture)
		promotion.PictureUrl = picUrl
		promotion.PicturePublicID = publicId
	}

	promotion.PaymentRequirements = paymentOptions
	promotion.PromotionDetails = promotionDetails
	promotion.Name = input.Name
	promotion.Description = input.Description
	promotion.Price = input.Price
	promotion.ExpiredDate = input.ExpiredDate

	err = u.promotionRepo.DeletePaymentRequirements(id)
	if err != nil {
		return nil, err
	}

	err = u.promotionRepo.DeletePromotionDetails(id)
	if err != nil {
		return nil, err
	}

	promotionRes, err := u.promotionRepo.UpdatePromotion(*promotion)
	if err != nil {
		return nil, err
	}

	return promotionRes, nil
}

func (u *promotionUsecaseImpl) DeletePromotion(id uint) error {
	promotion, err := u.promotionRepo.GetPromotionById(id)

	if err != nil {
		return err
	}

	if promotion.PicturePublicID != "" {
		u.mediaUsecase.FileDelete(promotion.PicturePublicID)
	}

	err = u.promotionRepo.DeletePromotion(*promotion)

	if err != nil {
		return err
	}

	return nil
}

func (u *promotionUsecaseImpl) CreatePromotionOrder(orderRequest dto.PromotionOrderRequest) (*entity.Order, error) {
	promotion, err := u.promotionRepo.GetPromotionById(orderRequest.PromotionID)
	if promotion == nil {
		return nil, domain.ErrPromotionNotFound
	}
	if err != nil {
		return nil, err
	}

	if promotion.ExpiredDate.Before(time.Now()) {
		return nil, domain.ErrPromotionNotFound
	}

	paymentOption, err := u.paymentOptionRepo.GetPaymentOptionById(orderRequest.PaymentOptionID)
	if paymentOption == nil {
		return nil, domain.ErrPaymentOptionNotFound
	}

	order := entity.Order{
		CouponID:        orderRequest.CouponID,
		PaymentOptionID: orderRequest.PaymentOptionID,
	}

	totalPrice := 0

	for _, orderDetailRequest := range orderRequest.OrderDetailRequest {
		menuOptionStr, _ := json.Marshal(orderDetailRequest.MenuOptions)
		err := u.menuUsecase.IsValidMenuOptions(orderDetailRequest.MenuID, string(menuOptionStr))
		if err != nil {
			return nil, domain.ErrMalformedRequest
		}

		menu, _ := u.menuUsecase.GetMenuById(orderDetailRequest.MenuID)
		if menu == nil {
			return nil, domain.ErrMenuNotFound
		}

		for _, menuOption := range orderDetailRequest.MenuOptions {
			for _, optionList := range menuOption.MenuOptionLists {
				totalPrice += optionList.Price * orderDetailRequest.Quantity
			}
		}
	}

	totalPrice = totalPrice + promotion.Price

	orderedMenus := promotion.Name

	order.TotalPrice = totalPrice
	order.OrderedMenus = orderedMenus

	var coupon *entity.Coupon
	if orderRequest.CouponID != nil {
		coupon, err = u.couponUsecase.GetCouponById(*orderRequest.CouponID)
		if coupon == nil {
			return nil, domain.ErrCouponNotFound
		}

		userCoupon, _ := u.couponUsecase.GetUserCouponByFK(coupon.ID, orderRequest.UserID)
		if userCoupon == nil {
			return nil, domain.ErrUserCouponNotFound
		}
		order.TotalPrice = order.TotalPrice - coupon.Discount
	}

	order.UserID = orderRequest.UserID
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

	orderRes, err := u.orderUsecase.CreateOrderProcess(order, orderDetails, deliveryData)
	if err != nil {
		return nil, err
	}

	return orderRes, nil
}
