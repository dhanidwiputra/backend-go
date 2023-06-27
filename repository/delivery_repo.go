package repository

import (
	"final-project-backend/entity"

	"gorm.io/gorm"
)

type DeliveryRepository interface {
	GetDeliveryByID(id uint) (*entity.Delivery, error)
	CreateDelivery(entity.Delivery) (*entity.Delivery, error)
	UpdateDeliveryStatus(entity.Delivery) (*entity.Delivery, error)
}

type deliveryRepositoryImpl struct {
	db *gorm.DB
}

type DeliveryRepoConfig struct {
	DB *gorm.DB
}

func NewDeliveryRepository(c DeliveryRepoConfig) DeliveryRepository {
	return &deliveryRepositoryImpl{db: c.DB}
}

func (r *deliveryRepositoryImpl) CreateDelivery(delivery entity.Delivery) (*entity.Delivery, error) {
	err := r.db.Create(&delivery).Error

	if err != nil {
		return nil, err
	}
	return &delivery, nil
}

func (r *deliveryRepositoryImpl) GetDeliveryByID(id uint) (*entity.Delivery, error) {
	var delivery entity.Delivery
	err := r.db.First(&delivery, id).Error

	if err != nil {
		return nil, err
	}

	return &delivery, nil
}

func (r *deliveryRepositoryImpl) UpdateDeliveryStatus(delivery entity.Delivery) (*entity.Delivery, error) {
	err := r.db.Model(&delivery).Update("status", delivery.Status).Error

	if err != nil {
		return nil, err
	}

	return &delivery, nil
}
