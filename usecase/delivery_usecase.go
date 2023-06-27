package usecase

import (
	"final-project-backend/domain"
	"final-project-backend/entity"
	"final-project-backend/repository"
	"time"
)

var DeliveryStatus = map[string]string{
	"pending":    "Pending",
	"on the way": "On the way",
	"delivered":  "Delivered",
}

type DeliveryUsecase interface {
	CreateDelivery(entity.Delivery) (*entity.Delivery, error)
	UpdateDeliveryStatus(entity.Delivery) (*entity.Delivery, error)
	GetDeliveryById(id uint) (*entity.Delivery, error)
}

type deliveryUsecaseImpl struct {
	deliveryRepo repository.DeliveryRepository
}

type DeliveryUsecaseConfig struct {
	DeliveryRepo repository.DeliveryRepository
}

func NewDeliveryUsecase(c DeliveryUsecaseConfig) DeliveryUsecase {
	return &deliveryUsecaseImpl{
		deliveryRepo: c.DeliveryRepo,
	}
}

func (d *deliveryUsecaseImpl) CreateDelivery(delivery entity.Delivery) (*entity.Delivery, error) {
	deliveryRes, err := d.deliveryRepo.CreateDelivery(delivery)
	if err != nil {
		return nil, err
	}

	return deliveryRes, nil
}

func (d *deliveryUsecaseImpl) GetDeliveryById(id uint) (*entity.Delivery, error) {
	deliveryRes, err := d.deliveryRepo.GetDeliveryByID(id)
	if err != nil {
		return nil, err
	}

	return deliveryRes, nil
}

func (d *deliveryUsecaseImpl) UpdateDeliveryStatus(delivery entity.Delivery) (*entity.Delivery, error) {
	if _, ok := DeliveryStatus[delivery.Status]; !ok {
		return nil, domain.ErrInvalidDeliveryStatus
	}

	if delivery.Status == "on the way" {
		delivery.DeliveryDate = time.Now()
	}

	deliveryRes, err := d.deliveryRepo.UpdateDeliveryStatus(delivery)
	if err != nil {
		return nil, err
	}

	return deliveryRes, nil
}
