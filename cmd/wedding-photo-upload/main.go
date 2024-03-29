package main

import (
	"github.com/RyanBreaker/go-photo-upload/internal/logger"
	"github.com/RyanBreaker/go-photo-upload/internal/routes"
	"github.com/gin-gonic/gin"
	"log"
	"log/slog"
	"net/http"
	"os"
)

var isProduction = os.Getenv("ENV") == "production"

func main() {
	slog.SetDefault(logger.Logger)

	router := gin.Default()

	router.LoadHTMLFiles("static/index.html")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{"isProduction": isProduction})
	})

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
