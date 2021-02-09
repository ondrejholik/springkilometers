package main

import (
	"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/contrib/sessions"
	handlers "github.com/ondrejholik/springkilometers/handlers"
	mid "github.com/ondrejholik/springkilometers/middleware"
)

func initializeRoutes() {

	// Use the setUserStatus middleware for every route to set a flag
	// indicating whether the request was from an authenticated user or not
	router.Use(mid.SetUserStatus())
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Static("/static", "./static")
	router.Use(sessions.Sessions("mysession", sessions.NewCookieStore([]byte("secret"))))

	// Use db in context

	// Handle the index route
	router.GET("/", handlers.ShowIndexPage)

	// Group user related routes together
	userRoutes := router.Group("/u")
	{
		// Handle the GET requests at /u/login
		// Show the login page
		// Ensure that the user is not logged in by using the middleware
		userRoutes.GET("/login", mid.EnsureNotLoggedIn(), handlers.ShowLoginPage)

		// Handle POST requests at /u/login
		// Ensure that the user is not logged in by using the middleware
		userRoutes.POST("/login", mid.EnsureNotLoggedIn(), handlers.PerformLogin)

		// Handle GET requests at /u/logout
		// Ensure that the user is logged in by using the middleware
		userRoutes.GET("/logout", mid.EnsureLoggedIn(), handlers.Logout)

		// Handle the GET requests at /u/register
		// Show the registration page
		// Ensure that the user is not logged in by using the middleware
		userRoutes.GET("/register", mid.EnsureNotLoggedIn(), handlers.ShowRegistrationPage)

		// Handle POST requests at /u/register
		// Ensure that the user is not logged in by using the middleware
		userRoutes.POST("/register", mid.EnsureNotLoggedIn(), handlers.Register)

		userRoutes.GET("/trips", mid.EnsureLoggedIn(), handlers.MyTrips)

		userRoutes.GET("/view/:id", handlers.ShowUser)
		userRoutes.GET("/username/:username", handlers.ShowUserByUsername)
	}

	tripRoutes := router.Group("/trip")
	{
		tripRoutes.GET("/all", handlers.ShowTripsPage)
		tripRoutes.GET("/create", mid.EnsureLoggedIn(), handlers.ShowTripCreationPage)
		tripRoutes.GET("/join/:id", mid.EnsureLoggedIn(), handlers.GetTrip)
		tripRoutes.GET("/update/:id", mid.EnsureLoggedIn(), handlers.ShowTripUpdatePage)
		tripRoutes.GET("/view/:id", handlers.GetTrip)

		tripRoutes.POST("/create", mid.EnsureLoggedIn(), handlers.CreateTrip)
		tripRoutes.POST("/delete/:id", mid.EnsureLoggedIn(), handlers.DeleteTrip)
		tripRoutes.POST("/join/:id", mid.EnsureLoggedIn(), handlers.JoinTrip)
		tripRoutes.POST("/disjoin/:id", mid.EnsureLoggedIn(), handlers.DisjoinTrip)
		tripRoutes.POST("/update/:id", mid.EnsureLoggedIn(), handlers.UpdateTrip)
	}

	router.NoRoute(handlers.NoRoute)
}

// SetDBMiddleware --
/*
func SetDBMiddleware(next http.Handler) http.Handler {
	database := models.Setup()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timeoutContext, _ := context.WithTimeout(context.Background(), time.Second)
		ctx := context.WithValue(r.Context(), "DB", database.WithContext(timeoutContext))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
*/
