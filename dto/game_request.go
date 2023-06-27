package dto

type GameRequest struct {
	Answer string `json:"answer" binding:"required"`
	GameID uint   `json:"game_id"`
	UserID uint   `json:"user_id"`
}
