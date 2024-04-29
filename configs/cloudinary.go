package configs

import (
	"github.com/cloudinary/cloudinary-go/v2"
	"log"
	"os"
)

var Cloudinary *cloudinary.Cloudinary

func InitializeCloudinary() {
	name := os.Getenv("CLOUDINARY_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_SECRET_KEY")

	var err error
	Cloudinary, err = cloudinary.NewFromParams(name, apiKey, apiSecret)
	if err != nil {
		log.Fatalf("Failed to initialize Cloudinary: %v", err)
	}
}
