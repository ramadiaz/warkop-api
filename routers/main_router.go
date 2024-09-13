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

	api.GET("/ping", compHandler.Ping)
	api.POST("/key/register", compHandler.GenerateAPIKey)

	userRouter := api.Group("/user")
	userRouter.Use(middleware.APIKeyAuth(config.InitDB()))
	{
		userRouter.POST("/register", compHandler.RegisterUser)

		verificationRouter := userRouter.Group("/verif")
		verificationRouter.Use(middleware.AuthMiddleware())
		{
			verificationRouter.POST("/resend", compHandler.ResendVerification)
		}
	}
}
