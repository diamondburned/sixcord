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

var imageCache = map[string]image.Image{}

// DownloadImage downloads the image. Supported formats are
// GIF, JPEG and PNG.
func DownloadImage(url string) (image.Image, error) {
	if i, ok := imageCache[url]; ok {
		return i, nil
	}

	r, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	img, _, err := image.Decode(r.Body)
	if err != nil {
		return nil, err
	}

	imageCache[url] = img

	return img, nil
}
