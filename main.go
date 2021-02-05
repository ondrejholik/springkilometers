package main

import (
	"github.com/gin-gonic/gin"
	models "github.com/ondrejholik/springkilometers/models"
  "log"
)

var router *gin.Engine

func main() {

	//err := godotenv.Load()
	//if err != nil {
	//log.Panic(err)
	//}

	//dburl := os.Getenv("DATABASE_URL")
	//user := os.Getenv("USERNAME")
	//dbname := os.Getenv("DATABASE")
	//dbpass := os.Getenv("PASSWORD")
	//dbport := os.Getenv("PORT")
	//dbhostname := os.Getenv("HOST")

  log.Println("Connecting to database")
	models.Setup()
  log.Println("Connected")
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
