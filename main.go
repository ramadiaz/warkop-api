package main

import (
	"log"
	"os"
	"warkop-api/config"
	"warkop-api/repositories"
	"warkop-api/routers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	// Initialize DB and run auto-migrations
	db := config.InitDB()
	repositories.NewComponentRepository(db)

	port := os.Getenv("PORT")
	environment := os.Getenv("ENVIRONMENT")

	r := gin.New()
	r.Use(gin.Logger())

	configCors := cors.DefaultConfig()
	configCors.AllowOrigins = []string{"*"}
	configCors.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	configCors.AllowHeaders = []string{"*"}
	configCors.ExposeHeaders = []string{"Content-Length"}
	configCors.AllowCredentials = true
	r.Use(cors.New(configCors))

	api := r.Group("/api")
	routers.CompRouter(api)

	if environment == "production" {
		host := "0.0.0.0"
		server := host + ":" + port
		err := r.Run(server)
		if err != nil {
			log.Fatal("Error starting the server: ", err)
		}
	} else if environment == "development" {
		host := "localhost"
		server := host + ":" + port
		err := r.Run(server)
		if err != nil {
			log.Fatal("Error starting the server: ", err)
		}
	} else {
		log.Fatal("ENV ERROR: {ENVIRONMENT} UNKNOWN")
	}
}
