package dto

import "mime/multipart"

type MenuFormRequest struct {
	Name        string               `form:"name" binding:"required"`
	Description string               `form:"description" binding:"required"`
	Price       int                  `form:"price" binding:"required"`
	Picture     multipart.FileHeader `form:"picture"`
	Categories  []uint               `form:"categories" binding:"required"`
	MenuOptions string               `form:"menu_options,omitempty" binding:"required"`
}
