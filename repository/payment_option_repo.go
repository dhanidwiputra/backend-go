package repository

import (
	"final-project-backend/entity"

	"gorm.io/gorm"
)

type PaymentOptionRepository interface {
	GetAllPaymentOptions() ([]entity.PaymentOption, error)
	GetPaymentOptionById(id uint) (*entity.PaymentOption, error)
}

type paymentOptionRepositoryImpl struct {
	db *gorm.DB
}

type PaymentOptionRepositoryConfig struct {
	DB *gorm.DB
}

func NewPaymentOptionRepository(c PaymentOptionRepositoryConfig) PaymentOptionRepository {
	return &paymentOptionRepositoryImpl{
		db: c.DB,
	}
}

func (p *paymentOptionRepositoryImpl) GetAllPaymentOptions() ([]entity.PaymentOption, error) {
	var paymentOptions []entity.PaymentOption
	err := p.db.Find(&paymentOptions).Error
	if err != nil {
		return nil, err
	}

	return paymentOptions, nil
}

func (p *paymentOptionRepositoryImpl) GetPaymentOptionById(id uint) (*entity.PaymentOption, error) {
	var paymentOption entity.PaymentOption
	err := p.db.First(&paymentOption, id).Error
	if err != nil {
		return nil, err
	}

	return &paymentOption, nil
}
