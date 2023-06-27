package handler

import (
	"encoding/json"
	"errors"
	"final-project-backend/domain"
	"final-project-backend/dto"
	"final-project-backend/entity"
	"final-project-backend/util"
	"net/http"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func (h *Handler) CreateMenu(c *gin.Context) {
	var input dto.MenuFormRequest
	err := c.ShouldBindWith(&input, binding.FormMultipart)
	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInvalidBody.Error(), "INVALID_BODY_REQUEST", http.StatusBadRequest)
		return
	}

	var menuOptionsJson []entity.MenuOption
	err = json.Unmarshal([]byte(input.MenuOptions), &menuOptionsJson)

	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInvalidBody.Error(), "INVALID_BODY_REQUEST", http.StatusBadRequest)
		return
	}

	for _, menuOption := range menuOptionsJson {
		if menuOption.Title == "" || menuOption.Type == "" || menuOption.Max == 0 {
			util.ResponseErrorJSON(c, domain.ErrInvalidBody.Error(), "INVALID_BODY_REQUEST", http.StatusBadRequest)
			return
		}
		for _, menuOptionList := range menuOption.MenuOptionLists {
			if menuOptionList.Name == "" {
				util.ResponseErrorJSON(c, domain.ErrInvalidBody.Error(), "INVALID_BODY_REQUEST", http.StatusBadRequest)
				return
			}
		}
	}

	menu := entity.Menu{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		MenuOptions: input.MenuOptions,
	}

	var picUrl, publicId string
	if input.Picture.Size != 0 {
		picUrl, publicId, err = h.mediaUsecase.FileUpload(input.Picture)
		menu.PictureUrl = picUrl
		menu.PicturePublicId = publicId
	}

	menuRes, err := h.menuUsecase.CreateMenu(menu)

	if err == domain.ErrInvalidBody {
		util.ResponseErrorJSON(c, domain.ErrInvalidBody.Error(), "INVALID_BODY_REQUEST", http.StatusBadRequest)
		return
	}

	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInternalServer.Error(), "INSERT_MENU_FAILED", http.StatusInternalServerError)
		return
	}

	var menuCategory = []entity.CategoriesMenu{}
	input.Categories = util.UniqueUint(input.Categories)

	for _, category := range input.Categories {
		menuCategory = append(menuCategory, entity.CategoriesMenu{
			CategoryID: category,
			MenuID:     menuRes.ID,
		})
	}

	err = h.menuUsecase.CreateMenuCategories(menuCategory)
	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInvalidCategory.Error(), "INSERT_MENU_CATEGORY_FAILED", http.StatusBadRequest)
		return
	}

	dtoMenuRes := dto.MenuResponse{
		ID:          menuRes.ID,
		Name:        menuRes.Name,
		Description: menuRes.Description,
		Price:       menuRes.Price,
		PictureUrl:  menuRes.PictureUrl,
		MenuOptions: menuOptionsJson,
	}

	util.ResponseSuccesJSON(c, dtoMenuRes, http.StatusOK)
}

func sortedSquares(nums []int) []int {
	// lastInd := len(nums) - 1

	for i := 0; i < len(nums); i++ {
		nums[i] = nums[i] ^ 2
		sort.Ints(nums)
	}
	return nums
}

func (h *Handler) GetMenus(c *gin.Context) {
	query := dto.Query{}
	err := c.ShouldBindQuery(&query)
	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInvalidParams.Error(), "BAD_REQUEST", http.StatusBadRequest)
		return
	}

	menus, err := h.menuUsecase.GetMenus(query)
	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInternalServer.Error(), "GET_MENU_FAILED", http.StatusInternalServerError)
		return
	}

	var dtoMenus []dto.MenuResponse
	for _, menu := range menus {

		var menuOptions []entity.MenuOption
		json.Unmarshal([]byte(menu.MenuOptions), &menuOptions)

		dtoMenus = append(dtoMenus, dto.MenuResponse{
			ID:              menu.ID,
			Name:            menu.Name,
			Description:     menu.Description,
			Price:           menu.Price,
			PictureUrl:      menu.PictureUrl,
			AvgRating:       menu.AvgRating,
			UserRatingCount: menu.UserRatingCount,
			MenuOptions:     menuOptions,
			Categories:      menu.Categories,
		})
	}

	util.ResponseSuccesJSON(c, dtoMenus, http.StatusOK)
}

func (h *Handler) GetMenuById(c *gin.Context) {
	menuId := c.Param("id")

	menuIdInt, err := strconv.ParseInt(menuId, 10, 64)
	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInvalidBody.Error(), "INVALID_BODY_REQUEST", http.StatusBadRequest)
		return
	}

	menuIdUint := uint(menuIdInt)
	menu, err := h.menuUsecase.GetMenuById(menuIdUint)

	if errors.Is(err, domain.ErrMenuNotFound) {
		util.ResponseErrorJSON(c, domain.ErrMenuNotFound.Error(), "MENU_NOT_FOUND", http.StatusBadRequest)
		return
	}
	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInternalServer.Error(), "GET_MENU_FAILED", http.StatusInternalServerError)
		return
	}

	var menuOptions []entity.MenuOption
	json.Unmarshal([]byte(menu.MenuOptions), &menuOptions)

	dtoMenuRes := dto.MenuResponse{
		ID:              menu.ID,
		Name:            menu.Name,
		Description:     menu.Description,
		Price:           menu.Price,
		PictureUrl:      menu.PictureUrl,
		AvgRating:       menu.AvgRating,
		UserRatingCount: menu.UserRatingCount,
		MenuOptions:     menuOptions,
		Categories:      menu.Categories,
	}

	util.ResponseSuccesJSON(c, dtoMenuRes, http.StatusOK)
}

func (h *Handler) UpdateMenu(c *gin.Context) {
	menuId := c.Param("id")

	menuIdInt, err := strconv.ParseInt(menuId, 10, 64)
	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInvalidBody.Error(), "INVALID_BODY_REQUEST", http.StatusBadRequest)
		return
	}

	menuIdUint := uint(menuIdInt)
	menu, err := h.menuUsecase.GetMenuById(menuIdUint)

	if errors.Is(err, domain.ErrMenuNotFound) {
		util.ResponseErrorJSON(c, domain.ErrMenuNotFound.Error(), "MENU_NOT_FOUND", http.StatusBadRequest)
		return
	}
	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInternalServer.Error(), "GET_MENU_FAILED", http.StatusInternalServerError)
		return
	}

	var input dto.MenuFormRequest
	err = c.MustBindWith(&input, binding.FormMultipart)
	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInvalidBody.Error(), "INVALID_BODY_REQUEST", http.StatusBadRequest)
		return
	}

	var menuOptionsJson []entity.MenuOption
	err = json.Unmarshal([]byte(input.MenuOptions), &menuOptionsJson)

	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInvalidBody.Error(), "INVALID_BODY_REQUEST", http.StatusBadRequest)
		return
	}

	for _, menuOption := range menuOptionsJson {
		if menuOption.Title == "" || menuOption.Type == "" || menuOption.Max == 0 {
			util.ResponseErrorJSON(c, domain.ErrInvalidBody.Error(), "INVALID_BODY_REQUEST", http.StatusBadRequest)
			return
		}
		for _, menuOptionList := range menuOption.MenuOptionLists {
			if menuOptionList.Name == "" {
				util.ResponseErrorJSON(c, domain.ErrInvalidBody.Error(), "INVALID_BODY_REQUEST", http.StatusBadRequest)
				return
			}
		}
	}

	if (menu.PictureUrl != "") && (input.Picture.Size != 0) {
		err = h.mediaUsecase.FileDelete(menu.PicturePublicId)
	}

	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrDeleteImage.Error(), "DELETE_IMAGE_FAILED", 400)
		return
	}

	menu.Name = input.Name
	menu.Description = input.Description
	menu.Price = input.Price
	menu.MenuOptions = input.MenuOptions

	var picUrl, publicId string
	if input.Picture.Size != 0 {
		picUrl, publicId, err = h.mediaUsecase.FileUpload(input.Picture)
		menu.PictureUrl = picUrl
		menu.PicturePublicId = publicId
	}

	var menuCategory = []entity.CategoriesMenu{}
	input.Categories = util.UniqueUint(input.Categories)

	for _, category := range input.Categories {
		menuCategory = append(menuCategory, entity.CategoriesMenu{
			CategoryID: category,
			MenuID:     menuIdUint,
		})
	}

	err = h.menuUsecase.UpdateMenuCategories(menuIdUint, menuCategory)
	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInvalidCategory.Error(), "UPDATE_MENU_CATEGORY_FAILED", http.StatusBadRequest)
		return
	}
	menu.Categories = nil
	menu, err = h.menuUsecase.UpdateMenu(*menu)

	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrUpdateMenu.Error(), "UPDATE_MENU_FAILED", http.StatusInternalServerError)
		return
	}

	dtoMenuRes := dto.MenuResponse{
		ID:          menu.ID,
		Name:        menu.Name,
		Description: menu.Description,
		Price:       menu.Price,
		PictureUrl:  menu.PictureUrl,
		MenuOptions: menuOptionsJson,
	}

	util.ResponseSuccesJSON(c, dtoMenuRes, http.StatusOK)
}

func (h *Handler) DeleteMenu(c *gin.Context) {
	menuId := c.Param("id")

	menuIdInt, err := strconv.ParseInt(menuId, 10, 64)
	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInvalidBody.Error(), "INVALID_BODY_REQUEST", http.StatusBadRequest)
		return
	}

	menuIdUint := uint(menuIdInt)

	err = h.menuUsecase.DeleteMenu(menuIdUint)
	if errors.Is(err, domain.ErrMenuNotFound) {
		util.ResponseErrorJSON(c, domain.ErrMenuNotFound.Error(), "MENU_NOT_FOUND", http.StatusNotFound)
		return
	}

	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInternalServer.Error(), "DELETE_MENU_FAILED", http.StatusInternalServerError)
		return
	}

	util.ResponseSuccesJSON(c, nil, http.StatusNoContent)
}

func (h *Handler) ToggleFavoriteMenu(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		util.ResponseErrorJSON(c, domain.ErrUnauthorized.Error(), "UNAUTHORIZED", http.StatusUnauthorized)
		return
	}

	menuId := c.Param("id")
	uint64MenuId, err := strconv.ParseUint(menuId, 10, 64)
	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInvalidBody.Error(), "INVALID_BODY_REQUEST", http.StatusBadRequest)
		return
	}

	uintMenuId := uint(uint64MenuId)
	userId := user.(dto.UserResponse).ID

	favoriteMenu, err := h.menuUsecase.ToggleFavoriteMenu(userId, uintMenuId)
	if errors.Is(err, domain.ErrMenuNotFound) {
		util.ResponseErrorJSON(c, domain.ErrMenuNotFound.Error(), "MENU_NOT_FOUND", http.StatusNotFound)
		return
	}

	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInternalServer.Error(), "TOGGLE_FAVORITE_MENU_FAILED", http.StatusInternalServerError)
		return
	}

	if favoriteMenu == nil {
		util.ResponseSuccesJSON(c, nil, http.StatusNoContent)
		return
	}

	util.ResponseSuccesJSON(c, favoriteMenu, http.StatusOK)
}

func (h *Handler) GetFavoriteMenus(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		util.ResponseErrorJSON(c, domain.ErrUnauthorized.Error(), "UNAUTHORIZED", http.StatusUnauthorized)
		return
	}

	userId := user.(dto.UserResponse).ID

	favoriteMenus, err := h.menuUsecase.GetUserFavoriteMenus(userId)
	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInternalServer.Error(), "GET_FAVORITE_MENU_FAILED", http.StatusInternalServerError)
		return
	}

	util.ResponseSuccesJSON(c, favoriteMenus, http.StatusOK)
}

func (h *Handler) GetCategories(c *gin.Context) {
	categories, err := h.menuUsecase.GetCategories()
	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInternalServer.Error(), "GET_MENU_FAILED", http.StatusInternalServerError)
		return
	}

	util.ResponseSuccesJSON(c, categories, http.StatusOK)
}
