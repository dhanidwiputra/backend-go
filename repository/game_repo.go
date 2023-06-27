package repository

import (
	"encoding/json"
	"final-project-backend/dto"
	"final-project-backend/entity"
	"io/ioutil"
	"net/http"

	"gorm.io/gorm"
)

type GameRepository interface {
	CreateGame(entity.Game) (*entity.Game, *dto.TriviaAPIResponse, error)
	GetGameByID(uint) (*entity.Game, error)
	AnswerGameQuestion(entity.Game) (*entity.Game, error)
	GetGameLeaderboard(uint) (*entity.GameLeaderboard, error)
	GetGamesLeaderboard() ([]entity.GameLeaderboard, error)
	CreateGameLeaderboard(entity.GameLeaderboard) (*entity.GameLeaderboard, error)
	UpdateGameLeaderboard(entity.GameLeaderboard) (*entity.GameLeaderboard, error)
}

type gameRepositoryImpl struct {
	db *gorm.DB
}

type GameRepoConfig struct {
	DB *gorm.DB
}

func NewGameRepository(c GameRepoConfig) GameRepository {
	return &gameRepositoryImpl{db: c.DB}
}

func (r *gameRepositoryImpl) CreateGame(game entity.Game) (*entity.Game, *dto.TriviaAPIResponse, error) {
	res, err := http.Get("https://the-trivia-api.com/api/questions?limit=1")
	if err != nil {
		return nil, nil, err
	}

	responseData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, nil, err
	}

	var responseObject []dto.TriviaAPIResponse
	json.Unmarshal(responseData, &responseObject)

	game.Answer = responseObject[0].CorrectAnswer
	game.Difficulty = responseObject[0].Difficulty

	err = r.db.Create(&game).Error
	if err != nil {
		return nil, nil, err
	}

	return &game, &responseObject[0], nil
}

func (r *gameRepositoryImpl) GetGameByID(id uint) (*entity.Game, error) {
	var game entity.Game
	err := r.db.First(&game, id).Error

	if err != nil {
		return nil, err
	}

	return &game, nil
}

func (r *gameRepositoryImpl) AnswerGameQuestion(game entity.Game) (*entity.Game, error) {
	err := r.db.Save(&game).Error

	if err != nil {
		return nil, err
	}

	return &game, nil
}

func (r *gameRepositoryImpl) GetGameLeaderboard(userId uint) (*entity.GameLeaderboard, error) {
	var gameLeaderboard *entity.GameLeaderboard

	err := r.db.Find(&gameLeaderboard, "user_id = ?", userId).Error
	if err != nil {
		return nil, err
	}

	return gameLeaderboard, nil
}

func (r *gameRepositoryImpl) GetGamesLeaderboard() ([]entity.GameLeaderboard, error) {
	var gameLeaderboards []entity.GameLeaderboard

	err := r.db.Preload("User").Order("accumulated_score desc").Find(&gameLeaderboards).Error
	if err != nil {
		return nil, err
	}

	return gameLeaderboards, nil
}

func (r *gameRepositoryImpl) CreateGameLeaderboard(gameLeaderboard entity.GameLeaderboard) (*entity.GameLeaderboard, error) {
	err := r.db.Create(&gameLeaderboard).Error
	if err != nil {
		return nil, err
	}

	return &gameLeaderboard, nil
}

func (r *gameRepositoryImpl) UpdateGameLeaderboard(gameLeaderboard entity.GameLeaderboard) (*entity.GameLeaderboard, error) {
	err := r.db.Save(&gameLeaderboard).Error
	if err != nil {
		return nil, err
	}

	return &gameLeaderboard, nil
}
