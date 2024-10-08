package routers

import (
	"warkop-api/config"
	"warkop-api/handlers"
	"warkop-api/middleware"
	"warkop-api/repositories"
	"warkop-api/services"

	"github.com/gin-gonic/gin"
)

func CompRouter(api *gin.RouterGroup) {
	api.Use(middleware.ClientTracker(config.InitDB()))

	compRepository := repositories.NewComponentRepository(config.InitDB())
	compService := services.NewService(compRepository)
	compHandler := handlers.NewCompHandlers(compService)

	api.Static("/user/image", "./public/user/image")

	api.GET("/ping", compHandler.Ping)
	api.POST("/key/register", compHandler.GenerateAPIKey)

	userRouter := api.Group("/user")
	userRouter.Use(middleware.APIKeyAuth(config.InitDB()))
	{
		userRouter.POST("/register", compHandler.RegisterUser)
		userRouter.POST("/login", compHandler.LoginUser)
		userRouter.POST("/request-reset", compHandler.RequestResetPassword)
		userRouter.POST("/verify-reset", compHandler.VerifyResetPassword)
		userRouter.POST("/reset-password", compHandler.ResetPassword)

		verificationRouter := userRouter.Group("/verif")
		{
			verificationRouter.POST("/resend", compHandler.ResendVerification)
			verificationRouter.POST("/verify", compHandler.VerifyAccount)
		}

		profileRouter := userRouter.Group("/profile")
		profileRouter.Use(middleware.AuthMiddleware())
		{
			profileRouter.POST("/upload", compHandler.UploadUserProfile)
			profileRouter.GET("/get", compHandler.GetUserProfile)
		}
	}

	menuRouter := api.Group("/menu")
	menuRouter.Use(middleware.APIKeyAuth(config.InitDB()))
	menuRouter.Use(middleware.AuthMiddleware())
	{
		menuRouter.POST("/register", compHandler.RegisterMenu)
		menuRouter.GET("/getall", compHandler.GetAllMenu)
	}

	transactionRouter := api.Group("/transaction")
	transactionRouter.Use(middleware.APIKeyAuth(config.InitDB()))
	transactionRouter.Use(middleware.AuthMiddleware())
	{
		transactionRouter.POST("/register", compHandler.RegisterTransaction)
		transactionRouter.GET("/get", compHandler.GetTransaction)
		transactionRouter.GET("/history", compHandler.GetTransactionHistory)
	}

}
