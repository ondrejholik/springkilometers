// models.user.go

package springkilometers

import (
	"encoding/hex"
	"errors"
	"log"
	"math/rand"
	"strings"
	"time"

	"golang.org/x/crypto/sha3"
	"gorm.io/gorm"
)

// User --
type User struct {
	ID         int       `json:"id"`
	Username   string    `json:"username"`
	Password   string    `json:"password"`
	Salt       string    `json:"salt"	`
	CreatedOn  time.Time `json:"created_on"`
	DeletedOn  time.Time `json:"deleted_on"`
	ModifiedOn time.Time `json:"modified_on"`
	UpdatedOn  time.Time `json:"updated_on"`
	Trips      []Trip    `gorm:"many2many:user_trip;"`
}

// UserPage --
type UserPage struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	CreatedOn time.Time `json:"created_on"`
	Trips     []Trip    `gorm:"many2many:user_trip;"`
	Km        float64   `json:"km"`
	Kmbike    float64   `json:"kmbike"`
	Kmwalk    float64   `json:"kmwalk"`
	AvgKm     float64   `json:"avgkm"`
}

// Result of database query
type Result struct {
	ID        int
	Username  string
	Km        float64
	Avgkm     float64
	Tripcount int
}

// GetUsersScore --
func GetUsersScore() []Result {
	var result []Result
	db.Table("users").Select("users.id, users.username, SUM(trips.km) as km").Joins("JOIN user_trip ON users.id = user_trip.user_id").Joins("JOIN trips ON user_trip.trip_id = trips.id").Group("users.id, users.username").Order("km desc").Scan(&result)
	return result
}

// GetUserPage --
func GetUserPage(id int) UserPage {
	// ---------------//
	// - score
	// - km on bike
	// - km walking
	// - trips
	// ---------------//

	var userpage UserPage
	var result Result
	var user User
	db.Table("users").Select("users.id, users.username, AVG(trips.km) as avgkm, SUM(trips.km) as km, COUNT(trips.id) as tripcount").Joins("JOIN user_trip ON users.id = user_trip.user_id").Joins("JOIN trips ON user_trip.trip_id = trips.id").Group("users.id, users.username").First(&result, id)
	db.Table("users").Select("SUM(trips.km) as kmbike").Joins("JOIN user_trip ON users.id = user_trip.user_id").Joins("JOIN trips ON user_trip.trip_id = trips.id").Where("trips.withbike = ?", true).Group("users.id").First(&userpage, id)
	//db.Table("users").Select("users.id, users.username, trips.*").Joins("JOIN user_trip ON users.id = user_trip.user_id").Joins("JOIN trips ON user_trip.trip_id = trips.id").First(&userpage, id)
	db.Preload("Trips", func(db *gorm.DB) *gorm.DB {
		return db.Order("trips.timestamp DESC")
	}).First(&user, id)

	userpage.Km = result.Km
	userpage.AvgKm = result.Avgkm
	userpage.Kmwalk = userpage.Km - userpage.Kmbike
	userpage.Trips = user.Trips
	userpage.ID = user.ID
	userpage.Username = user.Username

	return userpage
}

// GetUserByUsername --
func GetUserByUsername(username string) (int, error) {
	var user User
	err := db.Table("users").Select("users.id").Where("users.username = ?", username).First(&user).Error
	if err != nil {
		return -1, err
	}
	return user.ID, nil
}

// random salt with given length
func salting(n int) string {
	rand.Seed(time.Now().UnixNano())
	var letterBytes = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

// sha3 hash algorithm
func crypting(s string) string {
	h := sha3.New512()
	h.Write([]byte(s))
	pass := h.Sum(nil)
	return hex.EncodeToString(pass)
}

// IsUserValid --
// Check if the username and password combination is valid
func IsUserValid(username, password string) bool {

	if isUsernameAvailable(username) {
		return false
	}

	var user User
	db.Where("username = ?", username).First(&user)

	var pass string = crypting(password + user.Salt)
	log.Println(pass == user.Password)

	return pass == user.Password
}

// UserJoinsTrip --
func UserJoinsTrip(username string, trip Trip) {
	var user User
	db.Where("username = ?", username).First(&user)
	db.Model(&user).Association("Trips").Append(&trip)
}

// UserDisjoinsTrip --
func UserDisjoinsTrip(username string, trip Trip) {
	var user User
	db.Where("username = ?", username).First(&user)
	db.Model(&user).Association("Trips").Delete(&trip)
}

// RegisterNewUser a new user with the given username and password
func RegisterNewUser(username, password string) error {
	var user User
	if strings.TrimSpace(password) == "" {
		return errors.New("The password can't be empty")
	} else if !isUsernameAvailable(username) {
		return errors.New("The username isn't available")
	}

	salt := salting(10)
	pass := crypting(password + salt)

	user = User{Username: username, Password: pass, Salt: salt}
	result := db.Create(&user) // pass pointer of data to Create
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// Check if the supplied username is available
func isUsernameAvailable(username string) bool {
	var count int64 = 0
	err := db.Table("users").Where("username = ?", username).Count(&count)
	if err != nil {
		log.Println(err)
	}
	if count > 0 {
		return false
	}
	return true
}
