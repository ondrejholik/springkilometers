package springkilometers

import (
	"log"
	"strconv"
	"time"

	"gorm.io/gorm"
)

// Trip model
type Trip struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Content    string    `json:"content"`
	WithBike   bool      `json:"with_bike"`
	Km         float64   `json:"km"`
	CreatedOn  time.Time `json:"created_on"`
	DeletedOn  time.Time `json:"deleted_on"`
	ModifiedOn time.Time `json:"modified_on"`
	UpdatedOn  time.Time `json:"updated_on"`
	Users      []User    `gorm:"many2many:user_trip;"`
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

	if trip.ID >= 0 {
		return &trip, nil
	}
	return nil, nil
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

// DeleteTripByID --
func DeleteTripByID(id int) (bool, error) {
	db.Delete(&Trip{}, id)
	return true, nil
}

// CreateNewTrip trip with all users
func CreateNewTrip(name, content, kilometersCount, withbike string) (*Trip, error) {
	kmc, err := strconv.ParseFloat(kilometersCount, 64)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	wb := withbike == "withbike"

	newTrip := Trip{Name: name, Content: content, Km: kmc, WithBike: wb}

	// TODO: New database record  with $newTrip
	result := db.Create(&newTrip) // pass pointer of data to Create
	if result.Error != nil {
		log.Println(result.Error)
		return nil, result.Error
	}
	return &newTrip, nil

	// TODO: Each user who append in createNewTrip add to trip_user. With values trip_id, user_id.
}

// UpdateTrip --
func UpdateTrip(title, content, kilometersCount string) bool {
	// TODO: Update all values in trips.
	// TODO: Delete all users in trip_user with this specific trip_id
	// TODO: Each user who append in createNewTrip add to trip_user. With values trip_id, user_id.
	return false
}
