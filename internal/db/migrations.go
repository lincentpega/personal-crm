package db

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/lincentpega/personal-crm/internal/log"
)

const (
	migrationsPath = "file://db/migrations"
	dbName         = "postgres"
)

func ExecMigrations(db *sql.DB, log *log.Logger) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(migrationsPath, dbName, driver)
	if err != nil {
		return err
	}
    defer m.Close()

	if err = m.Up(); err != nil {
        if err == migrate.ErrNoChange {
            log.InfoLog.Println("No migrations to apply")
            return nil
        }
		return err
	}

	return nil
}
