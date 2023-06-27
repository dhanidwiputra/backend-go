package dto

type UserResponse struct {
	ID              uint   `json:"id"`
	Email           string `json:"email,omitempty"`
	Phone           string `json:"phone,omitempty"`
	FullName        string `json:"full_name,omitempty"`
	Username        string `json:"username,omitempty"`
	GamesAttempt    int    `json:"games_attempt"`
	PictureUrl      string `json:"picture_url,omitempty"`
	PicturePublicId string `json:"public_id,omitempty"`
	Role            string `json:"role"`
	AccessToken     string `json:"access_token,omitempty"`
}
