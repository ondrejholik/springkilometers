package springkilometers

import (
	"net/http"
	"strconv"

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
		"title": "Create New Trip"}, "create-trip.html")
}

// GetTrip --
func GetTrip(c *gin.Context) {
	// Check if the article ID is valid
	if tripID, err := strconv.Atoi(c.Param("trip_id")); err == nil {
		// Check if the article exists
		if trip, err := models.GetTripByID(tripID); err == nil {
			// Call the render function with the title, article and the name of the
			// template
			Render(c, gin.H{
				"title":   trip.Title,
				"payload": trip}, "trip.html")

		} else {
			// If the article is not found, abort with an error
			c.AbortWithError(http.StatusNotFound, err)
		}

	} else {
		// If an invalid article ID is specified in the URL, abort with an error
		c.AbortWithStatus(http.StatusNotFound)
	}
}

// CreateTrip --
func CreateTrip(c *gin.Context) {
	// Obtain the POSTed title and content values
	title := c.PostForm("title")
	content := c.PostForm("content")
	kilometersCount := c.PostForm("kmc")
	// Slice of users
	//users := c.PostForm("users")

	if a, err := models.CreateNewTrip(title, content, kilometersCount); err == nil {
		// If the article is created successfully, show success message
		Render(c, gin.H{
			"title":   "Submission Successful",
			"payload": a}, "trip-successful.html")
	} else {
		// if there was an error while creating the article, abort with an error
		c.AbortWithStatus(http.StatusBadRequest)
	}
}
