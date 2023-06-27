package entity

import "gorm.io/gorm"

type GameLeaderboard struct {
	gorm.Model
	UserID           uint `json:"user_id"`
	AccumulatedScore int  `json:"accumulated_score"`
	User             User `json:"user"`
}
