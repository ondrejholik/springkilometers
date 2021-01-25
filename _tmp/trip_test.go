// models.article_test.go

package springkilometers

import "testing"

// Test the function that fetches all articles
func TestGetAllTrips(t *testing.T) {
	alist := GetAllTrips()

	// Check that the length of the list of articles returned is the
	// same as the length of the global variable holding the list
	if len(alist) != len(tripList) {
		t.Fail()
	}

	// Check that each member is identical
	for i, v := range alist {
		if v.Content != tripList[i].Content ||
			v.ID != tripList[i].ID ||
			v.Title != tripList[i].Title ||
			v.KilometersCount != tripList[i].KilometersCount {

			t.Fail()
			break
		}
	}
}

// Test the function that fetche an Article by its ID
func TestGetTripByID(t *testing.T) {
	a, err := getTripByID(0)

	if err != nil || a.ID != 0 || a.Title != "Lorem" || a.Content != "Ipsum" || a.KilometersCount != 1.5 {
		t.Fail()
	}
}

// Test the functiogetAllTrips creates a new article
func TestCreateNewTrip(t *testing.T) {
	// get the original count of articles
	originalLength := len(GetAllTrips())

	// add another article
	a, err := createNewTrip("New test title", "New test content", "12.5")

	// get the new count of articles
	allTrips := getAllTrips()
	newLength := len(allTrips)

	if err != nil || newLength != originalLength+1 ||
		a.Title != "New test title" || a.Content != "New test content" || a.KilometersCount != 12.5 {

		t.Fail()
	}
}
