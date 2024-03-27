package dropbox

import (
	"github.com/RyanBreaker/go-photo-upload/internal/oauth"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
	"log"
	"mime/multipart"
	"strings"
	"sync"
	"time"
)

var uploadWG sync.WaitGroup

type Photo struct {
	FilePath string
	Data     multipart.File
}

func QueuePhotos(photos []Photo) {
	log.Println("Queueing", len(photos), "photos")

	// TODO: Counter for when files for given batch are done?
	uploadWG.Wait()
	for _, photo := range photos {
		uploadWG.Add(1)
		go uploadPhoto(photo)
	}
}

func uploadPhoto(photo Photo) {
	defer uploadWG.Done()

	dbxClient := files.New(dropbox.Config{
		Token: oauth.GetAccessToken(),
	})

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
