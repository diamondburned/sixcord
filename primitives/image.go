package primitives

import (
	"image"
	"net/http"
	"time"

	// Shadow imports for supported image types
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

var httpClient = &http.Client{
	Timeout: 2 * time.Second,
}

// DownloadImage downloads the image. Supported formats are
// GIF, JPEG and PNG.
func DownloadImage(url string) (image.Image, error) {
	r, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	img, _, err := image.Decode(r.Body)
	if err != nil {
		return nil, err
	}

	return img, nil
}
