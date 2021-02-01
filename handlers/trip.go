package springkilometers

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/muesli/smartcrop"
	"github.com/muesli/smartcrop/nfnt"
	"github.com/nfnt/resize"
	models "github.com/ondrejholik/springkilometers/models"
	"golang.org/x/image/webp"
)

// ShowTripsPage --
func ShowTripsPage(c *gin.Context) {
	trips := models.GetTrips()

	// Call the render function with the name of the template to render
	Render(c, gin.H{
		"title":   "Trips Page",
		"payload": trips}, "trips.html")
}

// ShowTripCreationPage --
func ShowTripCreationPage(c *gin.Context) {
	// Call the render function with the name of the template to render
	Render(c, gin.H{
		"title": "Create New Trip"}, "trip-create.html")
}

// ShowTripUpdatePage --
func ShowTripUpdatePage(c *gin.Context) {
	// Call the render function with the name of the template to render
	if tripID, err := strconv.Atoi(c.Param("id")); err == nil {
		// Check if the article exists
		if trip, err := models.GetTripByID(tripID); err == nil {
			// Call the render function with the title, article and the name of the
			// template
			log.Println(trip)
			Render(c, gin.H{
				"title":   trip.Name,
				"payload": trip}, "trip-update.html")

		} else {
			// If the article is not found, abort with an error
			c.AbortWithError(http.StatusNotFound, err)
			log.Println(err)
		}

	} else {
		// If an invalid article ID is specified in the URL, abort with an error
		c.AbortWithStatus(http.StatusNotFound)
		log.Println(err)
	}
}

// DeleteTrip --
func DeleteTrip(c *gin.Context) {
	// Call the render function with the name of the template to render
	if tripID, err := strconv.Atoi(c.Param("id")); err == nil {
		// Check if the article exists
		if trip, err := models.GetTripByID(tripID); err == nil {
			models.DeleteTripByID(trip.ID)
			MyTrips(c)
		} else {
			// If the article is not found, abort with an error
			c.AbortWithError(http.StatusNotFound, err)
			log.Println(err)
		}

	} else {
		// If an invalid article ID is specified in the URL, abort with an error
		c.AbortWithStatus(http.StatusNotFound)
		log.Println(err)
	}
}

func getImageName(filename string) string {
	name := "trip"
	// Random part
	rand.Seed(time.Now().UnixNano())
	rnd := rand.Int()
	// Hash file name
	// we dont care about security here, it is just an image which can be accessed by everybody
	h := sha1.New()
	h.Write([]byte(filename))
	bs := h.Sum(nil)

	return fmt.Sprintf("%s_%d%x", name, rnd, bs)
}

// UpdateTrip --
func UpdateTrip(c *gin.Context) {

	if tripID, err := strconv.Atoi(c.Param("id")); err == nil {
		// Check if the article exists
		name := c.PostForm("name")
		content := c.PostForm("content")
		kilometersCount := c.PostForm("km")
		withbike := c.PostForm("withbike")

		session := sessions.Default(c)
		username := session.Get("current_user")

		if _, err := models.UpdateTrip(tripID, username.(string), name, content, kilometersCount, withbike); err == nil {
			MyTrips(c)
		} else {
			// if there was an error while creating the article, abort with an error
			c.AbortWithStatus(http.StatusBadRequest)
		}

	} else {
		// If an invalid article ID is specified in the URL, abort with an error
		c.AbortWithStatus(http.StatusNotFound)
		log.Println(err)
	}

}

// GetTrip --
func GetTrip(c *gin.Context) {
	// Check if the article ID is valid
	if tripID, err := strconv.Atoi(c.Param("id")); err == nil {
		// Check if the article exists
		if trip, err := models.GetTripByIDWithUsers(tripID); err == nil {
			// Call the render function with the title, article and the name of the
			// template
			Render(c, gin.H{
				"title":   trip.Name,
				"message": "",
				"payload": trip}, "trip.html")

		} else {
			// If the article is not found, abort with an error
			c.AbortWithError(http.StatusNotFound, err)
			log.Println(err)
		}

	} else {
		// If an invalid article ID is specified in the URL, abort with an error
		c.AbortWithStatus(http.StatusNotFound)
		log.Println(err)
	}
}

// CreateTrip --
func CreateTrip(c *gin.Context) {
	// Obtain the POSTed title and content values

	var newtrip *models.Trip
	var err error
	name := c.PostForm("name")
	content := c.PostForm("content")
	kilometersCount := c.PostForm("km")
	withbike := c.PostForm("withbike")

	_, header, err := c.Request.FormFile("image")
	if err != nil {
		c.AbortWithError(500, http.ErrMissingFile)
	}

	// TODO: Check image size (max 15MB)
	if header.Size > 15000000 {
		log.Println("Error image too big")
	}

	file, err := header.Open()
	if err != nil {
		log.Panic(err)
	}

	// Get file name x
	filename := getImageName(header.Filename)
	filetype := header.Header["Content-Type"][0]

	session := sessions.Default(c)
	username := session.Get("current_user")

	if newtrip, err = models.CreateNewTrip(username.(string), name, content, kilometersCount, withbike); err == nil {
		// If the article is created successfully, show success message
		MyTripsSuccess(c)

	} else {
		// if there was an error while creating the article, abort with an error
		c.AbortWithStatus(http.StatusBadRequest)
	}

	// Begin compression of image on background
	go compression(*newtrip, file, filename, filetype)
}

func trimFirstRune(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}

// GCD greates common divisor
func GCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
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
		image, err = jpeg.Decode(file)
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
	trip.Tiny = fmt.Sprintf("./static/tiny_%s.jpg", filename)
	trip.Small = fmt.Sprintf("./static/small_%s.jpg", filename)
	trip.Medium = fmt.Sprintf("./static/medium_%s.jpg", filename)
	trip.Large = fmt.Sprintf("./static/large_%s.jpg", filename)

	// Resize image to
	// Tiny 	->    80x 60
	tiny := resize.Resize(80, 60, image, resize.Bilinear)

	// Small 	->	 160x120
	small := resize.Resize(160, 120, image, resize.Bilinear)

	// Medium 	-> 	 640x480
	medium := resize.Resize(640, 480, image, resize.Bilinear)

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
	jpeg.Encode(&buf, tiny, nil)
	if err = ioutil.WriteFile(trip.Tiny, buf.Bytes(), 0666); err != nil {
		log.Println(err)
	}

	buf.Reset()
	jpeg.Encode(&buf, small, nil)
	if err = ioutil.WriteFile(trip.Small, buf.Bytes(), 0666); err != nil {
		log.Println(err)
	}

	buf.Reset()
	jpeg.Encode(&buf, medium, nil)
	if err = ioutil.WriteFile(trip.Medium, buf.Bytes(), 0666); err != nil {
		log.Println(err)
	}

	buf.Reset()
	jpeg.Encode(&buf, large, nil)
	if err = ioutil.WriteFile(trip.Large, buf.Bytes(), 0666); err != nil {
		log.Println(err)
	}

	trip.Tiny = trimFirstRune((trip.Tiny))
	trip.Small = trimFirstRune((trip.Small))
	trip.Medium = trimFirstRune((trip.Medium))
	trip.Large = trimFirstRune((trip.Large))

	// Save paths to database
	models.UpdateTripStruct(trip)
}
