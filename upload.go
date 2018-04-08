package imageupload

import (
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/nfnt/resize"
)

// UploadFile function is a simple helper function that uploads and saves an image on the server
// Params:
// r: to get the picture from multi-part form using key "get_picture"
// location: the path on server where you wish to save file. Ex: /users/images/
// ID: unique string ID for the image
// size: to resize the image, the function will keep the aspect ratio intact
func UploadFile(r *http.Request, location string, ID string, size uint) (string, error) {
	var path string
	file, hdr, err := r.FormFile("get_picture")
	if err != nil {
		return path, nil
	}

	ext := string(hdr.Filename[len(hdr.Filename)-3:])
	if _, err := checkFileExtension(ext); err != nil {
		return "", err
	}
	defer file.Close()

	path, err = saveFile(file, location, ID, ext, size)
	if err != nil {
		return "", err
	}

	return path, nil
}

// CheckFileExtension function checks whether the uploaded file is an image or not
func checkFileExtension(ext string) (int8, error) {
	switch ext {
	case "jpg":
		fallthrough
	case "JPG":
		return JPG, nil
	case "png":
		fallthrough
	case "PNG":
		return PNG, nil
	case "gif":
		fallthrough
	case "GIF":
		return GIF, nil
	default:
		log.Println("File is NOT an image")
		return 0, errors.New("File is NOT an image")
	}
}

// SaveFile function helps in uploading the profile picture of user
func saveFile(src multipart.File, location, id, ext string, size uint) (string, error) {
	name := id + ".jpg"
	path := "." + location + name
	var img image.Image
	var op jpeg.Options
	op.Quality = 50

	dst, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	e, _ := checkFileExtension(ext)
	switch e {
	case JPG:
		img, err = decodeJPG(src, size)
		if err != nil {
			return "", err
		}
	case PNG:
		img, err = decodePNG(src, size)
		if err != nil {
			return "", err
		}
	case GIF:
		img, err = decodeGIF(src, size)
		if err != nil {
			return "", err
		}
	}

	if err := jpeg.Encode(dst, img, &op); err != nil {
		return "", err
	}

	path = location + name
	return path, err
}

// DecodeJPG function decodes JPG image
func decodeJPG(src multipart.File, size uint) (image.Image, error) {
	img, err := jpeg.Decode(src)
	img = resize.Resize(size, 0, img, resize.Lanczos3)
	return img, err
}

// DecodePNG function decodes PNG image
func decodePNG(src multipart.File, size uint) (image.Image, error) {
	img, err := png.Decode(src)
	img = resize.Resize(size, 0, img, resize.Lanczos3)
	return img, err
}

// DecodeGIF function decodes GIF image
func decodeGIF(src multipart.File, size uint) (image.Image, error) {
	img, err := gif.Decode(src)
	img = resize.Resize(size, 0, img, resize.Lanczos3)
	return img, err
}
