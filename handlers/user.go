// handlers.user.go

package springkilometers

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	cache "github.com/go-redis/cache/v8"
	models "github.com/ondrejholik/springkilometers/models"
)

// Err --
type Err struct {
	Code    int
	Message string
}

// JwtWrapper wraps the signing key and the issuer
type JwtWrapper struct {
	SecretKey       string
	Issuer          string
	ExpirationHours int64
}

// JwtClaim adds email as a claim to the token
type JwtClaim struct {
	UserID   int
	Username string
	jwt.StandardClaims
}

// GenerateToken --
func (j *JwtWrapper) GenerateToken(userID int, username string) (signedToken string, err error) {
	claims := &JwtClaim{
		UserID:   userID,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(j.ExpirationHours)).Unix(),
			Issuer:    j.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err = token.SignedString([]byte(j.SecretKey))
	if err != nil {
		return
	}

	return
}

// ValidateToken validates the jwt token
func (j *JwtWrapper) ValidateToken(signedToken string) (claims *JwtClaim, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JwtClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(j.SecretKey), nil
		},
	)

	if err != nil {
		return
	}

	claims, ok := token.Claims.(*JwtClaim)
	if !ok {
		err = errors.New("Couldn't parse claims")
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("JWT is expired")
		return
	}

	return

}

// ClaimsUser --
func ClaimsUser(c *gin.Context) (*JwtClaim, error) {
	var claims *JwtClaim
	jwtWrapper := JwtWrapper{
		SecretKey: os.Getenv("ACCESS_SECRET"),
		Issuer:    "AuthService",
	}

	cookie, err := c.Cookie("token")
	if err != nil {
		Render(c, gin.H{
			"message": err.Error(),
			"title":   "Unauthorized",
		}, "error.html")
		return nil, err
	}
	claims, err = jwtWrapper.ValidateToken(cookie)
	if err != nil {
		Render(c, gin.H{
			"message": err.Error(),
			"title":   "Unauthorized",
		}, "error.html")
		return nil, err
	}
	return claims, nil
}

// GetCurrentUser --
func GetCurrentUser(c *gin.Context) (*JwtClaim, error) {
	var claims *JwtClaim
	jwtWrapper := JwtWrapper{
		SecretKey: os.Getenv("ACCESS_SECRET"),
		Issuer:    "AuthService",
	}

	cookie, err := c.Cookie("token")
	if err != nil {
		return nil, err
	}
	claims, err = jwtWrapper.ValidateToken(cookie)
	if err != nil {
		return nil, err
	}

	return claims, nil
}

// NoRoute --
func NoRoute(c *gin.Context) {
	Render(c, gin.H{
		"title":   "Error",
		"payload": Err{Code: 404, Message: "Not found"}}, "error.html")
}

//ShowIndexPage --
func ShowIndexPage(c *gin.Context) {
	result := models.GetUsersScore()
	Render(c, gin.H{
		"title":   "Jarní Kilometry 2021",
		"payload": result}, "index.html")
}

// ShowUser --
func ShowUser(c *gin.Context) {

	if userID, err := strconv.Atoi(c.Param("id")); err == nil {

		var userpage models.UserPage

		if err := models.MyCache.Get(models.Ctx, "user:"+c.Param("id"), &userpage); err != nil {
			userpage = models.GetUserPage(userID)
			if err := models.MyCache.Set(&cache.Item{
				Ctx:   models.Ctx,
				Key:   "user:" + c.Param("id"),
				Value: userpage,
				TTL:   time.Hour,
			}); err != nil {
				panic(err)
			}

		}

		Render(c, gin.H{
			"title":   "User",
			"payload": userpage,
		}, "user.html")

	} else {
		// If an invalid article ID is specified in the URL, abort with an error
		c.AbortWithStatus(http.StatusNotFound)
		err := Err{Code: 404, Message: "Not found"}
		Render(c, gin.H{
			"message": err,
			"title":   "404 Not found",
		}, "error.html")
	}

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
	if claims, err := ClaimsUser(c); err == nil {
		result := models.GetUserTrips(claims.UserID)
		Render(c, gin.H{
			"title":   "My trips",
			"payload": result}, "user-trips.html")
	}
}

// MyTripsSuccess --
func MyTripsSuccess(c *gin.Context) {
	if claims, err := ClaimsUser(c); err == nil {
		result := models.GetUserTrips(claims.UserID)
		Render(c, gin.H{
			"title":   "Trip successfuly added",
			"payload": result}, "user-trips-success.html")

	}
}

// JoinTrip --
func JoinTrip(c *gin.Context) {
	// Check if the article ID is valid
	if tripID, err := strconv.Atoi(c.Param("id")); err == nil {
		// Check if the user exists

		if trip, err := models.GetTripByID(tripID); err == nil {

			claims, err := ClaimsUser(c)
			models.MyCache.Delete(models.Ctx, "user:"+strconv.Itoa(claims.UserID))
			models.MyCache.Delete(models.Ctx, "achievments:"+strconv.Itoa(claims.UserID))
			models.MyCache.Delete(models.Ctx, "trip:"+strconv.Itoa(tripID))
			if err == nil {
				hasUser := models.TripHasUser(tripID, claims.Username)
				models.UserJoinsTrip(claims.Username, *trip)
				models.TripJoinsUser(claims.Username, *trip)
				trip, _ = models.GetTripByIDWithUsers(tripID)
				Render(c, gin.H{
					"title":    "Successful joined trip",
					"isjoined": hasUser,
					"message":  "Successful joined trip",
					"payload":  trip}, "trip.html")
			}

		} else {
			c.AbortWithError(http.StatusNotFound, err)
			c.AbortWithStatus(http.StatusNotFound)
			err := Err{Code: 404, Message: "Not found"}
			Render(c, gin.H{
				"message": err,
				"title":   "Not found",
			}, "error.html")
		}

	} else {
		c.AbortWithStatus(http.StatusNotFound)
		c.AbortWithStatus(http.StatusNotFound)
		err := Err{Code: 404, Message: "Not found"}
		Render(c, gin.H{
			"message": err,
			"title":   "Not found",
		}, "error.html")
	}
}

// DisjoinTrip --
func DisjoinTrip(c *gin.Context) {
	if tripID, err := strconv.Atoi(c.Param("id")); err == nil {

		if trip, err := models.GetTripByID(tripID); err == nil {
			claims, err := ClaimsUser(c)
			models.MyCache.Delete(models.Ctx, "user:"+strconv.Itoa(claims.UserID))
			models.MyCache.Delete(models.Ctx, "trip:"+strconv.Itoa(tripID))
			models.MyCache.Delete(models.Ctx, "achievments:"+strconv.Itoa(claims.UserID))
			if err == nil {
				// TODO: replace with ID
				models.UserDisjoinsTrip(claims.Username, *trip)
				models.TripDisjoinsUser(claims.Username, *trip)
				trip, _ = models.GetTripByIDWithUsers(tripID)
				Render(c, gin.H{
					"title":   trip.Name,
					"message": "Trip disjoined!",
					"payload": trip}, "trip.html")
			}

		} else {
			c.AbortWithError(http.StatusNotFound, err)
			c.AbortWithStatus(http.StatusNotFound)
			err := Err{Code: 404, Message: "Not found"}
			Render(c, gin.H{
				"message": err,
				"title":   "Not found",
			}, "error.html")
		}

	} else {
		c.AbortWithStatus(http.StatusNotFound)
		c.AbortWithStatus(http.StatusNotFound)
		err := Err{Code: 404, Message: "Not found"}
		Render(c, gin.H{
			"message": err,
			"title":   "Not found",
		}, "error.html")
	}
}

// PerformLogin --
func PerformLogin(c *gin.Context) {
	// Obtain the POSTed username and password values
	username := c.PostForm("username")
	password := c.PostForm("password")

	// Check if the username/password combination is valid

	if userID, passed := models.IsUserValid(username, password); passed {
		// If the username/password is valid set the token in a cookie
		jwtWrapper := JwtWrapper{
			SecretKey:       os.Getenv("ACCESS_SECRET"),
			Issuer:          "AuthService",
			ExpirationHours: 24,
		}

		token, err := jwtWrapper.GenerateToken(userID, username)
		if err != nil {
			log.Panic(err)
		}

		c.SetCookie("token", token, 86400, "", "", false, true)
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

// Logout --
func Logout(c *gin.Context) {
	// Clear the cookie
	c.SetCookie("token", "", -1, "", "", false, true)

	// Redirect to the home page
	c.Redirect(http.StatusTemporaryRedirect, "/")
	c.Abort()
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

	if userID, err := models.RegisterNewUser(username, password); err == nil {
		// If the user is created, set the token in a cookie and log the user in
		jwtWrapper := JwtWrapper{
			SecretKey:       os.Getenv("ACCESS_SECRET"),
			Issuer:          "AuthService",
			ExpirationHours: 24,
		}

		token, err := jwtWrapper.GenerateToken(userID, username)
		if err != nil {
			log.Panic(err)
		}
		c.SetCookie("token", token, 10800, "", "", false, true)
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

// ShowSettings --
func ShowSettings(c *gin.Context) {
	claims, err := ClaimsUser(c)
	if err != nil {
		// TODO: error page
		log.Panic(err)
	}
	user, err := models.GetUserByID(claims.UserID)
	log.Println(user.Avatar)
	if err == nil {
		Render(c, gin.H{
			"payload": user,
			"title":   "Settings",
		}, "user-settings.html")
	}
}

// SettingsUpdate --
func SettingsUpdate(c *gin.Context) {
	// User can set
	// - profile picture
	//
	// -----NICE-TO-HAVE----------
	// - password
	// - room / group

	avatar := c.PostForm("avatar_input")
	claims, err := ClaimsUser(c)
	models.MyCache.Delete(models.Ctx, "user:"+strconv.Itoa(claims.UserID))

	if err == nil {

		if err := models.UpdateSettings(claims.UserID, avatar); err == nil {
			ShowSettings(c)
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
