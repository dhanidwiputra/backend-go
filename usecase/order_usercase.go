package usecase

import (
	"final-project-backend/domain"
	"final-project-backend/dto"
	"final-project-backend/entity"
	"final-project-backend/repository"
	"time"
)

type OrderUsecase interface {
	GetAllOrders(dto.UserResponse, dto.Query) ([]entity.Order, error)
	GetOrderByID(id uint) (*entity.Order, error)
	GetAllPaymentOptions() ([]entity.PaymentOption, error)
	CreateOrder(entity.Order) (*entity.Order, error)
	CreateOrderProcess(entity.Order, []entity.OrderDetail, entity.Delivery) (*entity.Order, error)
	CreateBatchOrderDetails([]entity.OrderDetail) (*[]entity.OrderDetail, error)
	CreateCustomerReview(dto.CustomerReviewRequest) (*entity.CustomerReview, error)
	GetCustomerReviewsByMenuId(id uint) (*[]entity.CustomerReview, error)
	GetTransactionTotalByDate(date time.Time) (int64, error)
}

type orderUsecaseImpl struct {
	orderRepo      repository.OrderRepository
	cartRepo       repository.CartRepository
	menuRepo       repository.MenuRepository
	paymentOptRepo repository.PaymentOptionRepository
	deliveryRepo   repository.DeliveryRepository
}

type OrderUsecaseConfig struct {
	OrderRepo      repository.OrderRepository
	CartRepo       repository.CartRepository
	MenuRepo       repository.MenuRepository
	PaymentOptRepo repository.PaymentOptionRepository
	DeliveryRepo   repository.DeliveryRepository
}

func NewOrderUsecase(c OrderUsecaseConfig) OrderUsecase {
	return &orderUsecaseImpl{
		orderRepo:      c.OrderRepo,
		cartRepo:       c.CartRepo,
		menuRepo:       c.MenuRepo,
		paymentOptRepo: c.PaymentOptRepo,
		deliveryRepo:   c.DeliveryRepo,
	}
}

func (o *orderUsecaseImpl) GetAllOrders(user dto.UserResponse, query dto.Query) ([]entity.Order, error) {
	return o.orderRepo.GetAllOrders(user, query)
}

func (o *orderUsecaseImpl) GetOrderByID(id uint) (*entity.Order, error) {
	order, err := o.orderRepo.GetOrderByID(id)
	if order == nil {
		return nil, domain.ErrOrderNotFound
	}

	if err != nil {
		return nil, err
	}

	return order, nil
}

func (o *orderUsecaseImpl) GetAllPaymentOptions() ([]entity.PaymentOption, error) {
	return o.paymentOptRepo.GetAllPaymentOptions()
}

func (o *orderUsecaseImpl) CreateOrder(b entity.Order) (*entity.Order, error) {
	order, err := o.orderRepo.CreateOrder(b)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (o *orderUsecaseImpl) CreateOrderProcess(order entity.Order, details []entity.OrderDetail, delivery entity.Delivery) (*entity.Order, error) {
	orderRes, err := o.orderRepo.CreateOrderProcess(order, details, delivery)
	if err != nil {
		return nil, err
	}

	return orderRes, nil
}

func (o *orderUsecaseImpl) CreateBatchOrderDetails(b []entity.OrderDetail) (*[]entity.OrderDetail, error) {
	orderDetails, err := o.orderRepo.CreateBatchOrderDetails(b)
	if err != nil {
		return nil, err
	}

	return orderDetails, nil
}

func (o *orderUsecaseImpl) CreateCustomerReview(input dto.CustomerReviewRequest) (*entity.CustomerReview, error) {
	orderDetail, err := o.orderRepo.GetOrderDetailsByID(input.OrderDetailID)
	if err != nil {
		return nil, err
	}

	if orderDetail.IsReviewed {
		return nil, domain.ErrDuplicateReview
	}

	order, err := o.orderRepo.GetOrderByID(orderDetail.OrderID)
	if order.UserID != input.UserID {
		return nil, domain.ErrUnauthorized
	}

	customerReview := entity.CustomerReview{
		OrderDetailID: input.OrderDetailID,
		Rating:        input.Rating,
		Review:        input.Review,
		UserID:        input.UserID,
	}
	customerReviewRes, err := o.orderRepo.CreateCustomerReviewProcess(customerReview)
	if err != nil {
		return nil, err
	}

	return customerReviewRes, nil
}

func (o *orderUsecaseImpl) GetCustomerReviewsByMenuId(id uint) (*[]entity.CustomerReview, error) {
	return o.orderRepo.GetCustomerReviewsByMenuId(id)
}

func (o *orderUsecaseImpl) GetTransactionTotalByDate(date time.Time) (int64, error) {
	return o.orderRepo.GetTransactionTotalByDate(date)
}
