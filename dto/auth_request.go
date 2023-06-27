package dto

import "mime/multipart"

type LoginRequest struct {
	Identifier string `json:"identifier" binding:"required"`
	Password   string `json:"password" binding:"required,min=8"`
}

type RegisterRequest struct {
	FullName       string               `form:"full_name" binding:"required"`
	Phone          string               `form:"phone" binding:"required"`
	Username       string               `form:"username" binding:"required,min=6,max=16"`
	Email          string               `form:"email" binding:"required,email"`
	Password       string               `form:"password" binding:"required,min=8,max=16"`
	ProfilePicture multipart.FileHeader `form:"profile_picture"`
}

type RegisterData struct {
	FullName       string `json:"full_name"`
	PhoneNumber    string `json:"phone_number"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	ProfilePicture string `json:"profile_picture"`
	PublicId       string `json:"public_id"`
}
