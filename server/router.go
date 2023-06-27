package server

import (
	"final-project-backend/handler"
	"final-project-backend/middleware"
	"final-project-backend/usecase"
	"final-project-backend/util"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type RouterConfig struct {
	AuthUsecase      usecase.AuthUsecase
	UserUsecase      usecase.UserUsecase
	CouponUsecase    usecase.CouponUsecase
	MenuUsecase      usecase.MenuUsecase
	MediaUsecase     usecase.MediaUsecase
	CartUsecase      usecase.CartUsecase
	OrderUsecase     usecase.OrderUsecase
	DeliveryUsecase  usecase.DeliveryUsecase
	GameUsecase      usecase.GameUsecase
	PromotionUsecase usecase.PromotionUsecase
}

func NewRouter(c RouterConfig) *gin.Engine {
	r := gin.Default()
	r.NoRoute(func(ctx *gin.Context) {
		util.ResponseErrorJSON(ctx, "page not found", "NOT_FOUND", 404)
	})

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET, POST, PUT, DELETE"},
		AllowHeaders:     []string{"Content-Type", "Content-Length", "Accept-Encoding", "Authorization", "Cache-Control"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	h := handler.New(handler.HandlerConfig{
		AuthUsecase:      c.AuthUsecase,
		UserUsecase:      c.UserUsecase,
		CouponUsecase:    c.CouponUsecase,
		MenuUsecase:      c.MenuUsecase,
		MediaUsecase:     c.MediaUsecase,
		CartUsecase:      c.CartUsecase,
		OrderUsecase:     c.OrderUsecase,
		DeliveryUsecase:  c.DeliveryUsecase,
		GameUsecase:      c.GameUsecase,
		PromotionUsecase: c.PromotionUsecase,
	})

	v1 := r.Group("/api/v1")

	v1.POST("/login", h.Login)
	v1.POST("/register", h.Register)
	v1.Static("/docs", "swaggerui")
	v1.GET("/promotions", h.GetPromotions)
	v1.GET("/promotions/:id", h.GetPromotionById)
	v1.GET("/menus", h.GetMenus)
	v1.GET("/menus/:id", h.GetMenuById)
	v1.GET("categories", h.GetCategories)
	v1.POST("/upload", h.UploadImage)

	v1.Use(middleware.Authorize, h.HasValidToken)
	v1.GET("/payment-options", h.GetAllPaymentOptions)
	v1.GET("/user-details", h.GetUserDetails)
	v1.DELETE("/user/photos", h.DeleteUserPhoto)
	v1.PUT("/user-details", h.UpdateUserDetails)
	v1.POST("/carts", h.AddToCart)
	v1.GET("/carts", h.GetCartItems)
	v1.DELETE("/carts/:id", h.DeleteCartItem)
	v1.PUT("/carts/:id", h.UpdateCartItem)
	v1.GET("/user-coupons", h.GetUserCoupons)
	v1.POST("/orders", h.CreateOrder)
	v1.GET("/orders", h.GetAllOrders)
	v1.PUT("/menus/:id/favorites", h.ToggleFavoriteMenu)
	v1.GET("/menus/favorites", h.GetFavoriteMenus)
	v1.POST("/customer-reviews", h.CreateCustomerReview)
	v1.POST("/games", h.CreateGame)
	v1.PUT("/games/:id", h.AnswerGameQuestion)
	v1.GET("/game-leaderboards", h.GetGameLeaderboard)
	v1.POST("/promotions/:id/orders", h.CreatePromotionOrder)
	v1.DELETE("/carts", h.EmptyCart)

	v1.Use(middleware.AuthorizeAdmin)
	v1.GET("/users/:id", h.GetUserByID)
	v1.POST("/reset-game", h.ResetGamesAttempt)
	v1.GET("/coupons", h.GetCoupons)
	v1.POST("/coupons", h.CreateCoupon)
	v1.PUT("/coupons/:id", h.UpdateCoupon)
	v1.DELETE("/coupons/:id", h.DeleteCouponById)
	v1.POST("/menus", h.CreateMenu)
	v1.PUT("/menus/:id", h.UpdateMenu)
	v1.DELETE("/menus/:id", h.DeleteMenu)
	v1.PUT("/deliveries/:id", h.UpdateDelivery)
	v1.POST("promotions", h.CreatePromotion)
	v1.PUT("/promotions/:id", h.UpdatePromotion)
	v1.DELETE("/promotions/:id", h.DeletePromotion)
	v1.GET("/customer-reviews/:id", h.GetCustomerReviewsByMenuId)
	v1.GET("/orders/count", h.GetTransactionTotalByDate)
	return r
}
