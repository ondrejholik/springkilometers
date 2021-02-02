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
	Author   string  `json:"author"`

	Tiny   string `json:"tiny"`
	Small  string `json:"small"`
	Medium string `json:"medium"`
	Large  string `json:"large"`

	Users []User `gorm:"many2many:user_trip;"`

	CreatedOn  time.Time `json:"created_on"`
	DeletedOn  time.Time `json:"deleted_on"`
	ModifiedOn time.Time `json:"modified_on"`
	UpdatedOn  time.Time `json:"updated_on"`
}

// GetUserTrips --
func GetUserTrips(username string) []Trip {
	var trips []Trip
	db.Table("trips").Where("trips.author = ?", username).Scan(&trips)
	return trips
}

// GetTrips --
// All trips with users sorted by date
func GetTrips() []Trip {
	var trips []Trip
	result := db.Find(&trips)

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

	log.Println(trip.Name)
	log.Println(trip.ID)
	if trip.ID >= 0 {
		return &trip, nil
	}
	return nil, nil
}

// GetTripByIDWithUsers --
// Return trip given id
func GetTripByIDWithUsers(id int) (*Trip, error) {
	var trip Trip
	err := db.Preload("Users").First(&trip, id).Error
	if err != nil && err != gorm.ErrRecordNotFound {
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
	log.Println("Count of user-trip : ", count == 1)
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

func tripBelongsToUsername(tripID int, username string) bool {
	var trip Trip
	db.Where("id = ?", tripID).First(&trip)
	return trip.Author == username
}

// CreateNewTrip trip with all users
func CreateNewTrip(username, name, content, kilometersCount, withbike string) (*Trip, error) {
	kmc, err := strconv.ParseFloat(kilometersCount, 64)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	wb := withbike == "on"

	newTrip := Trip{
		Name:     name,
		Content:  content,
		Km:       kmc,
		Withbike: wb,
		Author:   username,
		Tiny:     "/static/default/tiny.webp",
		Small:    "/static/default/small.webp",
		Medium:   "/static/default/medium.webp",
		Large:    "/static/default/large.webp",
	}

	result := db.Create(&newTrip) // pass pointer of data to Create
	if result.Error != nil {
		log.Println(result.Error)
		return nil, result.Error
	}

	// User, who created trip also "join" trip
	TripJoinsUser(username, newTrip)

	return &newTrip, nil
}

// UpdateTrip --
func UpdateTrip(id int, username, name, content, kilometersCount, withbike string) (*Trip, error) {
	kmc, err := strconv.ParseFloat(kilometersCount, 64)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	wb := withbike == "on"
	log.Println("withbike?", wb)

	trip, err := GetTripByID(id)
	if err != nil {
		return nil, err
	}
	trip.Name = name
	trip.Content = content
	trip.Km = kmc
	trip.Withbike = wb

	if !tripBelongsToUsername(id, username) {
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
		log.Fatal(err)
	}
	err = os.Remove("." + trip.Small)
	if err != nil {
		log.Fatal(err)
	}
	err = os.Remove("." + trip.Medium)
	if err != nil {
		log.Fatal(err)
	}
	err = os.Remove("." + trip.Large)
	if err != nil {
		log.Fatal(err)
	}

	db.Delete(&Trip{}, id)
	return true, nil
}
