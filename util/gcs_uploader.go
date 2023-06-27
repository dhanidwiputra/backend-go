package util

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"time"

	"cloud.google.com/go/storage"
)

type ClientUploader struct {
	cl         *storage.Client
	projectID  string
	bucketName string
	uploadPath string
}

const (
	projectID  = "group-project-shopee"       // FILL IN WITH YOURS
	bucketName = "final-project-stage-bucket" // FILL IN WITH YOURS
)

var clientUploader *ClientUploader

func init() {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "group-project-shopee-e9b78f49b28b.json") // FILL IN WITH YOUR FILE PATH
	client, err := storage.NewClient(context.Background())
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	clientUploader = &ClientUploader{
		cl:         client,
		bucketName: bucketName,
		projectID:  projectID,
		uploadPath: "public_assets/images/",
	}

}

type GCSUploader interface {
	UploadFile(file multipart.File, object string) (string, error)
}

type gCSUploaderImpl struct {
}

type GCSUploaderConfig struct {
}

func NewGCSUploader() GCSUploader {
	return &gCSUploaderImpl{}
}

// UploadFile uploads an object
func (g *gCSUploaderImpl) UploadFile(file multipart.File, object string) (string, error) {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()
	randString := RandomFileName(10)
	object = randString + RemoveSpaces(object)

	// Upload an object with storage.Writer.
	wc := clientUploader.cl.Bucket(clientUploader.bucketName).Object(clientUploader.uploadPath + object).NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return "", fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("Writer.Close: %v", err)
	}

	// [END upload_file]
	return fmt.Sprintf("https://storage.cloud.google.com/%s/%s", clientUploader.bucketName, clientUploader.uploadPath+object), nil
}
