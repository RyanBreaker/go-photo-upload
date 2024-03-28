package routes

import (
	dbx "github.com/RyanBreaker/go-photo-upload/internal/dropbox"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"path"
	"strings"
)

func UploadRoute(router *gin.Engine) {
	router.POST("/upload", func(c *gin.Context) {
		slog.Info("Processing request...")

		name := path.Clean(c.PostForm("name"))

		form, err := c.MultipartForm()
		if err != nil {
			if strings.Contains(err.Error(), "NextPart: EOF") {
				slog.Warn("Unexpected EOF, likely due to canceled transfer")
				return
			}
			slog.Error("Error processing multipart form", slog.Any("error", err))
			return
		}

		files := form.File["photos[]"]

		var photos []dbx.Photo
		for _, file := range files {
			data, _ := file.Open()
			photo := dbx.Photo{
				FilePath: path.Join("/Photos", name, file.Filename),
				Data:     data,
			}
			photos = append(photos, photo)
		}

		go dbx.QueuePhotos(&photos)

		c.Redirect(http.StatusSeeOther, "/?uploaded=true")
	})
}
