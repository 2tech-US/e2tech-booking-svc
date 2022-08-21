package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

const dbDriver = "postgres"

var db *sql.DB

func Connect(dbUrl string) *Queries {
	var err error
	db, err = sql.Open(dbDriver, dbUrl)
	if err != nil {
		log.Fatal("Failed to connect to db:", err)
	}

	return New(db)
}
