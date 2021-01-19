package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
)

var db *DB

// DB -------------------------
// Connection state to database
type DB struct {
	*pgx.DB
}

// GetDB connect to database
func GetDB() *DB {
	if db == nil {
		// Database initialize //
		config, err := pgx.ParseConfig("postgres://test:helloworld@localhost:26257/jarnikilometry?sslmode=require")
		if err != nil {
			log.Fatal("error configuring the database: ", err)
		}

		config.TLSConfig.ServerName = "localhost"

		// connect to jarnikilometry database
		conn, err := pgx.ConnectConfig(context.Background(), config)
		if err != nil {
			log.Fatal("error connecting to the database: ", err)
		}
		defer conn.Close(context.Background())
	}
	return db
}
