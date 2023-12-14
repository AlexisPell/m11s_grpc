package main

import (
	"errors"
	"flag"
	"fmt"
	// migrations
	"github.com/golang-migrate/migrate/v4"
	// sqlite3 driver
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	// driver to get migrations from files
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var storagePath, migrationsPath, migrationsTable string

	flag.StringVar(&storagePath, "storage-path", "", "path to storage")
	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.StringVar(&migrationsTable, "migrations-table", "", "name of migrations table")
	flag.Parse()

	if storagePath == "" {
		panic("storage path is required")
	}
	if migrationsPath == "" {
		panic("migrations path is required")
	}

	m, err := migrate.New(
		"file://"+migrationsPath,
		fmt.Sprintf("sqlite3://%s?x-migrations-table=%s", storagePath, migrationsTable),
	)
	if err != nil {
		panic(fmt.Errorf("migration lib failed: %w", err))
	}

	if err = m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("No migrations to apply")
			return
		} else {
			fmt.Println("Failed to migrate, path: ", migrationsPath, err)
			panic(err)
		}
	}

	fmt.Println("Migrations applied successfully")
}
