package entity

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FullName        string    `json:"full_name"`
	Phone           string    `json:"phone"`
	Email           string    `json:"email"`
	Username        string    `json:"username"`
	Password        string    `json:"password"`
	GamesAttempt    int       `json:"games_attempt"`
	RegisteredAt    time.Time `json:"registered_at"`
	RoleID          uint      `json:"role_id"`
	Role            Role      `json:"role"`
	PictureUrl      string    `json:"picture_url"`
	PicturePublicId string    `json:"picture_public_id"`
	AccessToken     string    `json:"access_token"`
	Coupons         []Coupon  `gorm:"many2many:users_coupons;"`
}
