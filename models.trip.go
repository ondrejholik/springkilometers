// models.trip.go

package main

import "errors"

type trip struct {
	ID              int     `json:"id"`
	Title           string  `json:"title"`
	Content         string  `json:"content"`
	KilometersCount float64 `json:"kmc"`
}

// For this demo, we're storing the article list in memory
// In a real application, this list will most likely be fetched
// from a database or from static files
var tripList = []trip{
	trip{ID: 0, Title: "Vylet na Hardegg", Content: "Jednoho krasneho dne jsme se vypravili za hranice svych moznosti...", KilometersCount: 20.0},
	trip{ID: 1, Title: "Vylet na kole do Znojma", Content: "Na kole az do Znojma", KilometersCount: 54.2},
}

// Return a list of all the articles
func getAllTrips() []trip {
	return tripList
}

// Fetch an article based on the ID supplied
func getTripByID(id int) (*trip, error) {
	for _, a := range tripList {
		if a.ID == id {
			return &a, nil
		}
	}
	return nil, errors.New("Trip not found")
}

// Create a new article with the title and content provided
func createNewTrip(title, content string, kilometersCount float64) (*trip, error) {
	// Set the ID of a new article to one more than the number of articles
	a := trip{ID: len(tripList) + 1, Title: title, Content: content, KilometersCount: kilometersCount}

	// Add the article to the list of articles
	tripList = append(tripList, a)

	return &a, nil
}
