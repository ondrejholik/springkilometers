package springkilometers

import (
	"crypto/sha1"
	"fmt"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	cache "github.com/go-redis/cache/v8"
	models "github.com/ondrejholik/springkilometers/models"
	"github.com/otiai10/opengraph"
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

// GetFileName -- generates renadom string
func GetFileName(filename string) string {
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
		go models.MyCache.Delete(models.Ctx, "trip:"+strconv.Itoa(tripID))

		var gpx *multipart.FileHeader
		var gpxname string
		var hasgpx bool = false

		name := c.PostForm("name")
		content := c.PostForm("content")
		kilometersCount := c.PostForm("km")
		withbike := c.PostForm("withbike")

		gpx, err = c.FormFile("gpx")
		if err != nil {
			gpxname = ""
		} else {
			gpxname, hasgpx = GpxHandling(gpx, tripID)
		}

		claims, err := ClaimsUser(c)
		if err == nil {
			models.UpdateTrip(tripID, claims.UserID, name, content, kilometersCount, withbike, gpxname, hasgpx)
			models.MyCache.Delete(models.Ctx, "user:"+strconv.Itoa(claims.UserID))
			MyTrips(c)
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
		gpxname = GetFileName(gpx.Filename)
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
	imagename := GetFileName(header.Filename)
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
