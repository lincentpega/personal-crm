package person

import (
	"time"
)

type Person struct {
	BirthDate    *time.Time
	FirstName    string
	LastName     *string
	SecondName   *string
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
