package springkilometers

import "log"

// Achievment --
type Achievment struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Group       string  `json:"group"`
	Description string  `json:"description"`
	Count       int     `json:"count"`
	MyCount     float64 `json:"mycount"`
	Done        bool    `json:"done"`
}

// Sumkm --
type Sumkm struct {
	Kmsum   float64 `json:"kmsum"`
	Bikesum float64 `json:"bikesum"`
	Walksum float64 `json:"walksum"`
	Maxkm   float64 `json:"maxkm"`
}

// VillageCount --
type VillageCount struct {
	VillageCount int `json:"village_count"`
}

// SundayTrip --
type SundayTrip struct {
	SundayTrip int `json:"sunday_trip"`
}

//GetAchievmentsByUserID --
func GetAchievmentsByUserID(userID int) []Achievment {
	// achievments bond to a trip -> trip_achievment
	// Some must be generated based on user data
	var sumkm Sumkm
	var sundayTrip SundayTrip
	var villageCount VillageCount
	var poiStats PoiStats
	var achmap map[string]float64
	achmap = make(map[string]float64)

	// Village count
	db.Raw("select count(distinct village_id) as village_count from users inner join user_trip on user_trip.user_id = users.id inner join trips on trips.id = user_trip.trip_id inner join trip_village on trip_village.trip_id = trips.id inner join villages on villages.id = trip_village.village_id where users.id = ?", userID).Scan(&villageCount)
	achmap["villagecount"] = float64(villageCount.VillageCount)
	// Poi count
	//db.Raw("select count(distinct poi_id) from users inner join user_trip on user_trip.user_id = users.id inner join trips on trips.id = user_trip.trip_id inner join trip_village on trip_village.trip_id = trips.id inner join trip_poi on trip_poi.trip_id = trips.id inner join pois on pois.id = trip_poi.poi_id where users.id = ?", userID)

	// Walk km count, Bike km count, overall km count
	var err error
	err = db.Raw("select sum(case when trips.withbike then trips.km * 0.25 else trips.km end) AS kmsum, sum(case when trips.withbike then trips.km end) AS bikesum, sum(case when not trips.withbike then trips.km end) AS walksum, max(trips.km) AS maxkm from users left join user_trip on user_trip.user_id = users.id left join trips on trips.id = user_trip.trip_id where users.id = ?", userID).Scan(&sumkm).Error
	if err != nil {
		log.Panic(err)
	}
	achmap["bikekm"] = sumkm.Bikesum
	achmap["sumkm"] = sumkm.Kmsum
	achmap["walksum"] = sumkm.Walksum
	achmap["maxtrip"] = sumkm.Maxkm
	// Sunday tripselect sum(case when trips.withbike then trips.km / 4 else trips.km end) as kmsum, sum(case when trips.withbike then trips.km else 0 end) as bike, sum(case when not trips.withbike then trips.km else 0 end) as walk from users left join user_trip on user_trip.user_id = users.id left join trips on trips.id = user_trip.trip_id where users.id =
	db.Raw("select count( case when 0 = (round(floor((trips.timestamp - 3600)/ 86400) + 4))::INTEGER % 7 then 1 end) as sunday_trip from user_trip left join trips on trips.id = user_trip.trip_id where user_trip.user_id = ?", userID).Scan(&sundayTrip)
	achmap["sunday"] = float64(sundayTrip.SundayTrip)

	// Get pois
	db.Raw("select max(pois.elevation) as max_peak, count(distinct pois.id) FILTER (WHERE type = 'ruin' or historic = 'castle' ) AS ruin_count, count(distinct pois.id) FILTER (WHERE type = 'attraction') as attraction_count, count(distinct pois.id) FILTER (WHERE type = 'station' OR type = 'halt') as station_count, count(distinct pois.id ) FILTER (WHERE type = 'viewpoint') as viewpoint_count, count(distinct pois.id) FILTER (WHERE type = 'place_of_worship') as worship_count from users inner join user_trip on user_trip.user_id = users.id inner join trips on trips.id = user_trip.trip_id inner join trip_poi on trip_poi.trip_id = trips.id inner join pois on pois.id = trip_poi.poi_id where users.id = ?", userID).Find(&poiStats)
	achmap["maxpeak"] = poiStats.MaxPeak
	achmap["station"] = float64(poiStats.StationCount)
	achmap["attraction"] = float64(poiStats.AttractionCount)
	achmap["castle"] = float64(poiStats.RuinCount)
	achmap["worship"] = float64(poiStats.WorshipCount)
	achmap["viewpoint"] = float64(poiStats.ViewpointCount)
	achmap["cityconnection"] = 0

	var achievments []Achievment
	rows, err := db.Table("achievments").Rows()
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var achievment Achievment
		db.ScanRows(rows, &achievment)
		if achmap[achievment.Group] >= float64(achievment.Count) {
			achievment.Done = true
		}
		achievment.MyCount = achmap[achievment.Group]
		achievments = append(achievments, achievment)
	}
	return achievments

}

//GetAchievmentsByUserIDForUserpage --
func GetAchievmentsByUserIDForUserpage(userID int) []Achievment {
	// achievments bond to a trip -> trip_achievment
	// Some must be generated based on user data
	var sumkm Sumkm
	var sundayTrip SundayTrip
	var villageCount VillageCount
	var poiStats PoiStats
	var achmap map[string]float64
	achmap = make(map[string]float64)

	// Village count
	db.Raw("select count(distinct village_id) as village_count from users inner join user_trip on user_trip.user_id = users.id inner join trips on trips.id = user_trip.trip_id inner join trip_village on trip_village.trip_id = trips.id inner join villages on villages.id = trip_village.village_id where users.id = ?", userID).Scan(&villageCount)
	achmap["villagecount"] = float64(villageCount.VillageCount)
	// Poi count
	//db.Raw("select count(distinct poi_id) from users inner join user_trip on user_trip.user_id = users.id inner join trips on trips.id = user_trip.trip_id inner join trip_village on trip_village.trip_id = trips.id inner join trip_poi on trip_poi.trip_id = trips.id inner join pois on pois.id = trip_poi.poi_id where users.id = ?", userID)

	// Walk km count, Bike km count, overall km count
	var err error
	err = db.Raw("select sum(case when trips.withbike then trips.km * 0.25 else trips.km end) AS kmsum, sum(case when trips.withbike then trips.km end) AS bikesum, sum(case when not trips.withbike then trips.km end) AS walksum, max(trips.km) AS maxkm from users left join user_trip on user_trip.user_id = users.id left join trips on trips.id = user_trip.trip_id where users.id = ?", userID).Scan(&sumkm).Error
	if err != nil {
		log.Panic(err)
	}
	achmap["bikekm"] = sumkm.Bikesum
	achmap["sumkm"] = sumkm.Kmsum
	achmap["walksum"] = sumkm.Walksum
	achmap["maxtrip"] = sumkm.Maxkm
	// Sunday tripselect sum(case when trips.withbike then trips.km / 4 else trips.km end) as kmsum, sum(case when trips.withbike then trips.km else 0 end) as bike, sum(case when not trips.withbike then trips.km else 0 end) as walk from users left join user_trip on user_trip.user_id = users.id left join trips on trips.id = user_trip.trip_id where users.id =
	db.Raw("select count( case when 0 = (round(floor((trips.timestamp - 3600)/ 86400) + 4))::INTEGER % 7 then 1 end) as sunday_trip from user_trip left join trips on trips.id = user_trip.trip_id where user_trip.user_id = ?", userID).Scan(&sundayTrip)
	achmap["sunday"] = float64(sundayTrip.SundayTrip)

	// Get pois
	db.Raw("select max(pois.elevation) as max_peak, count(distinct pois.id) FILTER (WHERE type = 'ruin' or historic = 'castle' ) AS ruin_count, count(distinct pois.id) FILTER (WHERE type = 'attraction') as attraction_count, count(distinct pois.id) FILTER (WHERE type = 'station' OR type = 'halt') as station_count, count(distinct pois.id ) FILTER (WHERE type = 'viewpoint') as viewpoint_count, count(distinct pois.id) FILTER (WHERE type = 'place_of_worship') as worship_count from users inner join user_trip on user_trip.user_id = users.id inner join trips on trips.id = user_trip.trip_id inner join trip_poi on trip_poi.trip_id = trips.id inner join pois on pois.id = trip_poi.poi_id where users.id = ?", userID).Find(&poiStats)
	achmap["maxpeak"] = poiStats.MaxPeak
	achmap["station"] = float64(poiStats.StationCount)
	achmap["attraction"] = float64(poiStats.AttractionCount)
	achmap["castle"] = float64(poiStats.RuinCount)
	achmap["worship"] = float64(poiStats.WorshipCount)
	achmap["viewpoint"] = float64(poiStats.ViewpointCount)
	achmap["cityconnection"] = 0

	var achievments []Achievment
	rows, err := db.Table("achievments").Rows()
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var achievment Achievment
		db.ScanRows(rows, &achievment)
		if achmap[achievment.Group] >= float64(achievment.Count) {
			achievments = append(achievments, achievment)
		}
	}
	return achievments

}
