// models.trip.go

package springkilometers

import (
	"strconv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"context"
)

type trip struct {
	gorm.Model
	ID              string `json:"id"`
	Title           string  `json:"title"`
	Content         string  `json:"content"`
	KilometersCount float64 `json:"kmc"`
	Date            string  `json:date`
}



// GetTrips --
// All trips with users sorted by date
// TODO: get trips 
func GetTrips() []trip {
	db, ok := ctx.Value("DB").(*gorm.DB)
	result := database.Find(&trip)
	return result
}

// GetTripByID
// Return trip given id
func GetTripByID(id string) (*trip, error) {
	// TODO(Ez): get trip by id
	// SQL: select * from trips where trip.id = $id
	database := GetDB()
	database.

	return nil, nil
}

// DeleteTripUserIDReference
func DeleteTripUserIDReference(id string) bool {
	// TODO(Ez): delete all trips from trip_user table with specific trip_id
	// SQL : DELETE FROM trip_user WHERE trip_user.trip_id = $id
}

// CreateNew trip with all users
func CreateNewTrip(title, content, kilometersCount string, users) (*trip, error) {
	kmc, err := strconv.ParseFloat(kilometersCount, 64)
	if err != nil {
		return nil, nil
	}
	newTrip := trip{Title: title, Content: content, KilometersCount: kmc}

	// TODO: New database record  with $newTrip

	// TODO: Each user who append in createNewTrip add to trip_user. With values trip_id, user_id.
	return &a, nil
}

// UpdateTrip --
func UpdateTrip(title, content, kilometersCount string, users) (*trip, error) {
	// TODO: Update all values in trips.
	// TODO: Delete all users in trip_user with this specific trip_id
	// TODO: Each user who append in createNewTrip add to trip_user. With values trip_id, user_id.
}

func deleteTrip(title, content, kilometersCount string, users) (*trip, error) {
	// TODO: Delete all users in trip_user with this specific trip_id
}
