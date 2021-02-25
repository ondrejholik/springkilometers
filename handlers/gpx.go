package springkilometers

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"

	models "github.com/ondrejholik/springkilometers/models"
)

// GpxHandling --
func GpxHandling(gpx *multipart.FileHeader, tripID int) (string, bool) {
	// Load gpx
	gpxname := GetFileName(gpx.Filename)
	filename := fmt.Sprintf("/static/gpx/%s.gpx", gpxname)
	log.Println(filename)

	gpxfile, err := gpx.Open()
	if err != nil {
		log.Panic(err)
		return "", false
	}

	// Save gpx
	err = saveGpx(tripID, gpxname, gpxfile)
	if err != nil {
		return "", false
	}
	log.Println(" ------- gpx success!")

	// Delete old gpx file by path
	trip, err := models.GetTripByID(tripID)
	if err != nil {
		log.Println("trip wasnt found!")
	}
	go deleteGpx(trip.Gpx)

	go func() {
		models.RemoveGpsFromTrip(tripID)
		models.AddGpsToTrip("."+filename, tripID)
	}()

	return filename, true
}

func saveGpx(tripID int, gpxname string, gpxfile multipart.File) error {

	filename := fmt.Sprintf("./static/gpx/%s.gpx", gpxname)
	gpxsaved, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer gpxsaved.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, gpxfile); err != nil {
		return err
	}

	if err = ioutil.WriteFile(filename, buf.Bytes(), 0666); err != nil {
		return err
	}

	go models.AddGpsToTrip(filename, tripID)
	return nil
}

func deleteGpx(gpxpath string) {
	if _, err := os.Stat(gpxpath); !os.IsNotExist(err) {
		os.Remove(gpxpath)
	}
}
