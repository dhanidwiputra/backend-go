package usecase

import (
	"final-project-backend/util"
	"mime/multipart"

	"github.com/go-playground/validator/v10"
)

var (
	validate = validator.New()
)

type MediaUsecase interface {
	FileUpload(file multipart.FileHeader) (string, string, error)
	FileDelete(publicId string) error
	UploadGCS(file multipart.File, object string) (string, error)
}

type mediaUsecaseImpl struct {
	mediaUploader util.MediaUploader
	gcsUploader   util.GCSUploader
}

type MediaUsecaseConfig struct {
	MediaUploader util.MediaUploader
	GCSUploader   util.GCSUploader
}

func NewMediaUsecase(c MediaUsecaseConfig) MediaUsecase {
	return &mediaUsecaseImpl{
		mediaUploader: c.MediaUploader,
		gcsUploader:   c.GCSUploader,
	}
}

func (m *mediaUsecaseImpl) FileUpload(file multipart.FileHeader) (string, string, error) {
	err := validate.Struct(file)
	if err != nil {
		return "", "", err
	}

	fileContent, err := file.Open()
	if err != nil {
		return "", "", err
	}

	uploadUrl, publicId, err := m.mediaUploader.ImageUploadHelper(fileContent)
	if err != nil {
		return "", "", err
	}
	return uploadUrl, publicId, nil
}

func (m *mediaUsecaseImpl) FileDelete(publicId string) error {
	err := validate.Var(publicId, "required")
	if err != nil {
		return err
	}

	err = m.mediaUploader.ImageDeleteHelper(publicId)
	if err != nil {
		return err
	}
	return nil
}

func (m *mediaUsecaseImpl) UploadGCS(file multipart.File, object string) (string, error) {
	res, err := m.gcsUploader.UploadFile(file, object)
	if err != nil {
		return "", err
	}

	return res, nil
}
