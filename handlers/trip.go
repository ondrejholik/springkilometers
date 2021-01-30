package springkilometers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	models "github.com/ondrejholik/springkilometers/models"
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
		if trip, err := models.GetTripByID(tripID); err == nil {
			// Call the render function with the title, article and the name of the
			// template
			log.Println(trip)
			Render(c, gin.H{
				"title":   trip.Name,
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

	name := c.PostForm("name")
	content := c.PostForm("content")
	kilometersCount := c.PostForm("km")
	withbike := c.PostForm("withbike")

	session := sessions.Default(c)
	username := session.Get("current_user")

	if a, err := models.CreateNewTrip(username.(string), name, content, kilometersCount, withbike); err == nil {
		// If the article is created successfully, show success message
		Render(c, gin.H{
			"title":   "Submission Successful",
			"payload": a}, "trip-successful.html")
	} else {
		// if there was an error while creating the article, abort with an error
		c.AbortWithStatus(http.StatusBadRequest)
	}
}
