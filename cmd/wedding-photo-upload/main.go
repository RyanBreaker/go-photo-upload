package main

import (
	"github.com/RyanBreaker/go-photo-upload/internal/routes"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func main() {
	log.SetOutput(os.Stdout)

	router := gin.Default()

	router.StaticFile("/", "./static/index.html")

	router.GET("/ping", func(c *gin.Context) {
		log.Println("pong")
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
		log.Println("Error starting server:", err.Error())
		return
	}
}
