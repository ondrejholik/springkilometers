package springkilometers

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"strconv"

	"github.com/disintegration/imageorient"
	"github.com/muesli/smartcrop"
	"github.com/muesli/smartcrop/nfnt"
	"github.com/nfnt/resize"
	models "github.com/ondrejholik/springkilometers/models"
	"golang.org/x/image/webp"
)

func getCropDimensions(img image.Image, width, height int) (int, int) {
	// if we don't have width or height set use the smaller image dimension as both width and height
	if width == 0 && height == 0 {
		bounds := img.Bounds()
		x := bounds.Dx()
		y := bounds.Dy()
		if x < y {
			width = x
			height = x
		} else {
			width = y
			height = y
		}
	}
	return width, height
}

func compression(trip models.Trip, file multipart.File, filename, filetype string) {
	var image image.Image
	var err error
	if filetype == "image/png" {
		image, err = png.Decode(file)
		if err != nil {
			log.Panic(err)
		}
	} else if filetype == "image/jpeg" {
		image, _, err = imageorient.Decode(file)
		if err != nil {
			log.Panic(err)
		}
	} else {
		image, err = webp.Decode(file)
		if err != nil {
			log.Panic(err)
		}
	}

	// Smart cropping
	image = crop(image, 1024, 768, true)

	// Put new paths to Trip struct
	trip.Tiny = fmt.Sprintf("./static/img/tiny_%s.jpg", filename)
	trip.Small = fmt.Sprintf("./static/img/small_%s.jpg", filename)
	trip.Medium = fmt.Sprintf("./static/img/medium_%s.jpg", filename)
	trip.Large = fmt.Sprintf("./static/img/large_%s.jpg", filename)

	// Resize image to
	// Tiny 	->    80x 60
	tiny := resize.Resize(80, 60, image, resize.Bilinear)

	// Small 	->	 160x120
	small := resize.Resize(160, 120, image, resize.Bilinear)

	// Medium 	-> 	896x672
	medium := resize.Resize(896, 672, image, resize.Bilinear)

	// Large 	-> 	1024x768
	large := resize.Resize(1024, 768, image, resize.Bilinear)

	tinyfile, err := os.Create(trip.Tiny)
	if err != nil {
		log.Fatal(err)
	}
	defer tinyfile.Close()

	smallfile, err := os.Create(trip.Small)
	if err != nil {
		log.Fatal(err)
	}
	defer smallfile.Close()

	mediumfile, err := os.Create(trip.Medium)
	if err != nil {
		log.Fatal(err)
	}
	defer mediumfile.Close()

	largefile, err := os.Create(trip.Large)
	if err != nil {
		log.Fatal(err)
	}
	defer largefile.Close()

	var buf bytes.Buffer
	jpeg.Encode(&buf, tiny, &jpeg.Options{Quality: 95})
	if err = ioutil.WriteFile(trip.Tiny, buf.Bytes(), 0666); err != nil {
		log.Println(err)
	}

	buf.Reset()
	jpeg.Encode(&buf, small, &jpeg.Options{Quality: 90})
	if err = ioutil.WriteFile(trip.Small, buf.Bytes(), 0666); err != nil {
		log.Println(err)
	}

	buf.Reset()
	jpeg.Encode(&buf, medium, &jpeg.Options{Quality: 80})
	if err = ioutil.WriteFile(trip.Medium, buf.Bytes(), 0666); err != nil {
		log.Println(err)
	}

	buf.Reset()
	jpeg.Encode(&buf, large, &jpeg.Options{Quality: 75})
	if err = ioutil.WriteFile(trip.Large, buf.Bytes(), 0666); err != nil {
		log.Println(err)
	}

	trip.Tiny = trimFirstRune((trip.Tiny))
	trip.Small = trimFirstRune((trip.Small))
	trip.Medium = trimFirstRune((trip.Medium))
	trip.Large = trimFirstRune((trip.Large))

	// Save paths to database
	models.UpdateTripStruct(trip)
	models.MyCache.Delete(models.Ctx, "trip:"+strconv.Itoa(trip.ID))
}

func crop(img image.Image, w, h int, resize bool) image.Image {
	width, height := getCropDimensions(img, w, h)
	resizer := nfnt.NewDefaultResizer()
	analyzer := smartcrop.NewAnalyzer(resizer)
	topCrop, _ := analyzer.FindBestCrop(img, width, height)

	type SubImager interface {
		SubImage(r image.Rectangle) image.Image
	}
	img = img.(SubImager).SubImage(topCrop)
	if resize && (img.Bounds().Dx() != width || img.Bounds().Dy() != height) {
		img = resizer.Resize(img, uint(width), uint(height))
	}
	return img
}
