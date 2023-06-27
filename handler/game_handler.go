package handler

import (
	"errors"
	"final-project-backend/domain"
	"final-project-backend/dto"
	"final-project-backend/entity"
	"final-project-backend/util"
	"math/rand"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateGame(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		util.ResponseErrorJSON(c, domain.ErrUnauthorized.Error(), "UNAUTHORIZED", 401)
		return
	}

	userId := user.(dto.UserResponse).ID
	gameRequest := entity.Game{
		UserID: userId,
	}

	game, triviaResponse, err := h.gameUsecase.CreateGame(gameRequest)
	if errors.Is(err, domain.ErrNoGamesAttempt) {
		util.ResponseErrorJSON(c, domain.ErrNoGamesAttempt.Error(), "NO_GAMES_ATTEMPT", 403)
		return
	}
	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInternalServer.Error(), "INTERNAL_SERVER_ERROR", 500)
		return
	}

	options := []string{triviaResponse.CorrectAnswer}
	for _, option := range triviaResponse.IncorrectAnswers {
		options = append(options, option)
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(options), func(i, j int) { options[i], options[j] = options[j], options[i] })

	triviaDTO := dto.TriviaAPIResponseDTO{
		ID:       game.ID,
		Question: triviaResponse.Question,
		Options:  options,
	}

	util.ResponseSuccesJSON(c, triviaDTO, 201)
}

func (h *Handler) AnswerGameQuestion(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		util.ResponseErrorJSON(c, domain.ErrUnauthorized.Error(), "UNAUTHORIZED", 401)
		return
	}

	var gameRequest dto.GameRequest
	if err := c.ShouldBindJSON(&gameRequest); err != nil {
		util.ResponseErrorJSON(c, domain.ErrInvalidRequest.Error(), "INVALID_REQUEST", 400)
		return
	}

	gameId := c.Param("id")
	gameIdInt, err := strconv.Atoi(gameId)
	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInvalidRequest.Error(), "INVALID_REQUEST", 400)
		return
	}

	userId := user.(dto.UserResponse).ID
	gameRequest.GameID = uint(gameIdInt)
	gameRequest.UserID = userId

	game, err := h.gameUsecase.AnswerGameQuestion(gameRequest)
	if errors.Is(err, domain.ErrGameNotFound) {
		util.ResponseErrorJSON(c, domain.ErrGameNotFound.Error(), "GAME_NOT_FOUND", 404)
		return
	}

	if errors.Is(err, domain.ErrUnauthorized) {
		util.ResponseErrorJSON(c, domain.ErrUnauthorized.Error(), "UNAUTHORIZED", 401)
		return
	}

	if errors.Is(err, domain.ErrGameAlreadyAnswered) {
		util.ResponseErrorJSON(c, domain.ErrGameAlreadyAnswered.Error(), "GAME_ALREADY_ANSWERED", 403)
		return
	}

	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInternalServer.Error(), "INTERNAL_SERVER_ERROR", 500)
		return
	}

	util.ResponseSuccesJSON(c, game, 200)
}

func (h *Handler) GetGameLeaderboard(c *gin.Context) {
	gamesLeaderboard, err := h.gameUsecase.GetGamesLeaderboard()
	if err != nil {
		util.ResponseErrorJSON(c, domain.ErrInternalServer.Error(), "INTERNAL_SERVER_ERROR", 500)
		return
	}

	dtoResp := []dto.GameLeaderboardResponse{}

	for _, game := range gamesLeaderboard {
		dtoResp = append(dtoResp, dto.GameLeaderboardResponse{
			AccumulatedScore: game.AccumulatedScore,
			Username:         game.User.Username,
		})
	}
	util.ResponseSuccesJSON(c, dtoResp, 200)
}
