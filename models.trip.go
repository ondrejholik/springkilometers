// models.trip.go

package main

import (
	"errors"
	"strconv"
)

type trip struct {
	ID              int     `json:"id"`
	Title           string  `json:"title"`
	Content         string  `json:"content"`
	KilometersCount float64 `json:"kmc"`
}

// All trips with users sorted by date
func getAllTrips() []trip {
	return tripList
}

// Return trip given id
func getTripByID(id int) (*trip, error) {
	for _, a := range tripList {
		if a.ID == id {
			return &a, nil
		}
	}
	return nil, errors.New("Trip not found")
}

// Create new trip with all users
func createNewTrip(title, content, kilometersCount string) (*trip, error) {
	kmc, err := strconv.ParseFloat(kilometersCount, 64)
	if err != nil {
		return nil, nil
	}
	a := trip{ID: len(tripList) + 1, Title: title, Content: content, KilometersCount: kmc}
	tripList = append(tripList, a)
	return &a, nil
}
