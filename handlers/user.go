// handlers.user.go

package springkilometers

import (
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	models "github.com/ondrejholik/springkilometers/models"
)

//ShowIndexPage --
func ShowIndexPage(c *gin.Context) {
	result := models.GetUsersScore()
	Render(c, gin.H{
		"title":   "Index",
		"payload": result}, "index.html")
}

// ShowLoginPage --
func ShowLoginPage(c *gin.Context) {
	// Call the render function with the name of the template to render
	Render(c, gin.H{
		"title": "Login",
	}, "login.html")
}

// MyTrips --
func MyTrips(c *gin.Context) {
	session := sessions.Default(c)
	currentUser := session.Get("current_user")
	result := models.GetUserTrips(currentUser.(string))
	Render(c, gin.H{
		"title":   "My trips",
		"payload": result}, "user-trips.html")

}

// JoinTrip --
func JoinTrip(c *gin.Context) {
	session := sessions.Default(c)
	// Check if the article ID is valid
	if tripID, err := strconv.Atoi(c.Param("id")); err == nil {
		// Check if the article exists
		if trip, err := models.GetTripByID(tripID); err == nil {
			currentUser := session.Get("current_user")
			log.Println(currentUser)

			models.UserJoinsTrip(currentUser.(string), *trip)
			models.TripJoinsUser(currentUser.(string), *trip)

			Render(c, gin.H{
				"title":   "Successful joined trip",
				"payload": trip}, "trip-joined.html")

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

// DisjoinTrip --
func DisjoinTrip(c *gin.Context) {
	session := sessions.Default(c)
	// Check if the article ID is valid
	if tripID, err := strconv.Atoi(c.Param("id")); err == nil {
		// Check if the article exists
		if trip, err := models.GetTripByID(tripID); err == nil {
			currentUser := session.Get("current_user")
			log.Println(currentUser)

			models.UserDisjoinsTrip(currentUser.(string), *trip)
			models.TripDisjoinsUser(currentUser.(string), *trip)
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

// PerformLogin --
func PerformLogin(c *gin.Context) {
	// Obtain the POSTed username and password values
	username := c.PostForm("username")
	password := c.PostForm("password")
	session := sessions.Default(c)

	// Check if the username/password combination is valid

	if models.IsUserValid(username, password) {
		// If the username/password is valid set the token in a cookie
		token := GenerateSessionToken()
		c.SetCookie("token", token, 3600, "", "", false, true)
		c.Set("is_logged_in", true)
		session.Set("current_user", username)
		session.Save()
		log.Println("settin:", username)

		Render(c, gin.H{
			"title": "Successful Login"}, "login-successful.html")

	} else {
		// If the username/password combination is invalid,
		// show the error message on the login page
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"ErrorTitle":   "Login Failed",
			"ErrorMessage": "Invalid credentials provided"})
	}
}

// GenerateSessionToken --
func GenerateSessionToken() string {
	// We're using a random 16 character string as the session token
	// This is NOT a secure way of generating session tokens
	// DO NOT USE THIS IN PRODUCTION
	// TODO: proper way to generate session token
	return strconv.FormatInt(rand.Int63(), 16)
}

// Logout --
func Logout(c *gin.Context) {
	// Clear the cookie
	c.SetCookie("token", "", -1, "", "", false, true)

	// Redirect to the home page
	c.Redirect(http.StatusTemporaryRedirect, "/")
}

// ShowRegistrationPage --
func ShowRegistrationPage(c *gin.Context) {
	// Call the render function with the name of the template to render
	Render(c, gin.H{
		"title": "Register"}, "register.html")
}

// Register --
func Register(c *gin.Context) {
	// Obtain the POSTed username and password values
	username := c.PostForm("username")
	password := c.PostForm("password")
	session := sessions.Default(c)

	if err := models.RegisterNewUser(username, password); err == nil {
		// If the user is created, set the token in a cookie and log the user in
		token := GenerateSessionToken()
		c.SetCookie("token", token, 3600, "", "", false, true)
		c.Set("is_logged_in", true)
		session.Set("current_user", username)
		session.Save()
		log.Println("settin:", username)

		Render(c, gin.H{
			"title": "Successful registration & Login"}, "login-successful.html")

	} else {
		// If the username/password combination is invalid,
		// show the error message on the login page
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"ErrorTitle":   "Registration Failed",
			"ErrorMessage": err.Error()})

	}
}

// Render one of HTML, JSON or CSV based on the 'Accept' header of the request
// If the header doesn't specify this, HTML is rendered, provided that
// the template name is present
func Render(c *gin.Context, data gin.H, templateName string) {
	loggedInInterface, _ := c.Get("is_logged_in")
	data["is_logged_in"] = loggedInInterface.(bool)

	switch c.Request.Header.Get("Accept") {
	case "application/json":
		// Respond with JSON
		c.JSON(http.StatusOK, data["payload"])
	case "application/xml":
		// Respond with XML
		c.XML(http.StatusOK, data["payload"])
	default:
		// Respond with HTML
		c.HTML(http.StatusOK, templateName, data)
	}
}
