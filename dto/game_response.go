package dto

import "final-project-backend/entity"

type TriviaAPIResponse struct {
	Category         string   `json:"category"`
	CorrectAnswer    string   `json:"correctAnswer"`
	Difficulty       string   `json:"difficulty"`
	IncorrectAnswers []string `json:"incorrectAnswers"`
	Question         string   `json:"question"`
}

type TriviaAPIResponseDTO struct {
	ID       uint     `json:"id"`
	Question string   `json:"question"`
	Options  []string `json:"options"`
}

type GameResponse struct {
	TriviaAPIResponse TriviaAPIResponseDTO `json:"trivia_api_response"`
	Game              entity.Game          `json:"game"`
}
