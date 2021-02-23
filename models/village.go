package springkilometers

import (
	"io/ioutil"
	"log"
	"math"
	"sync"

	"github.com/tkrajina/gpxgo/gpx"
)

// Gps --
type Gps struct {
	Lat float64
	Lon float64
	El  float64
}

// Village --
type Village struct {
	ID      int     `json:"id"`
	Village string  `json:"village"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	Type    string  `json:"type"`
	Trips   []Trip  `gorm:"many2many:trip_village;"`
}

// TripVillage --
type TripVillage struct {
	TripID    int `json:"trip_id"`
	VillageID int `json:"village_id"`
}

// LoadGpx --
func LoadGpx(filepath string) ([]Gps, gpx.GpxBounds) {
	gps := []Gps{}

	// Gpx analyzing
	gpxBytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatal("File not exists")
	}
	gpxFile, err := gpx.ParseBytes(gpxBytes)
	if err != nil {
		log.Println(err)
	}

	bounds := gpxFile.Bounds()

	var latbound float64 = 0.009
	var lonmaxbound float64 = 1 / (math.Cos(bounds.MaxLatitude) * 111.32)
	var lonminbound float64 = 1 / (math.Cos(bounds.MinLatitude) * 111.32)
	bounds.MaxLatitude += latbound
	bounds.MaxLongitude += lonmaxbound
	bounds.MinLatitude -= latbound
	bounds.MinLongitude -= lonminbound

	// Analyize/manipulate your track data here...
	for _, track := range gpxFile.Tracks {
		for _, segment := range track.Segments {
			for _, point := range segment.Points {
				gps = append(gps, Gps{Lat: point.Point.Latitude, Lon: point.Point.Longitude, El: point.Point.Elevation.Value()})
			}
		}
	}
	return gps, bounds
}

// LoadVillages --
// get data from database with bounds
// first from file
func LoadVillages(bounds gpx.GpxBounds) []Village {
	var villages []Village
	db.Table("villages").Where("lat BETWEEN ? AND ?", bounds.MinLatitude, bounds.MaxLatitude).Where("lon BETWEEN ? AND ?", bounds.MinLongitude, bounds.MaxLongitude).Scan(&villages)
	return villages
}

func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

func harvestineDistance(gps Gps, vil Village) float64 {
	var la1, lo1, la2, lo2, r float64
	la1 = gps.Lat * math.Pi / 180
	lo1 = gps.Lon * math.Pi / 180
	la2 = vil.Lat * math.Pi / 180
	lo2 = vil.Lon * math.Pi / 180

	r = 6378100 // Earth radius in METERS

	// calculate
	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)

	return 2 * r * math.Asin(math.Sqrt(h))
}

func monitorWorker(wg *sync.WaitGroup, cs chan int) {
	wg.Wait()
	close(cs)
}

func findGpx(village Village, gpx []Gps, c chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, gps := range gpx {
		if harvestineDistance(gps, village) < 1000 {
			c <- village.ID
			break
		}
	}
}

// AddVillagesToTrip --
func AddVillagesToTrip(gpxpath string, tripID int) {
	res := make(chan int)
	wg := &sync.WaitGroup{}
	gpx, bounds := LoadGpx(gpxpath)
	villages := LoadVillages(bounds)

	for _, village := range villages {
		wg.Add(1)
		go findGpx(village, gpx, res, wg)
	}

	go monitorWorker(wg, res)

	var tripVillage []TripVillage

	for i := range res {
		tripVillage = append(tripVillage, TripVillage{TripID: tripID, VillageID: i})
	}

	db.Table("trip_village").Create(&tripVillage)
}
