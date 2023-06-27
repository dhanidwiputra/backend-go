package repository

import (
	"final-project-backend/domain"
	"final-project-backend/entity"

	"gorm.io/gorm"
)

type CouponRepository interface {
	CreateCoupon(entity.Coupon) (*entity.Coupon, error)
	GetCouponById(id uint) (*entity.Coupon, error)
	GetUserCouponByFK(uint, uint) (*entity.UsersCoupon, error)
	DeleteUserCoupon(entity.UsersCoupon) error
	GetCouponsByUserId(id uint) ([]entity.UsersCoupon, error)
	GetCoupons() ([]entity.Coupon, error)
	GetCouponByDiscount(discount uint) (*entity.Coupon, error)
	UpdateCoupon(entity.Coupon) (*entity.Coupon, error)
	DeleteCoupon(entity.Coupon) error
	AssignCouponToUser(entity.UsersCoupon) (*entity.UsersCoupon, error)
	UpdateUserCoupon(entity.UsersCoupon) (*entity.UsersCoupon, error)
}

type couponRepositoryImpl struct {
	db *gorm.DB
}

type CouponRepoConfig struct {
	DB *gorm.DB
}

func NewCouponRepository(c CouponRepoConfig) CouponRepository {
	return &couponRepositoryImpl{db: c.DB}
}

func (r *couponRepositoryImpl) GetCouponById(id uint) (*entity.Coupon, error) {
	var coupon entity.Coupon

	res := r.db.Find(&coupon, id)

	if res.RowsAffected == 0 {
		return nil, domain.ErrCouponNotFound
	}

	return &coupon, nil
}

func (r *couponRepositoryImpl) CreateCoupon(coupon entity.Coupon) (*entity.Coupon, error) {
	err := r.db.Create(&coupon).Error

	if err != nil {
		return nil, err
	}

	return &coupon, nil
}

func (r *couponRepositoryImpl) AssignCouponToUser(userCoupon entity.UsersCoupon) (*entity.UsersCoupon, error) {
	err := r.db.Create(&userCoupon).Error

	if err != nil {
		return nil, err
	}

	return &userCoupon, nil
}

func (r *couponRepositoryImpl) UpdateUserCoupon(userCoupon entity.UsersCoupon) (*entity.UsersCoupon, error) {
	err := r.db.Save(&userCoupon).Error

	if err != nil {
		return nil, err
	}

	return &userCoupon, nil
}

func (r *couponRepositoryImpl) GetCoupons() ([]entity.Coupon, error) {
	var coupons []entity.Coupon

	err := r.db.Find(&coupons).Error

	if err != nil {
		return nil, err
	}

	return coupons, nil
}

func (r *couponRepositoryImpl) GetCouponsByUserId(id uint) ([]entity.UsersCoupon, error) {
	var coupons []entity.UsersCoupon

	err := r.db.Preload("Coupon").Where("user_id = ?", id).Find(&coupons).Error

	if err != nil {
		return nil, err
	}

	return coupons, nil
}

func (r *couponRepositoryImpl) GetCouponByDiscount(discount uint) (*entity.Coupon, error) {
	var coupon entity.Coupon

	err := r.db.Where("discount = ?", discount).Find(&coupon).Error

	if err != nil {
		return nil, err
	}

	return &coupon, nil
}

func (r *couponRepositoryImpl) GetUserCouponByFK(couponId uint, userId uint) (*entity.UsersCoupon, error) {
	var userCoupon entity.UsersCoupon

	err := r.db.Where("user_id = ? AND coupon_id = ?", userId, couponId).First(&userCoupon).Error

	if err != nil {
		return nil, err
	}

	return &userCoupon, nil
}

func (r *couponRepositoryImpl) DeleteUserCoupon(userCoupon entity.UsersCoupon) error {
	err := r.db.Delete(&userCoupon).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *couponRepositoryImpl) UpdateCoupon(coupon entity.Coupon) (*entity.Coupon, error) {
	err := r.db.Save(&coupon).Error

	if err != nil {
		return nil, err
	}

	return &coupon, nil
}

func (r *couponRepositoryImpl) DeleteCoupon(coupon entity.Coupon) error {
	err := r.db.Delete(&coupon).Error

	if err != nil {
		return err
	}

	return nil
}
