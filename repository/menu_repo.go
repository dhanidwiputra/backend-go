package repository

import (
	"final-project-backend/domain"
	"final-project-backend/dto"
	"final-project-backend/entity"

	"gorm.io/gorm"
)

type MenuRepository interface {
	CreateMenu(entity.Menu) (*entity.Menu, error)
	GetMenus(dto.Query) ([]entity.Menu, error)
	CreateMenuCategories([]entity.CategoriesMenu) error
	UpdateMenu(entity.Menu) (*entity.Menu, error)
	DeleteMenu(entity.Menu) error
	GetMenuById(uint) (*entity.Menu, error)
	DeleteMenuCategoriesByMenuId(uint) error
	IsValidMenuOptions(uint, string) error
	ToggleFavoriteMenu(uint, uint) (*entity.MenuFavorite, error)
	CreateUserFavoriteMenu(entity.MenuFavorite) (*entity.MenuFavorite, error)
	DeleteUserFavoriteMenu(entity.MenuFavorite) error
	GetUserFavoriteMenu(uint, uint) (*entity.MenuFavorite, error)
	GetUserFavoriteMenus(uint) ([]entity.Menu, error)
	GetCategories() ([]entity.Category, error)
}

type menuRepositoryImpl struct {
	db *gorm.DB
}

type MenuRepoConfig struct {
	DB *gorm.DB
}

func NewMenuRepository(c MenuRepoConfig) MenuRepository {
	return &menuRepositoryImpl{db: c.DB}
}

func (r *menuRepositoryImpl) CreateMenu(menu entity.Menu) (*entity.Menu, error) {
	err := r.db.Create(&menu).Error

	if err != nil {
		return nil, err
	}

	r.db.Preload("MenuOptions.MenuOptionLists").Preload("Categories").First(&menu)

	return &menu, nil
}

func (r *menuRepositoryImpl) GetMenus(query dto.Query) ([]entity.Menu, error) {
	var menus []entity.Menu
	err := r.db.Raw("SELECT distinct m.* FROM menus m JOIN categories_menus cm ON cm.menu_id = m.id JOIN categories c ON c.id = cm.category_id WHERE c.name ILIKE ? and m.name ILIKE ?  and m.deleted_at IS NULL", "%"+query.Category+"%", "%"+query.Search+"%").Scan(&menus).Error

	if err != nil {
		return nil, err
	}

	// PRELOAD CATEGORIES from generated entry on menus
	for i := range menus {
		r.db.Preload("Categories").First(&menus[i])
	}

	return menus, nil
}

func (r *menuRepositoryImpl) CreateMenuCategories(menuCategories []entity.CategoriesMenu) error {
	err := r.db.Create(&menuCategories).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *menuRepositoryImpl) DeleteMenuCategoriesByMenuId(menuId uint) error {
	res := r.db.Where("menu_id = ?", menuId).Unscoped().Delete(&entity.CategoriesMenu{})

	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *menuRepositoryImpl) UpdateMenu(menu entity.Menu) (*entity.Menu, error) {
	err := r.db.Save(&menu).Error

	if err != nil {
		return nil, err
	}

	err = r.db.Preload("Categories").Find(&menu).Error

	if err != nil {
		return nil, err
	}

	return &menu, nil

}

func (r *menuRepositoryImpl) DeleteMenu(menu entity.Menu) error {
	err := r.db.Delete(&menu).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *menuRepositoryImpl) GetMenuById(id uint) (*entity.Menu, error) {
	var menu entity.Menu
	err := r.db.Preload("Categories").Where("id = ?", id).First(&menu).Error

	if err != nil {
		return nil, err
	}

	return &menu, nil
}

func (r *menuRepositoryImpl) IsValidMenuOptions(menuId uint, option string) error {
	var menu *entity.Menu
	res := r.db.Where("id = ? and menu_options @> ?", menuId, option).First(&menu)
	if res.RowsAffected == 0 {
		return domain.ErrMalformedRequest
	}

	return nil
}

func (r *menuRepositoryImpl) ToggleFavoriteMenu(userId uint, menuId uint) (*entity.MenuFavorite, error) {
	favoriteMenu, err := r.GetUserFavoriteMenu(userId, menuId)
	if favoriteMenu != nil {
		err = r.DeleteUserFavoriteMenu(*favoriteMenu)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}

	favoriteMenu = &entity.MenuFavorite{
		MenuID: menuId,
		UserID: userId,
	}
	favoriteMenuRes, err := r.CreateUserFavoriteMenu(*favoriteMenu)
	if err != nil {
		return nil, err
	}
	return favoriteMenuRes, nil
}

func (r *menuRepositoryImpl) CreateUserFavoriteMenu(favoriteMenu entity.MenuFavorite) (*entity.MenuFavorite, error) {
	err := r.db.Create(&favoriteMenu).Error

	if err != nil {
		return nil, err
	}

	return &favoriteMenu, nil
}

func (r *menuRepositoryImpl) DeleteUserFavoriteMenu(favoriteMenu entity.MenuFavorite) error {
	err := r.db.Delete(&favoriteMenu).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *menuRepositoryImpl) GetUserFavoriteMenu(userId uint, menuId uint) (*entity.MenuFavorite, error) {
	var favoriteMenu entity.MenuFavorite
	err := r.db.Where("menu_id = ? and user_id = ?", menuId, userId).First(&favoriteMenu).Error

	if err != nil {
		return nil, err
	}

	return &favoriteMenu, nil
}

func (r *menuRepositoryImpl) GetUserFavoriteMenus(userId uint) ([]entity.Menu, error) {
	var favoriteMenus []entity.MenuFavorite
	err := r.db.Where("user_id = ?", userId).Find(&favoriteMenus).Error

	if err != nil {
		return nil, err
	}

	var menuIds []uint
	for _, menu := range favoriteMenus {
		menuIds = append(menuIds, menu.MenuID)
	}

	var menus []entity.Menu
	err = r.db.Where("id in (?)", menuIds).Find(&menus).Error

	if err != nil {
		return nil, err
	}

	return menus, nil
}

func (r *menuRepositoryImpl) GetCategories() ([]entity.Category, error) {
	var categories []entity.Category
	err := r.db.Find(&categories).Error

	if err != nil {
		return nil, err
	}

	return categories, nil
}
