package usecase

import (
	"final-project-backend/dto"
	"final-project-backend/repository"
	"final-project-backend/util"
)

type UserUsecase interface {
	GetUserByID(id uint) (*dto.UserResponse, error)
	GetUserByEmail(email string) (*dto.UserResponse, error)
	GetUserByUsername(username string) (*dto.UserResponse, error)
	GetUserByPhone(phone string) (*dto.UserResponse, error)
	UpdateUser(dto.UserResponse) (*dto.UserResponse, error)
	DeleteUserPhoto(id uint) error
	ResetGamesAttempt() error
}

type userUsecaseImpl struct {
	authUtil      util.AuthUtil
	userRepo      repository.UserRepository
	authUsecase   AuthUsecase
	mediaUploader util.MediaUploader
}

type UserUsecaseConfig struct {
	AuthUtil      util.AuthUtil
	AuthUsecase   AuthUsecase
	UserRepo      repository.UserRepository
	MediaUploader util.MediaUploader
}

func NewUserUsecase(c UserUsecaseConfig) UserUsecase {
	return &userUsecaseImpl{
		userRepo:      c.UserRepo,
		authUtil:      c.AuthUtil,
		authUsecase:   c.AuthUsecase,
		mediaUploader: c.MediaUploader,
	}
}

func (u *userUsecaseImpl) GetUserByID(id uint) (*dto.UserResponse, error) {
	user, err := u.userRepo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	userRes := dto.UserResponse{
		ID:              user.ID,
		Username:        user.Username,
		Email:           user.Email,
		Phone:           user.Phone,
		PictureUrl:      user.PictureUrl,
		FullName:        user.FullName,
		Role:            user.Role.Name,
		PicturePublicId: user.PicturePublicId,
		GamesAttempt:    user.GamesAttempt,
	}

	return &userRes, nil
}

func (u *userUsecaseImpl) GetUserByEmail(email string) (*dto.UserResponse, error) {
	user, err := u.userRepo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	userRes := dto.UserResponse{
		ID:              user.ID,
		Username:        user.Username,
		Email:           user.Email,
		Phone:           user.Phone,
		PictureUrl:      user.PictureUrl,
		FullName:        user.FullName,
		Role:            user.Role.Name,
		PicturePublicId: user.PicturePublicId,
		GamesAttempt:    user.GamesAttempt,
	}

	return &userRes, nil
}

func (u *userUsecaseImpl) GetUserByUsername(username string) (*dto.UserResponse, error) {
	user, err := u.userRepo.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	userRes := dto.UserResponse{
		ID:              user.ID,
		Username:        user.Username,
		Email:           user.Email,
		Phone:           user.Phone,
		PictureUrl:      user.PictureUrl,
		FullName:        user.FullName,
		Role:            user.Role.Name,
		PicturePublicId: user.PicturePublicId,
		GamesAttempt:    user.GamesAttempt,
	}

	return &userRes, nil
}

func (u *userUsecaseImpl) GetUserByPhone(phone string) (*dto.UserResponse, error) {
	user, err := u.userRepo.GetUserByPhone(phone)
	if err != nil {
		return nil, err
	}

	userRes := dto.UserResponse{
		ID:              user.ID,
		Username:        user.Username,
		Email:           user.Email,
		Phone:           user.Phone,
		PictureUrl:      user.PictureUrl,
		FullName:        user.FullName,
		Role:            user.Role.Name,
		PicturePublicId: user.PicturePublicId,
		GamesAttempt:    user.GamesAttempt,
	}

	return &userRes, nil
}

func (u *userUsecaseImpl) UpdateUser(input dto.UserResponse) (*dto.UserResponse, error) {
	userRes, err := u.userRepo.UpdateUser(input)
	if err != nil {
		return nil, err
	}

	return userRes, nil
}

func (u *userUsecaseImpl) DeleteUserPhoto(id uint) error {
	err := u.userRepo.DeletePhoto(id)
	if err != nil {
		return err
	}

	return nil
}

func (u *userUsecaseImpl) ResetGamesAttempt() error {
	err := u.userRepo.ResetGamesAttempt()
	if err != nil {
		return err
	}

	return nil
}
