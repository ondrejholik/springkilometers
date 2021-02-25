package springkilometers

import (
	"io/ioutil"
	"log"
	"math"
	"strconv"
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

// Poi -- Points of interest
type Poi struct {
	ID        int     `json:"id"`
	Type      string  `json:"type"`
	Name      string  `json:"name"`
	Historic  string  `json:"historic"`
	Elevation float64 `json:"elevation"`
	Lat       float64 `json:"lat"`
	Lon       float64 `json:"lon"`
}

// PoiStats --
type PoiStats struct {
	AttractionCount int     `json:"attraction_count"`
	PeakCount       int     `json:"peak_count"`
	RuinCount       int     `json:"ruin_count"`
	StationCount    int     `json:"station_count"`
	ViewpointCount  int     `json:"viewpoint_count"`
	WorshipCount    int     `json:"worship_count"`
	MaxPeak         float64 `json:"max_peak"`
}

// TripVillage --
type TripVillage struct {
	TripID    int `json:"trip_id"`
	VillageID int `json:"village_id"`
}

// TripPoi --
type TripPoi struct {
	TripID int `json:"trip_id"`
	PoiID  int `json:"poi_id"`
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
func LoadVillages(bounds gpx.GpxBounds) []Village {
	var villages []Village
	db.Table("villages").Where("lat BETWEEN ? AND ?", bounds.MinLatitude, bounds.MaxLatitude).Where("lon BETWEEN ? AND ?", bounds.MinLongitude, bounds.MaxLongitude).Scan(&villages)

	return villages
}

// LoadPoi --
func LoadPoi(bounds gpx.GpxBounds) []Poi {
	var poi []Poi
	db.Table("pois").Where("lat BETWEEN ? AND ?", bounds.MinLatitude, bounds.MaxLatitude).Where("lon BETWEEN ? AND ?", bounds.MinLongitude, bounds.MaxLongitude).Scan(&poi)

	return poi
}

func findGpxPoi(gpx *[]Gps, poi Poi, distanceFrom float64, c chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, gps := range *gpx {
		if harvestineDistancePoi(gps, poi) < distanceFrom {
			c <- poi.ID
			break
		}
	}
}

func findGpxVill(gpx *[]Gps, village Village, distanceFrom float64, c chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, gps := range *gpx {
		if harvestineDistanceVil(gps, village) < distanceFrom {
			c <- village.ID
			break
		}
	}
}

func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

func harvestineDistancePoi(gps Gps, poi Poi) float64 {
	var la1, lo1, la2, lo2, r float64
	la1 = gps.Lat * math.Pi / 180
	lo1 = gps.Lon * math.Pi / 180
	la2 = poi.Lat * math.Pi / 180
	lo2 = poi.Lon * math.Pi / 180

	r = 6378100 // Earth radius in METERS

	// calculate
	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)

	return 2 * r * math.Asin(math.Sqrt(h))
}
func harvestineDistanceVil(gps Gps, village Village) float64 {
	var la1, lo1, la2, lo2, r float64
	la1 = gps.Lat * math.Pi / 180
	lo1 = gps.Lon * math.Pi / 180
	la2 = village.Lat * math.Pi / 180
	lo2 = village.Lon * math.Pi / 180

	r = 6378100 // Earth radius in METERS

	// calculate
	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)

	return 2 * r * math.Asin(math.Sqrt(h))
}

func monitorWorker(wg *sync.WaitGroup, cs chan int) {
	wg.Wait()
	close(cs)
}

// AddGpsToTrip --
func AddGpsToTrip(gpxpath string, tripID int) {
	villchan := make(chan int)
	poichan := make(chan int)
	wgpoi := &sync.WaitGroup{}
	wgvill := &sync.WaitGroup{}
	gpx, bounds := LoadGpx(gpxpath)

	var pois = LoadPoi(bounds)
	var villages = LoadVillages(bounds)

	for _, poi := range pois {
		wgpoi.Add(1)
		go findGpxPoi(&gpx, poi, 64, poichan, wgpoi)
	}

	for _, village := range villages {
		wgvill.Add(1)
		go findGpxVill(&gpx, village, 1000, villchan, wgvill)
	}

	go monitorWorker(wgpoi, poichan)
	go monitorWorker(wgvill, villchan)

	var tripVillage []TripVillage
	var tripPoi []TripPoi

	for poi := range poichan {
		tripPoi = append(tripPoi, TripPoi{TripID: tripID, PoiID: poi})
	}

	for vil := range villchan {
		tripVillage = append(tripVillage, TripVillage{TripID: tripID, VillageID: vil})
	}

	db.Table("trip_poi").Create(&tripPoi)
	db.Table("trip_village").Create(&tripVillage)

	MyCache.Delete(Ctx, "trip:"+strconv.Itoa(tripID))

}

// RemoveGpsFromTrip --
func RemoveGpsFromTrip(tripID int) {
	db.Table("trip_poi").Delete(&TripPoi{}, &TripPoi{TripID: tripID})
	db.Table("trip_village").Delete(&TripVillage{}, &TripVillage{TripID: tripID})
}
