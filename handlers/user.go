// handlers.user.go

package springkilometers

import (
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	models "github.com/ondrejholik/springkilometers/models"
)

//ShowIndexPage --
func ShowIndexPage(c *gin.Context) {
	Render(c, gin.H{
		"title": "Index",
	}, "index.html")
}

// ShowLoginPage --
func ShowLoginPage(c *gin.Context) {
	// Call the render function with the name of the template to render
	Render(c, gin.H{
		"title": "Login",
	}, "login.html")
}

// PerformLogin --
func PerformLogin(c *gin.Context) {
	// Obtain the POSTed username and password values
	username := c.PostForm("username")
	password := c.PostForm("password")

	// Check if the username/password combination is valid

	if models.IsUserValid(username, password) {
		// If the username/password is valid set the token in a cookie
		token := GenerateSessionToken()
		c.SetCookie("token", token, 3600, "", "", false, true)
		c.Set("is_logged_in", true)

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

	if _, err := models.RegisterNewUser(username, password); err == nil {
		// If the user is created, set the token in a cookie and log the user in
		token := GenerateSessionToken()
		c.SetCookie("token", token, 3600, "", "", false, true)
		c.Set("is_logged_in", true)

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
