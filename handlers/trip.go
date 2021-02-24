package springkilometers

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/disintegration/imageorient"
	"github.com/gin-gonic/gin"
	cache "github.com/go-redis/cache/v8"
	"github.com/muesli/smartcrop"
	"github.com/muesli/smartcrop/nfnt"
	"github.com/nfnt/resize"
	models "github.com/ondrejholik/springkilometers/models"
	"github.com/otiai10/opengraph"
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
		// Check if the trip exists
		models.MyCache.Delete(models.Ctx, "trip:"+strconv.Itoa(tripID))

		if trip, err := models.GetTripByID(tripID); err == nil {
			// Logged user is the author of deleted trip
			claims, err := ClaimsUser(c)
			if err == nil && claims.UserID == trip.AuthorID {
				models.DeleteTripByID(trip.ID)
				MyTrips(c)
			} else {
				c.AbortWithError(http.StatusNotFound, err)
				log.Println(err)
			}

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

func getFileName(filename string) string {
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
		// Check if the trip exists
		models.MyCache.Delete(models.Ctx, "trip:"+strconv.Itoa(tripID))

		name := c.PostForm("name")
		content := c.PostForm("content")
		kilometersCount := c.PostForm("km")
		withbike := c.PostForm("withbike")

		claims, err := ClaimsUser(c)
		if err == nil {

			if _, err := models.UpdateTrip(tripID, claims.UserID, name, content, kilometersCount, withbike); err == nil {
				models.MyCache.Delete(models.Ctx, "user:"+strconv.Itoa(claims.UserID))
				MyTrips(c)
			} else {
				// if there was an error while creating the article, abort with an error
				c.AbortWithStatus(http.StatusBadRequest)
			}
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
		// Check if the trip exists

		var trip *models.Trip

		if err := models.MyCache.Get(models.Ctx, "trip:"+c.Param("id"), &trip); err != nil {
			trip, err = models.GetTripByIDWithUsers(tripID)
			if err != nil {
				c.AbortWithStatus(http.StatusNotFound)
				log.Println(err)
			}

			go models.MyCache.Set(&cache.Item{
				Ctx:   models.Ctx,
				Key:   "trip:" + c.Param("id"),
				Value: trip,
				TTL:   time.Hour,
			})
		}
		// Is logged user joined in current trip
		claims, err := GetCurrentUser(c)
		var userInfo *models.User = nil
		if err == nil {
			userInfo, err = models.GetUserByID(claims.UserID)

		}

		var hasUser bool
		if err == nil {
			hasUser = models.TripHasUser(tripID, claims.Username)
		} else {
			hasUser = false
		}

		// template
		Render(c, gin.H{
			"title":    trip.Name,
			"message":  "",
			"userinfo": userInfo,
			"isjoined": hasUser,
			"payload":  trip}, "trip.html")

	}
}

// CreateTrip --
func CreateTrip(c *gin.Context) {
	// Obtain the POSTed title and content values

	var newtrip *models.Trip
	var err error
	withgpx := true
	var gpxname string
	var gpxfile multipart.File

	name := c.PostForm("name")
	content := c.PostForm("content")
	kilometersCount := c.PostForm("km")
	withbike := c.PostForm("withbike")
	var mapycz string

	_, header, err := c.Request.FormFile("image")
	if err != nil {
		c.AbortWithError(500, http.ErrMissingFile)
	}
	_, gpx, err := c.Request.FormFile("gpx")
	if err != nil {
		withgpx = false
	} else {
		gpxname = getFileName(gpx.Filename)
		gpxfile, err = gpx.Open()
		if err != nil {
			withgpx = false
			log.Panic(err)
		}
	}

	file, err := header.Open()
	if err != nil {
		log.Panic(err)
	}

	// Get image file name
	imagename := getFileName(header.Filename)
	imagetype := header.Header["Content-Type"][0]

	// Get gpx file name

	if claims, err := ClaimsUser(c); err == nil {
		user, err := models.GetUserByID(claims.UserID)
		models.MyCache.Delete(models.Ctx, "user:"+strconv.Itoa(claims.UserID))
		if err != nil {
			c.AbortWithError(404, err)
		}
		if withgpx {
			filename := fmt.Sprintf("/static/gpx/%s.gpx", gpxname)

			if newtrip, err = models.CreateNewTrip(*user, name, content, kilometersCount, withbike, filename, ""); err == nil {
				// If the article is created successfully, show success message
				MyTripsSuccess(c)
				saveGpx(newtrip.ID, gpxname, gpxfile)

			} else {
				// if there was an error while creating the article, abort with an error
				c.AbortWithStatus(http.StatusBadRequest)
			}

		} else {
			mapycz = c.PostForm("mapycz")
			if mapycz != "" {
				ogp, err := opengraph.Fetch(mapycz)
				if err != nil {
					log.Fatal(err)
					mapycz = ""
				} else {
					mapycz = ogp.Image[0].URL
				}
			}

			if newtrip, err = models.CreateNewTrip(*user, name, content, kilometersCount, withbike, "", mapycz); err == nil {
				// If the article is created successfully, show success message
				MyTripsSuccess(c)

			} else {
				// if there was an error while creating the article, abort with an error
				c.AbortWithStatus(http.StatusBadRequest)
			}
		}

		// Begin compression of image on background
		go compression(*newtrip, file, imagename, imagetype)
	} else {
		Render(c, gin.H{
			"message": err,
			"title":   "Unauthorized",
		}, "error.html")
	}
}

func trimFirstRune(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
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

func saveGpx(tripid int, gpxname string, gpxfile multipart.File) {

	filename := fmt.Sprintf("./static/gpx/%s.gpx", gpxname)
	gpxsaved, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer gpxsaved.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, gpxfile); err != nil {
		log.Panic(err)
	}

	if err = ioutil.WriteFile(filename, buf.Bytes(), 0666); err != nil {
		log.Println(err)
	}

	go models.AddGpsToTrip(filename, tripid)
}
