package dropbox

import (
	"github.com/RyanBreaker/go-photo-upload/internal/oauth"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
	"io"
	"log/slog"
	"strings"
	"sync"
	"time"
)

var uploadWG sync.WaitGroup

type Photo struct {
	FilePath string
	Data     io.Reader
}

func QueuePhotos(photos *[]Photo) {
	amount := len(*photos)
	slog.Info("Queueing files", slog.Int("amount", amount))

	// Wait until all other uploads are done to start
	uploadWG.Wait()
	for _, photo := range *photos {
		uploadWG.Add(1)
		uploadPhoto(photo)
	}
	slog.Info("Done uploading files", slog.Int("amount", amount))
}

func uploadPhoto(photo Photo) {
	defer uploadWG.Done()

	dbxClient := files.New(dropbox.Config{
		Token: oauth.GetAccessToken(),
	})

	upload := files.NewUploadArg(photo.FilePath)
	upload.Autorename = true

	slog.Info("Uploading file", slog.String("destination", photo.FilePath))
	for {
		_, err := dbxClient.Upload(upload, photo.Data)

		if err == nil {
			break
		}

		if strings.Contains(err.Error(), "too_many") {
			slog.Warn("Rate limited, trying again in 3 seconds")
			time.Sleep(3 * time.Second)
		} else {
			slog.Error("Error uploading file:", slog.Any("error", err))
			return
		}
	}
}
