package repository

import (
	"final-project-backend/entity"

	"gorm.io/gorm"
)

type PromotionRepository interface {
	GetPromotions() ([]entity.Promotion, error)
	GetPromotionById(uint) (*entity.Promotion, error)
	CreatePromotion(entity.Promotion) (*entity.Promotion, error)
	UpdatePromotion(entity.Promotion) (*entity.Promotion, error)
	DeletePromotion(entity.Promotion) error
	DeletePromotionDetails(promotionID uint) error
	DeletePaymentRequirements(promotionID uint) error
}

type promotionRepositoryImpl struct {
	db *gorm.DB
}

type PromotionRepoConfig struct {
	DB *gorm.DB
}

func NewPromotionRepository(c PromotionRepoConfig) PromotionRepository {
	return &promotionRepositoryImpl{db: c.DB}
}

func (r *promotionRepositoryImpl) GetPromotions() ([]entity.Promotion, error) {
	var promotions []entity.Promotion
	err := r.db.Preload("PromotionDetails.Menu").Preload("PromotionDetails").Preload("PaymentRequirements.PaymentOption").Preload("PaymentRequirements").Find(&promotions).Error

	if err != nil {
		return nil, err
	}

	return promotions, nil
}

func (r *promotionRepositoryImpl) GetPromotionById(id uint) (*entity.Promotion, error) {
	var promotion entity.Promotion
	err := r.db.Preload("PromotionDetails.Menu").Preload("PromotionDetails").Preload("PaymentRequirements.PaymentOption").Preload("PaymentRequirements").First(&promotion, id).Error

	if err != nil {
		return nil, err
	}

	return &promotion, nil
}

func (r *promotionRepositoryImpl) CreatePromotion(promotion entity.Promotion) (*entity.Promotion, error) {
	err := r.db.Create(&promotion).Error
	if err != nil {
		return nil, err
	}

	r.db.Preload("PromotionDetails.Menu").Preload("PromotionDetails").Preload("PaymentRequirements.PaymentOption").Preload("PaymentRequirements").First(&promotion)

	return &promotion, nil
}

func (r *promotionRepositoryImpl) UpdatePromotion(promotion entity.Promotion) (*entity.Promotion, error) {
	err := r.db.Save(&promotion).Error
	if err != nil {
		return nil, err
	}

	r.db.Preload("PromotionDetails.Menu").Preload("PromotionDetails").Preload("PaymentRequirements.PaymentOption").Preload("PaymentRequirements").First(&promotion)

	return &promotion, nil
}

func (r *promotionRepositoryImpl) DeletePromotion(promotion entity.Promotion) error {
	err := r.db.Delete(&promotion).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *promotionRepositoryImpl) CreateBatchPromotionDetails(promotionDetails []entity.PromotionDetail) error {
	err := r.db.Create(&promotionDetails).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *promotionRepositoryImpl) CreateBatchPaymentRequirements(paymentRequirements []entity.PaymentRequirement) error {
	err := r.db.Create(&paymentRequirements).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *promotionRepositoryImpl) DeletePromotionDetails(promotionID uint) error {
	err := r.db.Where("promotion_id = ?", promotionID).Delete(&entity.PromotionDetail{}).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *promotionRepositoryImpl) DeletePaymentRequirements(promotionID uint) error {
	err := r.db.Where("promotion_id = ?", promotionID).Delete(&entity.PaymentRequirement{}).Error

	if err != nil {
		return err
	}

	return nil
}
