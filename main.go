package main

import (
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
	"github.com/gin-gonic/gin"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"
)

var dbxClient = files.New(dropbox.Config{Token: getToken()})

func getToken() string {
	token := os.Getenv("DBX_TOKEN")
	if token == "" {
		log.Panicln("DBX_TOKEN environment variable not set")
	}
	return token
}

func main() {
	router := gin.Default()

	router.StaticFile("/", "./index.html")

	router.GET("/ping", func(c *gin.Context) {
		log.Println("pong")
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.POST("/upload", func(c *gin.Context) {
		log.Println("Processing request.")
		name := c.PostForm("name")
		name = path.Clean(name)

		form, _ := c.MultipartForm()

		photos := form.File["photos[]"]
		for _, photo := range photos {
			filePath := path.Join("/Photos", name, path.Base(photo.Filename))
			file, _ := photo.Open()
			go uploadPhoto(filePath, file)
		}

		log.Println("Received", len(photos), "photos from", name)
		c.Redirect(http.StatusSeeOther, "/")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	err := router.Run(":" + port)
	if err != nil {
		log.Println("Error starting server: ", err.Error())
		return
	}
}

func uploadPhoto(filePath string, file multipart.File) {
	upload := files.NewUploadArg(filePath)
	upload.Autorename = true

	log.Println("Uploading photo to", filePath)
	_, err := dbxClient.Upload(upload, file)
	if err != nil {
		log.Println("Error uploading file:", err.Error())
		return
	}

	log.Println("Uploaded file:", filePath)
}
