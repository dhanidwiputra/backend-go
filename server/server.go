package server

import (
	"final-project-backend/db"
	"final-project-backend/repository"
	"final-project-backend/usecase"
	"final-project-backend/util"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func initRouter() *gin.Engine {
	userRepo := repository.NewUserRepository(repository.UserRepoConfig{
		DB: db.Get(),
	})

	couponRepo := repository.NewCouponRepository(repository.CouponRepoConfig{
		DB: db.Get(),
	})

	menuRepo := repository.NewMenuRepository(repository.MenuRepoConfig{
		DB: db.Get(),
	})

	cartRepo := repository.NewCartRepository(repository.CartRepoConfig{
		DB: db.Get(),
	})

	orderRepo := repository.NewOrderRepository(repository.OrderRepoConfig{
		DB: db.Get(),
	})

	paymentOptRepo := repository.NewPaymentOptionRepository(repository.PaymentOptionRepositoryConfig{
		DB: db.Get(),
	})

	deliveryRepo := repository.NewDeliveryRepository(repository.DeliveryRepoConfig{
		DB: db.Get(),
	})

	gameRepo := repository.NewGameRepository(repository.GameRepoConfig{
		DB: db.Get(),
	})

	promotionRepo := repository.NewPromotionRepository(repository.PromotionRepoConfig{
		DB: db.Get(),
	})

	mediaUploader := util.NewMediaUploaderUtil()
	gcsUploader := util.NewGCSUploader()
	mediaUsecase := usecase.NewMediaUsecase(usecase.MediaUsecaseConfig{
		MediaUploader: mediaUploader,
		GCSUploader:   gcsUploader,
	})

	userUsecase := usecase.NewUserUsecase(usecase.UserUsecaseConfig{
		AuthUtil:      util.NewAuthUtil(),
		UserRepo:      userRepo,
		MediaUploader: mediaUploader,
	})

	authUsecase := usecase.NewAuthUsecase(usecase.AuthUsecaseConfig{
		AuthUtil: util.NewAuthUtil(),
		UserRepo: userRepo,
	})

	couponUsecase := usecase.NewCouponUsecase(usecase.CouponUsecaseConfig{
		UserRepo:   userRepo,
		CouponRepo: couponRepo,
	})

	menuUsecase := usecase.NewMenuUsecase(usecase.MenuUsecaseConfig{
		MenuRepo:  menuRepo,
		OrderRepo: orderRepo,
	})

	cartUsecase := usecase.NewCartUsecase(usecase.CartUsecaseConfig{
		CartRepo: cartRepo,
	})

	orderUsecase := usecase.NewOrderUsecase(usecase.OrderUsecaseConfig{
		OrderRepo:      orderRepo,
		CartRepo:       cartRepo,
		MenuRepo:       menuRepo,
		PaymentOptRepo: paymentOptRepo,
		DeliveryRepo:   deliveryRepo,
	})

	deliveryUsecase := usecase.NewDeliveryUsecase(usecase.DeliveryUsecaseConfig{
		DeliveryRepo: deliveryRepo,
	})

	gameUsecase := usecase.NewGameUsecase(usecase.GameUsecaseConfig{
		GameRepo:      gameRepo,
		UserRepo:      userRepo,
		UserUsecase:   userUsecase,
		CouponRepo:    couponRepo,
		CouponUsecase: couponUsecase,
	})

	promotionUsecase := usecase.NewPromotionUsecase(usecase.PromotionUsecaseConfig{
		PromotionRepo:     promotionRepo,
		PaymentOptionRepo: paymentOptRepo,
		MediaUsecase:      mediaUsecase,
		MenuUsecase:       menuUsecase,
		CartUsecase:       cartUsecase,
		CouponUsecase:     couponUsecase,
		OrderUsecase:      orderUsecase,
	})

	r := NewRouter(RouterConfig{
		AuthUsecase:      authUsecase,
		UserUsecase:      userUsecase,
		CouponUsecase:    couponUsecase,
		MenuUsecase:      menuUsecase,
		MediaUsecase:     mediaUsecase,
		CartUsecase:      cartUsecase,
		OrderUsecase:     orderUsecase,
		DeliveryUsecase:  deliveryUsecase,
		GameUsecase:      gameUsecase,
		PromotionUsecase: promotionUsecase,
	})

	return r
}

func Init() {
	err := godotenv.Load(".env")

	if err != nil {
		panic(err)
	}

	dbErr := db.Connect()
	if dbErr != nil {
		fmt.Println("error connecting to DB")
	}

	r := initRouter()
	err = r.Run()

	if err != nil {
		fmt.Println("error while running server", err)
		return
	}
}
