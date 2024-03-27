package dropbox

import (
	"github.com/RyanBreaker/go-photo-upload/internal/oauth"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
	"io"
	"log"
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
	log.Println("Queueing", len(*photos), "files")

	n := 0
	// Wait until all other uploads are done to start
	uploadWG.Wait()
	for _, photo := range *photos {
		uploadWG.Add(1)
		uploadPhoto(photo)
		n++
	}
	log.Println("Uploaded", n, "files")
}

func uploadPhoto(photo Photo) {
	defer uploadWG.Done()

	dbxClient := files.New(dropbox.Config{
		Token: oauth.GetAccessToken(),
	})

	upload := files.NewUploadArg(photo.FilePath)
	upload.Autorename = true

	log.Println("Uploading file to", photo.FilePath)
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
}
