package db

import (
	"database/sql"
)


func Connect(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
        return nil, err
	}

	err = db.Ping()
	if err != nil {
        return nil, err
	}

	return db, nil
}
