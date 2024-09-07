package models

import "database/sql"

type Contact struct{}

type ContactModel struct {
	DB *sql.DB
}

func (m *ContactModel) Insert() {
}
