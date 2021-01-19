// models.trip.go

package main

import (
	"strconv"
)

type trip struct {
	ID              int     `json:"id"`
	Title           string  `json:"title"`
	Content         string  `json:"content"`
	KilometersCount float64 `json:"kmc"`
	Date            string  `json:date`
}

var db = db.GetDB
// All trips with users sorted by date
func getTrips() []trip {

}

// Return trip given id
func getTripByID(id int) (*trip, error) {
	// TODO(Ez): get trip by id
	// SQL: select * from trips where trip.id = $id
	return nil, nil
}

func deleteTripUserIDReference(id int) bool {
	// TODO(Ez): delete all trips from trip_user table with specific trip_id
	// SQL : DELETE FROM trip_user WHERE trip_user.trip_id = $id
}

// Create new trip with all users
func createNewTrip(title, content, kilometersCount string, users) (*trip, error) {
	kmc, err := strconv.ParseFloat(kilometersCount, 64)
	if err != nil {
		return nil, nil
	}
	// _newTrip := trip{Title: title, Content: content, KilometersCount: kmc}

	// TODO: New database record  with $newTrip

	// TODO: Each user who append in createNewTrip add to trip_user. With values trip_id, user_id.
	return &a, nil
}

func updateTrip(title, content, kilometersCount string, users) (*trip, error) {
	// TODO: Update all values in trips.
	// TODO: Delete all users in trip_user with this specific trip_id
	// TODO: Each user who append in createNewTrip add to trip_user. With values trip_id, user_id.
}

func deleteTrip(title, content, kilometersCount string, users) (*trip, error) {
	// TODO: Delete all users in trip_user with this specific trip_id
}
