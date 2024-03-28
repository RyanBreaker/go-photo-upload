package main

import (
	"github.com/RyanBreaker/go-photo-upload/internal/routes"
	"github.com/gin-gonic/gin"
	"log"
	"log/slog"
	"os"
)

func main() {
	log.SetOutput(os.Stdout)

	router := gin.Default()

	router.StaticFile("/", "./static/index.html")

	router.GET("/ping", func(c *gin.Context) {
		slog.Info("pong")
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	routes.OauthRoutes(router)
	routes.UploadRoute(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	err := router.Run(":" + port)
	if err != nil {
		log.Fatalln("Error starting server:", err.Error())
	}
}
