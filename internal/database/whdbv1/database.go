package whdbv1

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

type Database struct {
	db *sql.DB
}

func New(file string) *Database {
	// db, err := sql.Open("sqlite", "./"+file)
	db, err := sql.Open("sqlite", file)
	if err != nil {
		log.Fatal(err)
	}

	return &Database{
		db: db,
	}
}

func (d *Database) Close() error {
	return d.db.Close()
}
