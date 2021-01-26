package springkilometers

import (
	"log"
	"strconv"

	"gorm.io/gorm"
)

// Trip model
type Trip struct {
	Model
	TripID  int     `json:"trip_id"`
	Title   string  `json:"title"`
	Content string  `json:"content"`
	Km      float64 `json:"km"`
	//Date            string  `json:date`
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
	err := db.Select("trip_id").Where("trip_id = ? AND deleted_on = ? ", id, 0).First(&trip).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if trip.TripID >= 0 {
		return &trip, nil
	}
	return nil, nil
}

// ExistTripByID --
func ExistTripByID(id int) (bool, error) {
	var trip Trip
	err := db.Select("trip_id").Where("trip_id = ? AND deleted_on = ? ", id, 0).First(&trip).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if trip.TripID > 0 {
		return true, nil
	}
	return false, nil
}

// DeleteTripByID --
func DeleteTripByID(id int) (bool, error) {
	db.Delete(&Trip{}, id)
	return true, nil
}

// CreateNewTrip trip with all users
func CreateNewTrip(title, content, kilometersCount string) (*Trip, error) {
	kmc, err := strconv.ParseFloat(kilometersCount, 64)
	if err != nil {
		return nil, nil
	}
	newTrip := Trip{Title: title, Content: content, Km: kmc}

	// TODO: New database record  with $newTrip
	result := db.Create(&newTrip) // pass pointer of data to Create
	if result.Error != nil {
		log.Panic(result.Error)
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
