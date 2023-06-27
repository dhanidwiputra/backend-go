package usecase

import (
	"final-project-backend/domain"
	"final-project-backend/dto"
	"final-project-backend/entity"
	"final-project-backend/repository"
)

type CartUsecase interface {
	GetCartItems(userId uint) ([]entity.UserCartItem, error)
	GetCartItemById(id uint) (*entity.UserCartItem, error)
	GetCartItemsByDto(dto.CartItemData) (*entity.UserCartItem, error)
	AddToCart(dto.CartItemData) (*entity.UserCartItem, error)
	UpdateCartItem(entity.UserCartItem) (*entity.UserCartItem, error)
	DeleteCartItem(id uint) error
	EmptyCart(userId int) error
}

type cartUsecaseImpl struct {
	cartRepo repository.CartRepository
}

type CartUsecaseConfig struct {
	CartRepo repository.CartRepository
}

func NewCartUsecase(c CartUsecaseConfig) CartUsecase {
	return &cartUsecaseImpl{
		cartRepo: c.CartRepo,
	}
}

func (u *cartUsecaseImpl) GetCartItems(userId uint) ([]entity.UserCartItem, error) {
	cart, err := u.cartRepo.GetCartItems(userId)

	if err != nil {
		return nil, err
	}

	return cart, nil
}

func (u *cartUsecaseImpl) GetCartItemById(id uint) (*entity.UserCartItem, error) {
	cart, err := u.cartRepo.GetCartItemById(id)
	if cart == nil {
		return nil, domain.ErrCartItemNotFound
	}

	if err != nil {
		return nil, err
	}

	return cart, nil
}

func (u *cartUsecaseImpl) GetCartItemsByDto(item dto.CartItemData) (*entity.UserCartItem, error) {
	var cart *entity.UserCartItem
	cart, err := u.cartRepo.GetCartItemsByDto(item)

	if cart == nil {
		return nil, domain.ErrCartItemNotFound
	}
	if err != nil {
		return nil, err
	}

	return cart, nil
}

func (u *cartUsecaseImpl) AddToCart(item dto.CartItemData) (*entity.UserCartItem, error) {
	cartRes, err := u.cartRepo.AddToCart(item)

	if err != nil {
		return nil, err
	}

	return cartRes, nil
}

func (u *cartUsecaseImpl) UpdateCartItem(item entity.UserCartItem) (*entity.UserCartItem, error) {
	cartRes, err := u.cartRepo.UpdateCartItem(item)

	if err != nil {
		return nil, err
	}

	return cartRes, nil
}

func (u *cartUsecaseImpl) DeleteCartItem(id uint) error {
	cart, err := u.cartRepo.GetCartItemById(id)

	if err != nil {
		return domain.ErrCartItemNotFound
	}

	err = u.cartRepo.DeleteCartItem(*cart)

	if err != nil {
		return err
	}

	return nil
}

func (u *cartUsecaseImpl) EmptyCart(userId int) error {
	err := u.cartRepo.EmptyCart(userId)

	if err != nil {
		return err
	}

	return nil
}
