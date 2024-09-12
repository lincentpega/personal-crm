package notifications

import "time"

type Type string

const (
	KeepInTouch Type = "keep_in_touch"
)

type Status string

const (
	Pending Status = "pending"
	Raised  Status = "raised"
	Failed  Status = "failed"
)

type Notification struct {
	NotificationTime *time.Time
	Description      string
	Type             Type
	Status           Status
	PersonID         int
	ID               int
}
