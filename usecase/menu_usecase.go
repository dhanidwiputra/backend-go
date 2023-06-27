package usecase

import (
	"final-project-backend/domain"
	"final-project-backend/dto"
	"final-project-backend/entity"
	"final-project-backend/repository"
	"final-project-backend/util"
)

type MenuUsecase interface {
	CreateMenu(menu entity.Menu) (*entity.Menu, error)
	GetMenus(dto.Query) ([]entity.Menu, error)
	CreateMenuCategories([]entity.CategoriesMenu) error
	UpdateMenu(menu entity.Menu) (*entity.Menu, error)
	DeleteMenu(id uint) error
	UpdateMenuCategories(id uint, menuCategory []entity.CategoriesMenu) error
	GetMenuById(id uint) (*entity.Menu, error)
	IsValidMenuOptions(id uint, menuOptions string) error
	ToggleFavoriteMenu(userId uint, menuId uint) (*entity.MenuFavorite, error)
	GetUserFavoriteMenus(userId uint) ([]entity.Menu, error)
	GetCategories() ([]entity.Category, error)
}

type menuUsecaseImpl struct {
	menuUsecase   MenuUsecase
	menuRepo      repository.MenuRepository
	mediaUploader util.MediaUploader
	orderRepo     repository.OrderRepository
}

type MenuUsecaseConfig struct {
	MenuUsecase   MenuUsecase
	MenuRepo      repository.MenuRepository
	MediaUploader util.MediaUploader
	OrderRepo     repository.OrderRepository
}

func NewMenuUsecase(c MenuUsecaseConfig) MenuUsecase {
	return &menuUsecaseImpl{
		menuUsecase:   c.MenuUsecase,
		menuRepo:      c.MenuRepo,
		mediaUploader: c.MediaUploader,
		orderRepo:     c.OrderRepo,
	}
}

func (m *menuUsecaseImpl) CreateMenu(menu entity.Menu) (*entity.Menu, error) {
	resMenu, err := m.menuRepo.CreateMenu(menu)
	if err != nil {
		return nil, err
	}

	return resMenu, nil
}

func (m *menuUsecaseImpl) CreateMenuCategories(menuCategory []entity.CategoriesMenu) error {
	err := m.menuRepo.CreateMenuCategories(menuCategory)
	if err != nil {
		return err
	}

	return nil
}

func (m *menuUsecaseImpl) GetMenus(query dto.Query) ([]entity.Menu, error) {
	menus, err := m.menuRepo.GetMenus(query)
	if err != nil {
		return nil, err
	}

	return menus, nil
}

func (m *menuUsecaseImpl) GetMenuById(id uint) (*entity.Menu, error) {
	menu, err := m.menuRepo.GetMenuById(id)
	if menu == nil {
		return nil, domain.ErrMenuNotFound
	}

	if err != nil {
		return nil, domain.ErrInternalServer
	}

	return menu, nil
}

func (m *menuUsecaseImpl) UpdateMenu(menu entity.Menu) (*entity.Menu, error) {
	resMenu, err := m.menuRepo.UpdateMenu(menu)
	if err != nil {
		return nil, domain.ErrUpdateMenu
	}

	return resMenu, nil
}

func (m *menuUsecaseImpl) UpdateMenuCategories(id uint, menuCategory []entity.CategoriesMenu) error {
	err := m.menuRepo.DeleteMenuCategoriesByMenuId(id)
	if err != nil {
		return err
	}

	err = m.menuRepo.CreateMenuCategories(menuCategory)
	if err != nil {
		return err
	}

	return nil
}

func (m *menuUsecaseImpl) DeleteMenu(id uint) error {
	menu, err := m.GetMenuById(id)
	if err != nil {
		return domain.ErrMenuNotFound
	}

	err = m.menuRepo.DeleteMenu(*menu)
	if err != nil {
		return err
	}

	return nil
}

func (m *menuUsecaseImpl) IsValidMenuOptions(menuId uint, option string) error {
	err := m.menuRepo.IsValidMenuOptions(menuId, option)
	if err != nil {
		return domain.ErrMalformedRequest
	}

	return nil
}

func (m *menuUsecaseImpl) ToggleFavoriteMenu(userId uint, menuId uint) (*entity.MenuFavorite, error) {
	menu, _ := m.menuRepo.GetMenuById(menuId)
	if menu == nil {
		return nil, domain.ErrMenuNotFound
	}

	favoriteMenu, err := m.menuRepo.ToggleFavoriteMenu(userId, menuId)
	if err != nil {
		return nil, err
	}

	return favoriteMenu, nil
}

func (m *menuUsecaseImpl) GetUserFavoriteMenus(userId uint) ([]entity.Menu, error) {
	menus, err := m.menuRepo.GetUserFavoriteMenus(userId)
	if err != nil {
		return nil, err
	}

	return menus, nil
}

func (m *menuUsecaseImpl) GetCategories() ([]entity.Category, error) {
	categories, err := m.menuRepo.GetCategories()
	if err != nil {
		return nil, err
	}

	return categories, nil
}
