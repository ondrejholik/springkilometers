package springkilometers

import (
	"context"
	"net/http"
	"time"

	mid "github.com/ondrejholik/springkilometers/middleware"
	models "github.com/ondrejholik/springkilometers/models"
)

func initializeRoutes() {

	// Use the setUserStatus middleware for every route to set a flag
	// indicating whether the request was from an authenticated user or not
	router.Use(setUserStatus())

	// Use db in context
	router.Use(SetDBMiddleware())

	// Handle the index route
	router.GET("/", showIndexPage)

	// Group user related routes together
	userRoutes := router.Group("/u")
	{
		// Handle the GET requests at /u/login
		// Show the login page
		// Ensure that the user is not logged in by using the middleware
		userRoutes.GET("/login", mid.EnsureNotLoggedIn(), showLoginPage)

		// Handle POST requests at /u/login
		// Ensure that the user is not logged in by using the middleware
		userRoutes.POST("/login", mid.EnsureNotLoggedIn(), performLogin)

		// Handle GET requests at /u/logout
		// Ensure that the user is logged in by using the middleware
		userRoutes.GET("/logout", mid.EnsureLoggedIn(), logout)

		// Handle the GET requests at /u/register
		// Show the registration page
		// Ensure that the user is not logged in by using the middleware
		userRoutes.GET("/register", mid.EnsureNotLoggedIn(), showRegistrationPage)

		// Handle POST requests at /u/register
		// Ensure that the user is not logged in by using the middleware
		userRoutes.POST("/register", mid.EnsureNotLoggedIn(), register)
	}

	tripRoutes := router.Group("/trip")
	{
		tripRoutes.GET("/view/:trip_id", getTrip)
		tripRoutes.GET("/all", showTripsPage)
		tripRoutes.GET("/create", mid.EnsureLoggedIn(), showTripCreationPage)
		tripRoutes.POST("/create", mid.EnsureLoggedIn(), createTrip)
	}

}

// SetDBMiddleware --
func SetDBMiddleware(next http.Handler) http.Handler {
	db := models.ConnectToDB()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timeoutContext, _ := context.WithTimeout(context.Background(), time.Second)
		ctx := context.WithValue(r.Context(), "DB", db.WithContext(timeoutContext))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
