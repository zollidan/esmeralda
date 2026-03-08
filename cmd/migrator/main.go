package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var databaseURL, migrationPath, migrationTable string

	flag.StringVar(&databaseURL, "database-url", "", "database url")

	flag.StringVar(&migrationPath, "migration-path", "./migrations", "path to migration folder")

	flag.StringVar(&migrationTable, "migration-table", "", "name of migrations table")
	flag.Parse()

	if databaseURL == "" {
		panic(errors.New("database url required"))
	}

	if migrationPath == "" {
		panic(errors.New("migration path required"))
	}

	m, err := migrate.New(
		"",
		"",
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("Nothing to migrate")

			return
		}

		panic(err)
	}
}
