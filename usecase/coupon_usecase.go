package usecase

import (
	"final-project-backend/dto"
	"final-project-backend/entity"
	"final-project-backend/repository"

	"github.com/gofrs/uuid"
)

type CouponUsecase interface {
	CreateCoupon(dto.CouponRequest) (*entity.Coupon, error)
	AssignCouponToUser(entity.Coupon, dto.UserResponse) (*entity.UsersCoupon, error)
	GetCoupons() ([]entity.Coupon, error)
	GetCouponsByUserId(id uint) ([]entity.UsersCoupon, error)
	GetCouponByDiscount(discount uint) (*entity.Coupon, error)
	GetCouponById(id uint) (*entity.Coupon, error)
	GetUserCouponByFK(uint, uint) (*entity.UsersCoupon, error)
	UpdateCoupon(entity.Coupon) (*entity.Coupon, error)
	DeleteCoupon(entity.Coupon) error
	UpdateUserCoupon(entity.UsersCoupon) (*entity.UsersCoupon, error)
}

type couponUsecaseImpl struct {
	userRepo   repository.UserRepository
	couponRepo repository.CouponRepository
}

type CouponUsecaseConfig struct {
	UserRepo   repository.UserRepository
	CouponRepo repository.CouponRepository
}

func NewCouponUsecase(c CouponUsecaseConfig) CouponUsecase {
	return &couponUsecaseImpl{
		userRepo:   c.UserRepo,
		couponRepo: c.CouponRepo,
	}
}

func (c *couponUsecaseImpl) CreateCoupon(b dto.CouponRequest) (*entity.Coupon, error) {
	coupon := entity.Coupon{
		Code:         uuid.Must(uuid.NewV4()),
		IssuerID:     b.IssuerID,
		Discount:     b.Discount,
		Description:  b.Description,
		Availability: true,
	}

	resCoupon, err := c.couponRepo.CreateCoupon(coupon)

	if err != nil {
		return nil, err
	}

	return resCoupon, nil
}

func (c *couponUsecaseImpl) AssignCouponToUser(coupon entity.Coupon, user dto.UserResponse) (*entity.UsersCoupon, error) {
	userCoupon, err := c.couponRepo.GetUserCouponByFK(coupon.ID, user.ID)
	if userCoupon != nil {
		userCoupon.Stock += 1
		return c.couponRepo.AssignCouponToUser(*userCoupon)
	}

	userCoupon = &entity.UsersCoupon{
		CouponID: coupon.ID,
		UserID:   user.ID,
		Stock:    1,
	}

	userRes, err := c.couponRepo.AssignCouponToUser(*userCoupon)

	if err != nil {
		return nil, err
	}

	return userRes, nil
}

func (c *couponUsecaseImpl) GetCoupons() ([]entity.Coupon, error) {
	coupons, err := c.couponRepo.GetCoupons()

	if err != nil {
		return nil, err
	}

	return coupons, nil
}

func (c *couponUsecaseImpl) GetCouponById(id uint) (*entity.Coupon, error) {
	coupon, err := c.couponRepo.GetCouponById(id)

	if err != nil {
		return nil, err
	}

	return coupon, nil
}

func (c *couponUsecaseImpl) GetCouponByDiscount(discount uint) (*entity.Coupon, error) {
	coupon, err := c.couponRepo.GetCouponByDiscount(discount)

	if err != nil {
		return nil, err
	}

	return coupon, nil
}

func (c *couponUsecaseImpl) GetCouponsByUserId(id uint) ([]entity.UsersCoupon, error) {
	coupons, err := c.couponRepo.GetCouponsByUserId(id)

	if err != nil {
		return nil, err
	}

	return coupons, nil
}

func (c *couponUsecaseImpl) GetUserCouponByFK(couponId uint, userId uint) (*entity.UsersCoupon, error) {
	userCoupon, err := c.couponRepo.GetUserCouponByFK(couponId, userId)

	if err != nil {
		return nil, err
	}

	return userCoupon, nil
}

func (c *couponUsecaseImpl) DeleteUserCoupon(userCoupon entity.UsersCoupon) error {
	err := c.couponRepo.DeleteUserCoupon(userCoupon)

	if err != nil {
		return err
	}

	return nil
}

func (c *couponUsecaseImpl) UpdateCoupon(coupon entity.Coupon) (*entity.Coupon, error) {
	resCoupon, err := c.couponRepo.UpdateCoupon(coupon)

	if err != nil {
		return nil, err
	}

	return resCoupon, nil
}

func (c *couponUsecaseImpl) DeleteCoupon(coupon entity.Coupon) error {
	err := c.couponRepo.DeleteCoupon(coupon)

	if err != nil {
		return err
	}

	return nil
}

func (c *couponUsecaseImpl) UpdateUserCoupon(userCoupon entity.UsersCoupon) (*entity.UsersCoupon, error) {
	if userCoupon.Stock == 0 {
		err := c.couponRepo.DeleteUserCoupon(userCoupon)

		if err != nil {
			return nil, err
		}
	}

	resUserCoupon, err := c.couponRepo.UpdateUserCoupon(userCoupon)

	if err != nil {
		return nil, err
	}

	return resUserCoupon, nil
}
