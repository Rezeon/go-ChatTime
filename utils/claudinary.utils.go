package utils

import (
	"context"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// untuk upload image ke claudinary
func UploadImage(filePath string, folder string) (string, string, error) {
	cld, err := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL")) // untuk menyimpan url claudinary
	if err != nil {                                                // jika error akan di kembalikan
		return "", "", err
	}
	//untuk mengupload ke claudinary
	resp, err := cld.Upload.Upload(context.Background(), filePath, uploader.UploadParams{
		UniqueFilename: api.Bool(false),
		Overwrite:      api.Bool(true),
		Folder:         folder,
	})
	if err != nil {
		return "", "", err
	}
	//yang akan di kembalikan adalah url dan publikid
	return resp.SecureURL, resp.PublicID, nil
}

// berfungsi untuk menghapus dari claudinary
func DeleteFromCloudinary(publicID string) error {
	cld, err := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))
	if err != nil {
		return err
	}

	//yang diambil adalah publicid
	_, err = cld.Upload.Destroy(context.Background(), uploader.DestroyParams{PublicID: publicID})
	if err != nil {
		return err
	}

	return nil
}
