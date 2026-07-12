package migrations

import (
	"database/sql"
	"embed"
	"fmt"
)

//go:embed *.sql
var files embed.FS

func Run(db *sql.DB) error {
	entries, err := files.ReadDir(".")
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		statement, err := files.ReadFile(entry.Name())
		if err != nil {
			return err
		}

		if _, err := db.Exec(string(statement)); err != nil {
			return fmt.Errorf("run migration %s: %w", entry.Name(), err)
		}
	}

	return nil
}
