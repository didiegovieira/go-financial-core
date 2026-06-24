package postgres

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func NewConnection() (*sql.DB, error) {
	db, err := sql.Open("postgres", "user=youruser dbname=yourdb sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	return db, nil
}
