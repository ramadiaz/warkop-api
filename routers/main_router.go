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
}
