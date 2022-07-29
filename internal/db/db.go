package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

const dbDriver = "postgres"

func Connect(dbUrl string) *Queries {
	db, err := sql.Open(dbDriver, dbUrl)
	if err != nil {
		log.Fatal("Failed to connect to db:", err)
	}

	return New(db)
}
