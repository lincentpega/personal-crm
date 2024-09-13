package notifications

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lincentpega/personal-crm/internal/common/txcontext"
	"github.com/lincentpega/personal-crm/internal/db"
	"github.com/lincentpega/personal-crm/internal/models"
)

type NotificationRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) getDB(ctx context.Context) db.DB {
	if tx, ok := txcontext.GetTx(ctx); ok {
		return tx
	}
	return r.db
}

func (r *NotificationRepository) Insert(ctx context.Context, n *Notification) error {
	const stmt = `INSERT INTO notifications (person_id, type, status, notification_time, description)
	VALUES($1, $2, $3, $4, $5)
	RETURNING id`

	err := r.getDB(ctx).QueryRowContext(ctx, stmt, n.PersonID, n.Type, n.Status, n.NotificationTime.UTC(), n.Description).Scan(&n.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *NotificationRepository) UpdateNotificationStatus(ctx context.Context, notifID int, status Status) error {
	const stmt = `UPDATE notifications SET status = $1 WHERE id = $2`

	_, err := r.getDB(ctx).ExecContext(ctx, stmt, status, notifID)
	if err != nil {
		return err
	}

	return nil
}

func (r *NotificationRepository) Get(ctx context.Context, id int) (*Notification, error) {
	const stmt = `SELECT id, person_id, type, status, notification_time, description
	FROM notifications
	WHERE person_id = $1`

	var n Notification

	err := r.getDB(ctx).QueryRowContext(ctx, stmt, id).Scan(&n.ID, &n.PersonID, &n.Type, &n.Status, &n.NotificationTime, &n.Description)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrRecordNotFound
		}
		return nil, err
	}

	return &n, nil
}

func (r *NotificationRepository) GetAwaitingSend(ctx context.Context) ([]Notification, error) {
	const stmt = `SELECT id, person_id, type, status, notification_time, description
	FROM notifications
	WHERE status = 'pending' AND NOW() >= notification_time`

	var ns []Notification

	rows, err := r.getDB(ctx).QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var n Notification

		err := rows.Scan(&n.ID, &n.PersonID, &n.Type, &n.Status, &n.NotificationTime, &n.Description)
		if err != nil {
			return nil, err
		}

		ns = append(ns, n)
	}

	return ns, nil
}
