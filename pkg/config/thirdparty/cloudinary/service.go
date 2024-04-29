package cloudinary

import (
	"context"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func UploadImage(path string, folder string) (string, error) {
	var ctx = context.Background()
	uploadResult, err := Cloudinary.Upload.Upload(
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
