package dto

type GameLeaderboardResponse struct {
	Username         string `json:"username"`
	AccumulatedScore int    `json:"accumulated_score"`
}
