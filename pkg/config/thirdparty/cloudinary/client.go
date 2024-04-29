package cloudinary

import (
	"github.com/cloudinary/cloudinary-go/v2"
	"os"
)

var Cloudinary *cloudinary.Cloudinary

func NewCloudinaryClient() error {
	name := os.Getenv("CLOUDINARY_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_SECRET_KEY")

	_, err := cloudinary.NewFromParams(name, apiKey, apiSecret)

	return err
}
