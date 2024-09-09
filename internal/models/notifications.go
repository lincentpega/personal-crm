package models

import (
	"context"
	"database/sql"
	"time"
)

type Type int

const (
	KeepInTouch Type = iota
)

func (t Type) String() string {
	switch t {
	case KeepInTouch:
		return "KeepInTouch"
	default:
		return "WrongValue"
	}
}

type Notification struct {
	NotificationTime *time.Time
	Description      string
	Type             Type
	PersonID         int
}

type NotificationRepository struct {
	DB *sql.DB
}

func (r *NotificationRepository) Insert(ctx context.Context, n *Notification) error {
	const stmt = `INSERT INTO person_notifications (person_id, type, notification_time, description)
        VALUES($1, $2, $3, $4)`

	_, err := r.DB.ExecContext(ctx, stmt, n.PersonID, n.Type, n.NotificationTime, n.Description)
	if err != nil {
		return err
	}

	return nil
}

func (r *NotificationRepository) Get(ctx context.Context, id int) (*Notification, error) {
	return nil, nil
}
