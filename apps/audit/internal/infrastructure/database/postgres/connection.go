package postgres

import (
	"database/sql"
	"fmt"

	"github.com/didiegovieira/go-financial-core/apps/audit/internal/infrastructure/database/postgres/migrations"
	"github.com/didiegovieira/go-financial-core/apps/audit/internal/settings"
	_ "github.com/lib/pq"
)

func NewConnection() (*sql.DB, error) {
	database := settings.Settings.Database.Postgres
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		database.Host,
		database.User,
		database.Password,
		database.DatabaseName,
		database.Port,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err := migrations.Run(db); err != nil {
		return nil, err
	}

	return db, nil
}
