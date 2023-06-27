package repository

import (
	"final-project-backend/domain"
	"final-project-backend/dto"
	"final-project-backend/entity"
	"time"

	"gorm.io/gorm"
)

type OrderRepository interface {
	GetAllOrders(user dto.UserResponse, query dto.Query) ([]entity.Order, error)
	GetOrderByID(id uint) (*entity.Order, error)
	GetOrderDetailsByID(id uint) (*entity.OrderDetail, error)
	GetTransactionTotalByDate(date time.Time) (int64, error)
	CreateOrder(order entity.Order) (*entity.Order, error)
	CreateOrderProcess(order entity.Order, details []entity.OrderDetail, delivery entity.Delivery) (*entity.Order, error)
	CreateBatchOrderDetails([]entity.OrderDetail) (*[]entity.OrderDetail, error)
	UserHasOrdered(userId, menuId uint) (bool, error)
	UserHasReviewed(userId uint, orderDetailId uint) (bool, error)
	GetCustomerReviewsByMenuId(id uint) (*[]entity.CustomerReview, error)
	CreateCustomerReview(customerReview entity.CustomerReview) (*entity.CustomerReview, error)
	CreateCustomerReviewProcess(customerReview entity.CustomerReview) (*entity.CustomerReview, error)
}

type orderRepositoryImpl struct {
	db *gorm.DB
}

type OrderRepoConfig struct {
	DB *gorm.DB
}

func NewOrderRepository(c OrderRepoConfig) OrderRepository {
	return &orderRepositoryImpl{db: c.DB}
}

func (o *orderRepositoryImpl) GetAllOrders(user dto.UserResponse, query dto.Query) ([]entity.Order, error) {
	var orders []entity.Order
	var err error
	if user.Role == "admin" {
		err = o.db.Preload("OrderDetails.Menu.Categories").Preload("OrderDetails.Menu").Preload("OrderDetails").Preload("Delivery").
			Where("ordered_menus ILIKE ?", "%"+query.Search+"%").Order(query.SortBy + " " + query.Sort).Find(&orders).Error
	}

	if user.Role == "user" {
		err = o.db.Preload("OrderDetails.Menu.Categories").Preload("OrderDetails.Menu").Preload("OrderDetails").Preload("Delivery").
			Where("ordered_menus ILIKE ? AND user_id = ?", "%"+query.Search+"%", user.ID).Order(query.SortBy + " " + query.Sort).Find(&orders).Error
	}

	if err != nil {
		return nil, domain.ErrInvalidRequest
	}

	return orders, nil
}

func (o *orderRepositoryImpl) GetTransactionTotalByDate(dateStart time.Time) (int64, error) {
	var total int64
	err := o.db.Model(&entity.Order{}).Where("created_at >= ?", dateStart).Count(&total).Error
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (o *orderRepositoryImpl) GetOrderByID(id uint) (*entity.Order, error) {
	var order entity.Order
	err := o.db.Preload("Delivery").Preload("OrderDetails").First(&order, id).Error
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (o *orderRepositoryImpl) GetOrderDetailsByID(id uint) (*entity.OrderDetail, error) {
	var orderDetail entity.OrderDetail
	err := o.db.First(&orderDetail, id).Error
	if err != nil {
		return nil, err
	}

	return &orderDetail, nil
}

func (o *orderRepositoryImpl) CreateOrder(order entity.Order) (*entity.Order, error) {
	err := o.db.Create(&order).Error
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (o *orderRepositoryImpl) CreateOrderProcess(order entity.Order, details []entity.OrderDetail, delivery entity.Delivery) (*entity.Order, error) {
	tx := o.db.Begin()
	err := tx.Create(&order).Association("OrderDetails").Append(&details)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	delivery.OrderID = order.ID
	err = tx.Create(&delivery).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	o.db.Preload("Delivery").Preload("OrderDetails.Menu").Preload("OrderDetails").First(&order, order.ID)

	return &order, nil
}

func (o *orderRepositoryImpl) CreateBatchOrderDetails(orderDetails []entity.OrderDetail) (*[]entity.OrderDetail, error) {
	err := o.db.Create(&orderDetails).Error
	if err != nil {
		return nil, err
	}

	return &orderDetails, nil
}

func (o *orderRepositoryImpl) UserHasOrdered(userId, menuId uint) (bool, error) {
	var count int64
	err := o.db.Model(&entity.Order{}).Where("user_id = ? AND order_details.menu_id = ?", userId, menuId).Joins("JOIN order_details ON orders.id = order_details.order_id").Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (o *orderRepositoryImpl) UserHasReviewed(userId uint, orderDetailId uint) (bool, error) {
	var count int64
	err := o.db.Model(&entity.OrderDetail{}).Where("user_id = ? AND id = ?", userId, orderDetailId).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (o *orderRepositoryImpl) CreateCustomerReviewProcess(customerReview entity.CustomerReview) (*entity.CustomerReview, error) {
	tx := o.db.Begin()
	err := tx.Create(&customerReview).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	orderDetail := entity.OrderDetail{}
	err = tx.Model(&entity.OrderDetail{}).Where("id = ?", customerReview.OrderDetailID).First(&orderDetail).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	orderDetail.IsReviewed = true
	err = tx.Save(&orderDetail).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	menu := entity.Menu{}
	err = tx.Model(&entity.Menu{}).Where("id = ?", orderDetail.MenuID).First(&menu).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	menu.AvgRating = float64((menu.AvgRating*float64(menu.UserRatingCount) + float64(customerReview.Rating)) / float64(menu.UserRatingCount+1))
	menu.UserRatingCount = menu.UserRatingCount + 1
	err = tx.Save(&menu).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return &customerReview, nil
}

func (o *orderRepositoryImpl) GetCustomerReviewsByMenuId(menuId uint) (*[]entity.CustomerReview, error) {
	var customerReview []entity.CustomerReview
	err := o.db.Raw("SELECT * FROM customer_reviews JOIN users ON customer_reviews.user_id = users.id  JOIN order_details ON customer_reviews.order_detail_id = order_details.id WHERE order_details.menu_id = ?", menuId).Scan(&customerReview).Error
	if err != nil {
		return nil, err
	}

	return &customerReview, nil
}

func (o *orderRepositoryImpl) CreateCustomerReview(customerReview entity.CustomerReview) (*entity.CustomerReview, error) {
	err := o.db.Create(&customerReview).Error
	if err != nil {
		return nil, err
	}

	return &customerReview, nil
}
