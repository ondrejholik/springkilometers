package springkilometers

import (
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	"gorm.io/gorm"
)

// Trip model
type Trip struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Withbike bool    `json:"withbike"`
	Content  string  `json:"content"`
	Km       float64 `json:"km"`
	Mapycz   string  `json:"mapycz"`

	Tiny   string `json:"tiny"`
	Small  string `json:"small"`
	Medium string `json:"medium"`
	Large  string `json:"large"`
	Gpx    string `json:"gpx"`

	Timestamp int64 `json:"timestamp"`
	Year      int   `json:"year"`
	Month     int   `json:"month"`
	Day       int   `json:"day"`
	Hour      int   `json:"hour"`
	Minute    int   `json:"minute"`

	AuthorID int       `json:"author_id"`
	Users    []User    `gorm:"many2many:user_trip;"`
	Villages []Village `gorm:"many2many:trip_village;"`
	Pois     []Poi     `gorm:"many2many:trip_poi;"`
	Comments []CommentResult

	CreatedOn  time.Time `json:"created_on"`
	DeletedOn  time.Time `json:"deleted_on"`
	ModifiedOn time.Time `json:"modified_on"`
	UpdatedOn  time.Time `json:"updated_on"`
}

// TripAll for view showTrips
type TripAll struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Withbike bool    `json:"withbike"`
	Content  string  `json:"content"`
	Km       float64 `json:"km"`
	Avatar   string  `json:"avatar"`
	Username string  `json:"username"`

	Medium        string `json:"medium"`
	Small         string `json:"small"`
	CommentsCount int    `json:"comments_count"`
	HeartCount    int    `json:heart_count"`

	Timestamp int64 `json:"timestamp"`
	Year      int   `json:"year"`
	Month     int   `json:"month"`
	Day       int   `json:"day"`
	Hour      int   `json:"hour"`
	Minute    int   `json:"minute"`

	AuthorID int `json:"author_id"`
}

// GetUserTrips --
func GetUserTrips(userID int) []Trip {
	var trips []Trip
	db.Table("trips").Where("trips.author_id = ?", userID).Order("timestamp desc").Scan(&trips)
	return trips
}

// GetTrips --
// All trips with users sorted by date
func GetTrips() []TripAll {
	var trips []TripAll
	result := db.Table("trips").Joins("left join users on users.id = author_id").Joins("left join comments on comments.trip_id = trips.id").Group("trips.id, users.id").Order("timestamp desc").Select("users.ID as author_id, trips.*, users.Username, users.Avatar, COUNT(comments.id) as comments_count, count( case when comments.message = '❤️' then 1 end) as heart_count").Scan(&trips)
	if result.Error != nil {
		log.Panic(result.Error)
	}

	return trips
}

// GetTripByID --
// Return trip given id
func GetTripByID(id int) (*Trip, error) {
	var trip Trip
	err := db.Select("*").Where("id = ? ", id).First(&trip).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if trip.ID >= 0 {
		return &trip, nil
	}
	return nil, nil
}

// GetTripByIDWithUsers --
// Return trip given id
func GetTripByIDWithUsers(id int) (*Trip, error) {
	var trip Trip
	err := db.Preload("Users").Preload("Villages").Preload("Pois").First(&trip, id).Error
	trip.Comments = GetComments(id)
	if err != nil {
		return nil, err
	}

	if trip.ID >= 0 {
		return &trip, nil
	}
	return nil, nil
}

// TripHasUser --
func TripHasUser(id int, username string) bool {
	// select * from trips inner join user_trip on user_trip.trip_id = trips.id
	// inner join users on users.id = user_trip.user_id
	// where users.username = 'o' and trips.id = 349;
	var count int64
	db.Table("trips").Joins("INNER JOIN user_trip ON user_trip.trip_id = trips.id").Joins("INNER JOIN users on users.id = user_trip.user_id").Where("users.username = ? and trips.id = ?", username, id).Count(&count)
	return count == 1
}

// ExistTripByID --
func ExistTripByID(id int) (bool, error) {
	var trip Trip
	err := db.Select("id").Where("id = ?", id).First(&trip).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if trip.ID >= 0 {
		return true, nil
	}
	return false, nil
}

// TripJoinsUser --
func TripJoinsUser(username string, trip Trip) {
	var user User
	db.Where("username = ?", username).First(&user)
	db.Model(&trip).Association("Users").Append(&user)
}

// TripDisjoinsUser --
func TripDisjoinsUser(username string, trip Trip) {
	var user User
	db.Where("username = ?", username).First(&user)
	db.Model(&trip).Association("Users").Delete(&user)
}

func tripBelongsToUsername(tripID int, userID int) bool {
	var trip Trip
	db.Where("id = ?", tripID).First(&trip)
	return trip.AuthorID == userID
}

func getDate() (int, int, int, int, int) {
	currentTime := time.Now()
	minute := currentTime.Minute()
	hour := currentTime.Hour()
	day := currentTime.Day()
	month := currentTime.Month()
	year := currentTime.Year()

	return year, int(month), day, hour, minute
}

// CreateNewTrip trip with all users
func CreateNewTrip(user User, name, content, kilometersCount, withbike, gpxname, mapycz string) (*Trip, error) {
	kmc, err := strconv.ParseFloat(kilometersCount, 64)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	year, month, day, hour, minute := getDate()

	newTrip := Trip{
		Name:      name,
		Content:   content,
		Km:        kmc,
		Withbike:  withbike == "on",
		AuthorID:  user.ID,
		Tiny:      "/static/default/tiny.webp",
		Small:     "/static/default/small.webp",
		Medium:    "/static/default/medium.webp",
		Large:     "/static/default/large.webp",
		Timestamp: time.Now().Unix(),
		Year:      year,
		Month:     month,
		Day:       day,
		Hour:      hour,
		Minute:    minute,
		Gpx:       gpxname,
		Mapycz:    mapycz,
	}

	result := db.Create(&newTrip) // pass pointer of data to Create
	if result.Error != nil {
		log.Println(result.Error)
		return nil, result.Error
	}

	// User, who created trip also "join" trip
	TripJoinsUser(user.Username, newTrip)

	return &newTrip, nil
}

// UpdateTrip --
func UpdateTrip(id, userID int, name, content, kilometersCount, withbike string) (*Trip, error) {
	kmc, err := strconv.ParseFloat(kilometersCount, 64)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	wb := withbike == "on"

	trip, err := GetTripByID(id)
	if err != nil {
		return nil, err
	}
	trip.Name = name
	trip.Content = content
	trip.Km = kmc
	trip.Withbike = wb

	if !tripBelongsToUsername(id, userID) {
		err := errors.New("Username does not belong to trip")
		return nil, err
	}

	db.Save(&trip)
	return trip, nil
}

// UpdateTripStruct --
func UpdateTripStruct(trip Trip) (*Trip, error) {

	db.Save(&trip)
	return &trip, nil
}

// DeleteTripByID --
func DeleteTripByID(id int) (bool, error) {

	var trip Trip
	db.Table("trips").First(&trip, id)
	// Delete also imgs from path
	err := os.Remove("." + trip.Tiny)
	if err != nil {
		log.Println(err)
	}
	err = os.Remove("." + trip.Small)
	if err != nil {
		log.Println(err)
	}
	err = os.Remove("." + trip.Medium)
	if err != nil {
		log.Println(err)
	}
	err = os.Remove("." + trip.Large)
	if err != nil {
		log.Println(err)
	}

	err = os.Remove("." + trip.Gpx)
	if err != nil {
		log.Println(err)
	}

	db.Delete(&Trip{}, id)
	return true, nil
}
