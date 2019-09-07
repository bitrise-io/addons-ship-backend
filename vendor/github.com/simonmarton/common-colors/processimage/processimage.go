package processimage

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/nfnt/resize"
	"github.com/simonmarton/common-colors/calculator"
	"github.com/simonmarton/common-colors/color"
	"github.com/simonmarton/common-colors/models"
)

// FromURL ...
func FromURL(url string) ([]string, error) {
	calculator := calculator.New(models.CalculatorConfig{
		Algorithm:            "yiq",
		TransparencyTreshold: 10,
		IterationCount:       3,
		MinLuminance:         0.3,
		MaxLuminance:         0.9,
		DistanceThreshold:    20,
		MinSaturation:        0.3,
	})

	client := httpClient()

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	// Content-Type is broken currently on Bitrise
	// openImage(resp.Body, resp.Header.Get("Content-Type"))

	img, err := openImage(resp.Body, path.Ext(url))
	if err != nil {
		return nil, err
	}

	img = resizeImage(img, 32, 32)
	colors := colorsFromImage(img)

	commonColors, _ := calculator.GetCommonColors(colors)
	return calculator.GenrateGradientColors(commonColors), nil
}

func httpClient() *http.Client {
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	return &client
}

func openImage(file io.Reader, imageType string) (image.Image, error) {
	var img image.Image
	var err error

	// save original bytes for the future
	originalBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	file = ioutil.NopCloser(bytes.NewReader(originalBytes)) // reset file stream

	switch imageType {
	case ".jpg":
		fallthrough
	case ".jpeg":
		fallthrough
	case "image/jpg":
		fallthrough
	case "image/jpeg":
		img, err = jpeg.Decode(file)
		if err != nil {
			file = ioutil.NopCloser(bytes.NewReader(originalBytes)) // reset file stream
			img, err = png.Decode(file)
			if err != nil {
				return nil, err
			}
		}
	case ".png":
		fallthrough
	case "image/png":
		img, err = png.Decode(file)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("Not supported image format: %s", imageType)
	}

	return img, nil
}

func colorsFromImage(i image.Image) (result []color.Color) {
	b := i.Bounds()
	s := b.Size()

	for x := 0; x < s.X; x++ {
		for y := 0; y < s.Y; y++ {
			c := color.NewFromRGBA(i.At(x, y))

			result = append(result, c)
		}
	}

	return result
}

func resizeImage(img image.Image, width, height int) image.Image {
	bounds := img.Bounds()

	if bounds.Dx() > width || bounds.Dy() > height {
		return resize.Resize(uint(width), uint(height), img, resize.Lanczos3)
	}

	return img
}
