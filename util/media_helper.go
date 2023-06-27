package util

import (
	"context"
	"final-project-backend/config"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type MediaUploader interface {
	ImageUploadHelper(input interface{}) (string, string, error)
	ImageDeleteHelper(publicID string) error
}

type mediaUploaderImpl struct{}

func NewMediaUploaderUtil() MediaUploader {
	return &mediaUploaderImpl{}
}

func (m *mediaUploaderImpl) ImageUploadHelper(input interface{}) (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cloudinaryConfig := config.InitConfig().CloudinaryConfig
	cld, err := cloudinary.NewFromParams(cloudinaryConfig.CloudName, cloudinaryConfig.APIKey, cloudinaryConfig.APISecret)
	if err != nil {
		return "", "", err
	}

	uploadParam, err := cld.Upload.Upload(ctx, input, uploader.UploadParams{})
	if err != nil {
		return "", "", err
	}
	return uploadParam.SecureURL, uploadParam.PublicID, nil
}

func (m *mediaUploaderImpl) ImageDeleteHelper(publicID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cloudinaryConfig := config.InitConfig().CloudinaryConfig
	cld, err := cloudinary.NewFromParams(cloudinaryConfig.CloudName, cloudinaryConfig.APIKey, cloudinaryConfig.APISecret)
	if err != nil {
		return err
	}

	_, err = cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	if err != nil {
		return err
	}
	return nil
}
