package usecase_test

// func TestToggleFavoriteMenu(t *testing.T) {

// 	mockMenuRepo := &MockMenuRepository{}

// 	menuUsecase := menu.NewMenuUsecaseImpl(mockMenuRepo)

// 	userId := uint(1)
// 	menuId := uint(2)

// 	// Test when menu is found in repository
// 	mockMenuRepo.On("GetMenuById", menuId).Return(&entity.Menu{}, nil).Once()
// 	mockMenuRepo.On("ToggleFavoriteMenu", userId, menuId).Return(&entity.MenuFavorite{}, nil).Once()

// 	favorite, err := menuUsecase.ToggleFavoriteMenu(userId, menuId)

// 	if err != nil {
// 		t.Errorf("Expected no error, got %v", err)
// 	}

// 	if favorite == nil {
// 		t.Errorf("Expected favorite to be not nil")
// 	}

// 	// Test when menu is not found in repository
// 	mockMenuRepo.On("GetMenuById", menuId).Return(nil, domain.ErrMenuNotFound).Once()

// 	_, err = menuUsecase.ToggleFavoriteMenu(userId, menuId)

// 	if err != domain.ErrMenuNotFound {
// 		t.Errorf("Expected error to be ErrMenuNotFound, got %v", err)
// 	}

// 	mockMenuRepo.AssertExpectations(t)
// }
