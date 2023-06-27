package usecase

import (
	"final-project-backend/domain"
	"final-project-backend/dto"
	"final-project-backend/entity"
	"final-project-backend/repository"

	"github.com/gofrs/uuid"
)

var (
	score_hard   = 10000
	score_medium = 5000
	score_easy   = 1000
)

type GameUsecase interface {
	CreateGame(entity.Game) (*entity.Game, *dto.TriviaAPIResponse, error)
	AnswerGameQuestion(dto.GameRequest) (*entity.Game, error)
	GetGamesLeaderboard() ([]entity.GameLeaderboard, error)
}

type gameUsecaseImpl struct {
	gameRepo      repository.GameRepository
	userRepo      repository.UserRepository
	userUsecase   UserUsecase
	couponRepo    repository.CouponRepository
	couponUsecase CouponUsecase
}

type GameUsecaseConfig struct {
	GameRepo      repository.GameRepository
	UserRepo      repository.UserRepository
	UserUsecase   UserUsecase
	CouponRepo    repository.CouponRepository
	CouponUsecase CouponUsecase
}

func NewGameUsecase(c GameUsecaseConfig) GameUsecase {
	return &gameUsecaseImpl{
		gameRepo:      c.GameRepo,
		userRepo:      c.UserRepo,
		userUsecase:   c.UserUsecase,
		couponRepo:    c.CouponRepo,
		couponUsecase: c.CouponUsecase,
	}
}

func (u *gameUsecaseImpl) CreateGame(game entity.Game) (*entity.Game, *dto.TriviaAPIResponse, error) {
	user, err := u.userUsecase.GetUserByID(game.UserID)
	if err != nil {
		return nil, nil, err
	}

	if user.GamesAttempt <= 0 {
		return nil, nil, domain.ErrNoGamesAttempt
	}

	gameRes, triviaReponse, err := u.gameRepo.CreateGame(game)
	if err != nil {
		return nil, nil, err
	}

	err = u.userRepo.ReduceGamesAttempt(user.ID)
	if err != nil {
		return nil, nil, err
	}

	return gameRes, triviaReponse, nil
}

func (u *gameUsecaseImpl) AnswerGameQuestion(input dto.GameRequest) (*entity.Game, error) {
	game, err := u.gameRepo.GetGameByID(input.GameID)
	if err != nil {
		return nil, err
	}
	if game == nil {
		return nil, domain.ErrGameNotFound
	}
	if game.Score != nil {
		return nil, domain.ErrGameAlreadyAnswered
	}

	game.Score = new(int)

	if (game.Answer != input.Answer) && (input.Answer != "") {
		return u.gameRepo.AnswerGameQuestion(*game)
	}

	switch game.Difficulty {
	case "hard":
		game.Score = &score_hard
	case "medium":
		game.Score = &score_medium
	case "easy":
		game.Score = &score_easy
	}

	user, err := u.userUsecase.GetUserByID(game.UserID)
	if err != nil {
		return nil, err
	}

	coupon := entity.Coupon{
		Code:         uuid.Must(uuid.NewV4()),
		Description:  "Game Prize Coupon",
		Discount:     *game.Score,
		IssuerID:     user.ID,
		Availability: true,
	}

	couponRes, err := u.couponRepo.CreateCoupon(coupon)
	if err != nil {
		return nil, err
	}

	_, err = u.couponUsecase.AssignCouponToUser(*couponRes, *user)
	if err != nil {
		return nil, err
	}

	game.CouponID = &couponRes.ID
	gameRes, err := u.gameRepo.AnswerGameQuestion(*game)
	if err != nil {
		return nil, err
	}

	gameLeaderboard, _ := u.gameRepo.GetGameLeaderboard(game.UserID)
	if gameLeaderboard == nil {
		gameLeaderboard.AccumulatedScore = *game.Score
		gameLeaderboard.UserID = game.UserID
		gameLeaderboard, _ = u.gameRepo.CreateGameLeaderboard(*gameLeaderboard)

		return gameRes, nil
	}

	gameLeaderboard.AccumulatedScore += *game.Score
	gameLeaderboard.UserID = game.UserID

	_, err = u.gameRepo.UpdateGameLeaderboard(*gameLeaderboard)
	if err != nil {
		return nil, err
	}

	return gameRes, nil
}

func (u *gameUsecaseImpl) GetGamesLeaderboard() ([]entity.GameLeaderboard, error) {
	gamesLeaderboard, err := u.gameRepo.GetGamesLeaderboard()
	if err != nil {
		return nil, err
	}

	return gamesLeaderboard, nil
}
