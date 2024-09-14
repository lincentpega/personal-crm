package db

import (
	"database/sql"
	"path/filepath"
	"runtime"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/lincentpega/personal-crm/internal/log"
)

const (
	dbName = "postgres"
)

func ExecMigrations(db *sql.DB, log *log.Logger) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(getMigrationsSourceURL(), dbName, driver)
	if err != nil {
		return err
	}

	if err = m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.InfoLog.Println("No migrations to apply")
			return nil
		}
		return err
	} else {
		log.InfoLog.Println("Migrations applied successfully")
	}

	return nil
}

func getMigrationsSourceURL() string {
	_, filename, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(filename), "..", "..")
	migrationsPath := filepath.Join(projectRoot, "db", "migrations")

	return "file://" + migrationsPath
}
