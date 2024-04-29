package services

import (
	"context"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"sharePie-api/configs"
)

func UploadImage(path string, folder string) (string, error) {
	var ctx = context.Background()
	uploadResult, err := configs.Cloudinary.Upload.Upload(
		ctx,
		path,
		uploader.UploadParams{PublicID: folder,
			UniqueFilename: api.Bool(false),
			Overwrite:      api.Bool(true)})

	if err != nil {
		return "Could not upload your image", err
	}

	return uploadResult.SecureURL, err
}
