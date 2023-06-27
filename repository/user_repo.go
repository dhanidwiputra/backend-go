package repository

import (
	"final-project-backend/domain"
	"final-project-backend/dto"
	"final-project-backend/entity"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByID(id uint) (*entity.User, error)
	GetUserByEmail(email string) (*entity.User, error)
	GetUserByUsername(username string) (*entity.User, error)
	GetUserByPhone(phone string) (*entity.User, error)
	CreateUser(user entity.User) (*entity.User, error)
	InsertToken(id uint, token string) error
	UpdateUser(input dto.UserResponse) (*dto.UserResponse, error)
	DeletePhoto(id uint) error
	HasValidToken(id uint, token string) bool
	ReduceGamesAttempt(userId uint) error
	ResetGamesAttempt() error
}

type userRepositoryImpl struct {
	db *gorm.DB
}

type UserRepoConfig struct {
	DB *gorm.DB
}

func NewUserRepository(c UserRepoConfig) UserRepository {
	return &userRepositoryImpl{db: c.DB}
}

func (r *userRepositoryImpl) GetUserByID(id uint) (*entity.User, error) {
	var user *entity.User
	res := r.db.Preload("Role").Find(&user, id)
	if res.RowsAffected == 0 {
		return nil, domain.ErrResourceNotFound
	}
	return user, nil
}

func (r *userRepositoryImpl) GetUserByEmail(email string) (*entity.User, error) {
	var user *entity.User
	res := r.db.Preload("Role").Where("email LIKE ?", email).First(&user)
	if res.RowsAffected == 0 {
		return nil, domain.ErrInvalidEmail
	}
	return user, nil
}

func (r *userRepositoryImpl) GetUserByUsername(username string) (*entity.User, error) {
	var user *entity.User
	res := r.db.Preload("Role").Where("username LIKE ?", username).First(&user)
	if res.RowsAffected == 0 {
		return nil, domain.ErrInvalidUsername
	}
	return user, nil
}

func (r *userRepositoryImpl) GetUserByPhone(phone string) (*entity.User, error) {
	var user *entity.User
	res := r.db.Preload("Role").Where("phone LIKE ?", phone).First(&user)
	if res.RowsAffected == 0 {
		return nil, domain.ErrInvalidPhone
	}
	return user, nil
}

func (r *userRepositoryImpl) CreateUser(user entity.User) (*entity.User, error) {
	err := r.db.Debug().Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		return nil
	})
	return &user, err
}

func (r *userRepositoryImpl) InsertToken(id uint, token string) error {
	var user *entity.User
	res := r.db.Model(&user).Where("id = ?", id).Update("access_token", token)
	if res.RowsAffected == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

func (r *userRepositoryImpl) UpdateUser(input dto.UserResponse) (*dto.UserResponse, error) {
	var user *entity.User
	res := r.db.Model(&user).Where("id = ?", input.ID).Updates(input)
	if res.RowsAffected == 0 {
		return nil, domain.ErrUserNotFound
	}
	return &input, nil
}

func (r *userRepositoryImpl) ReduceGamesAttempt(userId uint) error {
	var user *entity.User
	res := r.db.Model(&user).Where("id = ?", userId).Update("games_attempt", gorm.Expr("games_attempt - ?", 1))
	if res.RowsAffected == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

func (r *userRepositoryImpl) DeletePhoto(id uint) error {
	var user *entity.User
	res := r.db.Model(&user).Where("id = ?", id).Update("profile_picture", []byte{})
	if res.RowsAffected == 0 {
		return domain.ErrResourceNotFound
	}
	return nil
}

func (r *userRepositoryImpl) HasValidToken(id uint, token string) bool {
	var user *entity.User
	res := r.db.Model(&user).Where("id = ? AND access_token = ?", id, token).First(&user)
	if res.RowsAffected == 0 {
		return false
	}
	return true
}

func (r *userRepositoryImpl) ResetGamesAttempt() error {
	var user *[]entity.User
	res := r.db.Model(&user).Where("games_attempt != ?", 3).Update("games_attempt", 3)
	if res.RowsAffected == 0 {
		return domain.ErrUserNotFound
	}

	return nil

}
