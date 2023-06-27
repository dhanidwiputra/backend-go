package dto

import "final-project-backend/entity"

type MenuResponse struct {
	ID              uint                `json:"id"`
	Name            string              `json:"name"`
	Description     string              `json:"description"`
	Price           int                 `json:"price"`
	PictureUrl      string              `json:"picture_url"`
	AvgRating       float64             `json:"avg_rating"`
	UserRatingCount int                 `json:"user_rating_count"`
	MenuOptions     []entity.MenuOption `json:"menu_options"`
	Categories      []entity.Category   `json:"categories,omitempty"`
}
