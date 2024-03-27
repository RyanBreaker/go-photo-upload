package routes

import (
	dbx "github.com/RyanBreaker/go-photo-upload/internal/dropbox"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"path"
)

func UploadRoute(router *gin.Engine) {
	router.POST("/upload", func(c *gin.Context) {
		log.Println("Processing request")
		name := c.PostForm("name")
		name = path.Clean(name)

		form, _ := c.MultipartForm()

		photos := form.File["photos[]"]
		log.Println("Received", len(photos), "photos from", name)

		var photosToUpload []dbx.Photo
		for _, photo := range photos {
			filePath := path.Join("/Photos", name, path.Base(photo.Filename))
			file, _ := photo.Open()
			photosToUpload = append(photosToUpload, dbx.Photo{
				FilePath: filePath,
				Data:     file,
			})
		}

		go dbx.QueuePhotos(photosToUpload)

		c.Redirect(http.StatusSeeOther, "/?uploaded=true")
	})
}
