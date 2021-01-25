package springkilometers

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func ConnectToDB() *gorm.DB {
	dsn := "host=localhost user=postgres password=postgres dbname=jarnikilometry port=9920 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("Error: connection to DB")
		log.Panic(err)
	}
	return db
}
