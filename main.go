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
	"strings"
	"sync"
	"time"
)

var dbxClient = files.New(dropbox.Config{Token: getToken()})
var uploadWG sync.WaitGroup

type Photo struct {
	FilePath string
	Data     multipart.File
}

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
		log.Println("Processing request")
		name := c.PostForm("name")
		name = path.Clean(name)

		form, _ := c.MultipartForm()

		photos := form.File["photos[]"]
		log.Println("Received", len(photos), "photos from", name)

		var photosToUpload []Photo
		for _, photo := range photos {
			filePath := path.Join("/Photos", name, path.Base(photo.Filename))
			file, _ := photo.Open()
			photosToUpload = append(photosToUpload, Photo{
				FilePath: filePath,
				Data:     file,
			})
		}

		go queuePhotos(photosToUpload)

		c.Redirect(http.StatusSeeOther, "/?uploaded=true")
	})

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

func queuePhotos(photos []Photo) {
	log.Println("Queueing", len(photos), "photos")

	uploadWG.Wait()
	for _, photo := range photos {
		uploadWG.Add(1)
		go uploadPhoto(photo)
	}
}

func uploadPhoto(photo Photo) {
	defer uploadWG.Done()

	upload := files.NewUploadArg(photo.FilePath)
	upload.Autorename = true

	log.Println("Uploading photo to", photo.FilePath)
	for {
		_, err := dbxClient.Upload(upload, photo.Data)

		if err == nil {
			break
		}

		if strings.Contains(err.Error(), "too_many") {
			log.Println("Rate limited, trying again in 3 seconds")
			time.Sleep(3 * time.Second)
		} else {
			log.Println("Error uploading file:", err.Error())
			return
		}
	}

	log.Println("Uploaded file:", photo.FilePath)
}
