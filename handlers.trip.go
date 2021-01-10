// handlers.article.go

package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func showTripsPage(c *gin.Context) {
	trips := getAllTrips()

	// Call the render function with the name of the template to render
	render(c, gin.H{
		"title":   "Trips Page",
		"payload": trips}, "trips.html")
}

func showTripCreationPage(c *gin.Context) {
	// Call the render function with the name of the template to render
	render(c, gin.H{
		"title": "Create New Trip"}, "create-trip.html")
}

func getTrip(c *gin.Context) {
	// Check if the article ID is valid
	if tripID, err := strconv.Atoi(c.Param("trip_id")); err == nil {
		// Check if the article exists
		if trip, err := getTripByID(tripID); err == nil {
			// Call the render function with the title, article and the name of the
			// template
			render(c, gin.H{
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

func createTrip(c *gin.Context) {
	// Obtain the POSTed title and content values
	title := c.PostForm("title")
	content := c.PostForm("content")
	kilometersCount := c.PostForm("kmc")

	if a, err := createNewTrip(title, content, kilometersCount); err == nil {
		// If the article is created successfully, show success message
		render(c, gin.H{
			"title":   "Submission Successful",
			"payload": a}, "trip-successful.html")
	} else {
		// if there was an error while creating the article, abort with an error
		c.AbortWithStatus(http.StatusBadRequest)
	}
}
