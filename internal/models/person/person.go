package person

import (
	"database/sql"
)

type Person struct {
	BirthDate    sql.NullTime
	FirstName    string
	LastName     sql.NullString
	SecondName   sql.NullString
	ContactInfos []ContactInfo
	JobInfos     []JobInfo
	Settings     Settings
	ID           int
}

type ContactInfo struct {
	Method string
	Data   string
}

type JobInfo struct {
	Company  string
	Position string
	Current  bool
}

type Settings struct {
	BirthdayNotify bool
}
