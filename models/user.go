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
	Avatar     string    `json:"avatar"`
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
	ID            int       `json:"id"`
	Username      string    `json:"username"`
	Avatar        string    `json:"avatar"`
	CreatedOn     time.Time `json:"created_on"`
	Trips         []Trip    `gorm:"many2many:user_trip;"`
	Km            float64   `json:"km"`
	Km100         float64
	Kmbike        float64 `json:"kmbike"`
	Kmwalk        float64 `json:"kmwalk"`
	AvgKm         float64 `json:"avgkm"`
	Maxkm         float64 `json:"maxkm"`
	VillagesCount int     `json:"villages_count`
	TripCount     int     `json:"trip_count`
	Achievments   Achievments
	//Villages      []Village
	Villages []Village `gorm:"many2many:trip_village;"`
	Pois     []Poi     `gorm:"many2many:trip_poi;"`
}

// Result of database query
type Result struct {
	ID        int
	Username  string
	Avatar    string
	Km        float64
	Km100     float64
	Avgkm     float64
	Maxkm     float64
	Tripcount int
}

// Achievments --
type Achievments struct {
	Walker1 bool
	Walker2 bool
	Walker3 bool

	Explorer1 bool
	Explorer2 bool
	Explorer3 bool

	Challenger1 bool
	Challenger2 bool
	Challenger3 bool

	Score int
}

// GetUsersScore --
func GetUsersScore() []Result {
	var result []Result

	db.Table("users").Select("users.id, users.username, users.avatar, SUM(trips.km) as km").Joins("JOIN user_trip ON users.id = user_trip.user_id").Joins("JOIN trips ON user_trip.trip_id = trips.id").Group("users.id, users.username").Order("km desc").Scan(&result)
	for i := range result {
		result[i].Km100 = result[i].Km
		if result[i].Km > 200 {
			result[i].Km100 = result[i].Km - 200

		} else if result[i].Km > 100 {
			result[i].Km100 = result[i].Km - 100
		}

	}
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
	var achievments Achievments
	var result Result
	var user User
	var villages []Village
	var pois []Poi
	db.Table("users").Select("users.id, users.avatar, users.username, AVG(trips.km) as avgkm, SUM(trips.km) as km, COUNT(trips.id) as tripcount, MAX(trips.km) as maxkm").Joins("JOIN user_trip ON users.id = user_trip.user_id").Joins("JOIN trips ON user_trip.trip_id = trips.id").Group("users.id, users.username").First(&result, id)
	db.Table("users").Select("SUM(trips.km) as kmbike").Joins("JOIN user_trip ON users.id = user_trip.user_id").Joins("JOIN trips ON user_trip.trip_id = trips.id").Where("trips.withbike = ?", true).Group("users.id").First(&userpage, id)
	//db.Table("users").Select("users.id, users.username, trips.*").Joins("JOIN user_trip ON users.id = user_trip.user_id").Joins("JOIN trips ON user_trip.trip_id = trips.id").First(&userpage, id)
	db.Preload("Trips", func(db *gorm.DB) *gorm.DB {
		return db.Order("trips.timestamp DESC")
	}).First(&user, id)

	// Villages
	db.Raw("select distinct villages.* from users inner join user_trip on user_trip.user_id = users.id inner join trips on trips.id = user_trip.trip_id inner join trip_village on trip_village.trip_id = trips.id inner join villages on villages.id = trip_village.village_id where users.username = ? order by villages.village", user.Username).Find(&villages)
	// POIs
	db.Raw("select max(pois.elevation) as max_peak, count(*) FILTER (WHERE type = 'ruin') AS ruin_count, count(*) FILTER (WHERE type = 'attraction') as attraction_count, count(*) FILTER (WHERE type = 'station' OR type = 'halt') as station_count, count(*) FILTER (WHERE type = 'viewpoint') as viewpoint_count, count(*) FILTER (WHERE type = 'peak' ) as peak_count, count(*) FILTER (WHERE type = 'place_of_worship') as worship_count from users inner join user_trip on user_trip.user_id = users.id inner join trips on trips.id = user_trip.trip_id inner join trip_poi on trip_poi.trip_id = trips.id inner join pois on pois.id = trip_poi.poi_id where users.username = ?", user.Username).Find(&poiStats)

	userpage.Km = result.Km
	userpage.AvgKm = result.Avgkm
	userpage.Avatar = result.Avatar
	userpage.Kmwalk = userpage.Km - userpage.Kmbike
	userpage.Trips = user.Trips
	userpage.ID = user.ID
	userpage.Username = user.Username
	userpage.VillagesCount = len(villages)
	userpage.Villages = villages
	userpage.Pois = pois
	userpage.TripCount = result.Tripcount
	userpage.Maxkm = result.Maxkm

	userpage.Km100 = userpage.Km
	if userpage.Km > 200 {
		userpage.Km100 = userpage.Km - 200

	} else if userpage.Km > 100 {
		userpage.Km100 = userpage.Km - 100
	}

	achievments.Walker1 = userpage.Km >= 100
	achievments.Walker2 = userpage.Km >= 200
	achievments.Walker3 = userpage.Km >= 300

	achievments.Explorer1 = len(villages) >= 5
	achievments.Explorer2 = len(villages) >= 10
	achievments.Explorer3 = len(villages) >= 15

	achievments.Challenger1 = userpage.Maxkm >= 10
	achievments.Challenger2 = userpage.Maxkm >= 20
	achievments.Challenger3 = userpage.Maxkm >= 30

	achievments.Score = 0
	if achievments.Walker1 && achievments.Walker2 && achievments.Walker3 {
		achievments.Score++
	}

	if achievments.Explorer1 && achievments.Explorer2 && achievments.Explorer3 {
		achievments.Score++
	}

	if achievments.Challenger1 && achievments.Challenger2 && achievments.Challenger3 {
		achievments.Score++
	}

	userpage.Achievments = achievments

	return userpage
}

// GetUserByID --
func GetUserByID(userID int) (*User, error) {
	var user User
	err := db.Table("users").Select("*").Where("users.id = ?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil

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
func IsUserValid(username, password string) (int, bool) {

	if isUsernameAvailable(username) {
		return 0, false
	}

	var user User
	db.Where("username = ?", username).First(&user)

	var pass string = crypting(password + user.Salt)

	return user.ID, pass == user.Password
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
func RegisterNewUser(username, password string) (int, error) {
	var user User
	if strings.TrimSpace(password) == "" {
		return -1, errors.New("The password can't be empty")
	} else if !isUsernameAvailable(username) {
		return -1, errors.New("The username isn't available")
	}

	salt := salting(10)
	pass := crypting(password + salt)

	user = User{Username: username, Password: pass, Salt: salt, Avatar: username}
	result := db.Create(&user) // pass pointer of data to Create
	if result.Error != nil {
		return -1, result.Error
	}

	return user.ID, nil
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

// UpdateSettings --
func UpdateSettings(userID int, avatar string) error {
	err := db.Table("users").Where("id = ?", userID).Update("avatar", avatar)
	if err != nil {
		return err.Error
	}
	return nil
}
