package repository

import (
	"final-project-backend/dto"
	"final-project-backend/entity"

	"gorm.io/gorm"
)

type CartRepository interface {
	AddToCart(dto.CartItemData) (*entity.UserCartItem, error)
	UpdateCartItem(entity.UserCartItem) (*entity.UserCartItem, error)
	DeleteCartItem(entity.UserCartItem) error
	GetCartItems(userId uint) ([]entity.UserCartItem, error)
	GetCartItemById(id uint) (*entity.UserCartItem, error)
	GetCartItemsByDto(dto.CartItemData) (*entity.UserCartItem, error)
	EmptyCart(userId int) error
}

type cartRepositoryImpl struct {
	db *gorm.DB
}

type CartRepoConfig struct {
	DB *gorm.DB
}

func NewCartRepository(c CartRepoConfig) CartRepository {
	return &cartRepositoryImpl{db: c.DB}
}

func (r *cartRepositoryImpl) AddToCart(item dto.CartItemData) (*entity.UserCartItem, error) {
	res, err := r.GetCartItemsByDto(item)
	if res != nil {
		res.Quantity += item.Quantity
		cart, err := r.UpdateCartItem(*res)
		return cart, err
	}

	cart := entity.UserCartItem{
		UserID:      item.UserID,
		MenuID:      item.MenuID,
		Quantity:    item.Quantity,
		MenuOptions: item.MenuOptions,
	}

	err = r.db.Create(&cart).Error

	if err != nil {
		return nil, err
	}

	err = r.db.Preload("Menu").Where("id = ?", cart.ID).First(&cart).Error
	if err != nil {
		return nil, err
	}

	return &cart, nil
}

func (r *cartRepositoryImpl) GetCartItems(userId uint) ([]entity.UserCartItem, error) {
	var carts []entity.UserCartItem
	err := r.db.Preload("Menu").Where("user_id = ?", userId).Order("created_at desc").Find(&carts).Error

	if err != nil {
		return nil, err
	}

	return carts, nil
}

func (r *cartRepositoryImpl) GetCartItemById(id uint) (*entity.UserCartItem, error) {
	var cart *entity.UserCartItem
	err := r.db.Preload("Menu").Where("id = ?", id).First(&cart).Error

	if err != nil {
		return nil, err
	}

	return cart, nil
}

func (r *cartRepositoryImpl) GetCartItemsByDto(item dto.CartItemData) (*entity.UserCartItem, error) {
	var cart *entity.UserCartItem
	err := r.db.Where("user_id = ? AND menu_id = ? and menu_options @> ? and menu_options <@ ?", item.UserID, item.MenuID, item.MenuOptions, item.MenuOptions).First(&cart).Error

	if err != nil {
		return nil, err
	}

	r.db.Preload("Menu").Where("id = ?", cart.ID).First(&cart)

	return cart, nil
}

func (r *cartRepositoryImpl) CreateCartItem(item entity.UserCartItem) (*entity.UserCartItem, error) {
	err := r.db.Create(&item).Error

	if err != nil {
		return nil, err
	}

	r.db.Preload("Menu").Where("id = ?", item.ID).First(&item)

	return &item, nil
}

func (r *cartRepositoryImpl) UpdateCartItem(item entity.UserCartItem) (*entity.UserCartItem, error) {
	if item.Quantity <= 0 {
		err := r.DeleteCartItem(item)
		return nil, err
	}

	err := r.db.Save(&item).Error

	if err != nil {
		return nil, err
	}

	r.db.Preload("Menu").Where("id = ?", item.ID).First(&item)

	return &item, nil
}

func (r *cartRepositoryImpl) DeleteCartItem(item entity.UserCartItem) error {
	err := r.db.Delete(&item).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *cartRepositoryImpl) EmptyCart(userId int) error {
	err := r.db.Where("user_id = ?", userId).Delete(&entity.UserCartItem{}).Error
	if err != nil {
		return err
	}

	return nil
}
