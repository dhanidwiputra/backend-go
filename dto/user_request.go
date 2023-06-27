package dto

import "mime/multipart"

type UpdateUserRequest struct {
	FullName       string               `form:"full_name" binding:"required"`
	Phone          string               `form:"phone" binding:"required"`
	ProfilePicture multipart.FileHeader `form:"profile_picture" binding:"required"`
}
