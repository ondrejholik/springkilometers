package main

import (
	"github.com/gin-gonic/gin"
	models "github.com/ondrejholik/springkilometers/models"
)

var router *gin.Engine

func main() {

	models.Setup()
	// Set Gin to production mode
	gin.SetMode(gin.ReleaseMode)

	// Set the router as the default one provided by Gin
	router = gin.Default()

	// Process the templates at the start so that they don't have to be loaded
	// from the disk again. This makes serving HTML pages very fast.
	router.LoadHTMLGlob("templates/*")

	// Initialize the routes
	initializeRoutes()

	// Start serving the application
	router.Run()
}
