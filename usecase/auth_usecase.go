package usecase

import (
	"errors"
	"final-project-backend/domain"
	"final-project-backend/dto"
	"final-project-backend/entity"
	"final-project-backend/repository"
	"final-project-backend/util"
)

const (
	userRoleName = "user"
)

type AuthUsecase interface {
	Login(dto.LoginRequest) (string, error)
	Register(dto.RegisterData) (*dto.UserResponse, error)
	HasValidToken(uint, string) bool
}

type authUsecaseImpl struct {
	authUtil util.AuthUtil
	userRepo repository.UserRepository
}

type AuthUsecaseConfig struct {
	AuthUtil util.AuthUtil
	UserRepo repository.UserRepository
}

func NewAuthUsecase(c AuthUsecaseConfig) AuthUsecase {
	return &authUsecaseImpl{
		authUtil: c.AuthUtil,
		userRepo: c.UserRepo,
	}
}

func (a *authUsecaseImpl) Login(data dto.LoginRequest) (string, error) {
	var user *entity.User
	var err error

	matchEmail, err := util.IsEmail(data.Identifier)
	if matchEmail {
		user, err = a.userRepo.GetUserByEmail(data.Identifier)
	}

	if errors.Is(err, domain.ErrInvalidEmail) {
		return "", domain.ErrInvalidEmail
	}

	matchPhone, err := util.IsPhone(data.Identifier)
	if matchPhone {
		user, err = a.userRepo.GetUserByPhone(data.Identifier)
	}

	if errors.Is(err, domain.ErrInvalidPhone) {
		return "", domain.ErrInvalidPhone
	}

	matchUsername, err := util.IsUsername(data.Identifier)
	if matchUsername {
		user, err = a.userRepo.GetUserByUsername(data.Identifier)
	}

	if errors.Is(err, domain.ErrInvalidUsername) {
		return "", domain.ErrInvalidUsername
	}

	if !matchEmail && !matchPhone && !matchUsername {
		return "", domain.ErrInvalidIdentifier
	}

	if err != nil {
		return "", domain.ErrInternalServer
	}
	if !a.authUtil.ComparePassword(user.Password, data.Password) {
		return "", domain.ErrInvalidPassword
	}
	var token string
	token, err = a.authUtil.GenerateAccessToken(user)

	if err != nil {
		return "", domain.ErrInternalServer
	}

	err = a.userRepo.InsertToken(user.ID, token)

	if err != nil {
		return "", domain.ErrInternalServer
	}

	return token, nil
}

func (a *authUsecaseImpl) Register(data dto.RegisterData) (*dto.UserResponse, error) {
	matchEmail, _ := util.IsEmail(data.Email)
	if !matchEmail {
		return nil, domain.ErrInvalidEmailFormat
	}
	user, _ := a.userRepo.GetUserByEmail(data.Email)
	if user != nil {
		return nil, domain.ErrDuplicateEmail
	}

	matchUsername, _ := util.IsUsername(data.Username)
	if !matchUsername {
		return nil, domain.ErrInvalidUsernameFormat
	}
	user, _ = a.userRepo.GetUserByUsername(data.Username)
	if user != nil {
		return nil, domain.ErrDuplicateUsername
	}

	matchPhone, _ := util.IsPhone(data.PhoneNumber)
	if !matchPhone {
		return nil, domain.ErrInvalidPhoneFormat
	}
	user, _ = a.userRepo.GetUserByPhone(data.PhoneNumber)
	if user != nil {
		return nil, domain.ErrDuplicatePhone
	}

	var err error
	plainPassword := data.Password
	data.Password, err = a.authUtil.HashPassword(data.Password)
	if err != nil {
		return nil, err
	}

	userRegisterData := entity.User{
		FullName:        data.FullName,
		Email:           data.Email,
		Phone:           data.PhoneNumber,
		Username:        data.Username,
		Password:        data.Password,
		RoleID:          1,
		AccessToken:     "",
		PictureUrl:      data.ProfilePicture,
		PicturePublicId: data.PublicId,
	}

	user, err = a.userRepo.CreateUser(userRegisterData)
	if err != nil {
		return nil, err
	}

	accessToken, err := a.Login(dto.LoginRequest{
		Identifier: user.Email,
		Password:   plainPassword,
	})

	if err != nil {
		return nil, err
	}

	registerResponse := dto.UserResponse{
		ID:              user.ID,
		FullName:        user.FullName,
		Email:           user.Email,
		Phone:           user.Phone,
		Username:        user.Username,
		AccessToken:     accessToken,
		Role:            userRoleName,
		PictureUrl:      user.PictureUrl,
		PicturePublicId: user.PicturePublicId,
	}

	return &registerResponse, nil
}

func (a *authUsecaseImpl) HasValidToken(id uint, token string) bool {
	user, err := a.userRepo.GetUserByID(id)
	if err != nil {
		return false
	}

	return user.AccessToken == token
}
