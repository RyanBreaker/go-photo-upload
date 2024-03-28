package routes

import (
	"bytes"
	"errors"
	dbx "github.com/RyanBreaker/go-photo-upload/internal/dropbox"
	"github.com/gin-gonic/gin"
	"io"
	"log/slog"
	"net/http"
	"path"
)

func UploadRoute(router *gin.Engine) {
	router.POST("/upload", func(c *gin.Context) {
		slog.Info("Processing request...")

		form, _ := c.Request.MultipartReader()

		var name string
		var photos []dbx.Photo
		for {
			part, err := form.NextPart()
			if err == io.EOF {
				break
			} else if errors.Is(err, io.ErrUnexpectedEOF) {
				slog.Warn("Error parsing multipart form, likely canceled transfer:", err)
				return
			}

			if part.FormName() == "name" {
				buf, _ := io.ReadAll(part)
				name = string(buf)
				continue
			} else if part.FileName() == "" {
				// Ignore any other fields
				continue
			}

			buf := &bytes.Buffer{}
			_, _ = io.Copy(buf, part)
			photo := dbx.Photo{
				FilePath: path.Join("/Photos", name, part.FileName()),
				Data:     buf,
			}
			photos = append(photos, photo)
		}

		go dbx.QueuePhotos(&photos)

		c.Redirect(http.StatusSeeOther, "/?uploaded=true")
	})
}
