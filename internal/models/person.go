package models

import (
	"context"
	"database/sql"
	"errors"
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

type PersonRepository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *PersonRepository {
	return &PersonRepository{DB: db}
}

func (m *PersonRepository) Get(ctx context.Context, id int) (*Person, error) {
	var p Person

	if err := m.fetchPerson(ctx, id, &p); err != nil {
		return nil, err
	}

	if err := m.fetchContactInfos(ctx, id, &p); err != nil {
		return nil, err
	}

	if err := m.fetchJobInfos(ctx, id, &p); err != nil {
		return nil, err
	}

	if err := m.fetchPersonSettings(ctx, id, &p); err != nil {
		return nil, err
	}

	return &p, nil
}

func (m *PersonRepository) Insert(ctx context.Context, p *Person) error {
	tx, err := m.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := m.insertPerson(ctx, tx, p); err != nil {
		return err
	}

	if err := m.insertContactInfos(ctx, tx, p); err != nil {
		return err
	}

	if err := m.insertJobInfos(ctx, tx, p); err != nil {
		return err
	}

	if err := m.insertSettings(ctx, tx, p); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (m *PersonRepository) fetchPerson(ctx context.Context, id int, p *Person) error {
	const stmt = `SELECT id, first_name, last_name, second_name, birth_date 
        FROM persons 
        WHERE id = $1`

	err := m.DB.QueryRowContext(ctx, stmt, id).Scan(&p.ID, &p.FirstName, &p.LastName, &p.SecondName, &p.BirthDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrRecordNotFound
		}
		return err
	}

	return nil
}

func (m *PersonRepository) fetchContactInfos(ctx context.Context, id int, p *Person) error {
	const stmt = `SELECT method_name, contact_data 
        FROM contact_infos 
        WHERE person_id = $1`

	rows, err := m.DB.QueryContext(ctx, stmt, id)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var c ContactInfo
		if err := rows.Scan(&c.Method, &c.Data); err != nil {
			return err
		}
		p.ContactInfos = append(p.ContactInfos, c)
	}

	return rows.Err()
}

func (m *PersonRepository) fetchJobInfos(ctx context.Context, id int, p *Person) error {
	const stmt = `SELECT company, job_position, is_current 
        FROM job_infos 
        WHERE person_id = $1`

	rows, err := m.DB.QueryContext(ctx, stmt, id)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var j JobInfo
		if err := rows.Scan(&j.Company, &j.Position, &j.Current); err != nil {
			return err
		}
		p.JobInfos = append(p.JobInfos, j)
	}

	return rows.Err()
}

func (m *PersonRepository) fetchPersonSettings(ctx context.Context, id int, p *Person) error {
	const stmt = `SELECT birthday_notify
        FROM person_settings
        WHERE person_id = $1`

	err := m.DB.QueryRowContext(ctx, stmt, id).Scan(&p.Settings.BirthdayNotify)
	if err != nil {
		return err
	}

	return nil
}

func (m *PersonRepository) insertPerson(ctx context.Context, tx *sql.Tx, p *Person) error {
	const stmt = `INSERT INTO persons (id, first_name, last_name, second_name, birth_date)
        VALUES($1, $2, $3, $4, $5)`

	_, err := tx.ExecContext(ctx, stmt, p.ID, p.FirstName, p.LastName, p.SecondName, p.BirthDate)
	if err != nil {
		return err
	}

	return nil
}

func (m *PersonRepository) insertContactInfos(ctx context.Context, tx *sql.Tx, p *Person) error {
	const stmt = `INSERT INTO contact_infos (person_id, method_name, contact_data) 
        VALUES($1, $2, $3)`

	for _, c := range p.ContactInfos {
		_, err := tx.ExecContext(ctx, stmt, p.ID, c.Method, c.Data)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *PersonRepository) insertJobInfos(ctx context.Context, tx *sql.Tx, p *Person) error {
	const stmt = `INSERT INTO job_infos (person_id, company, job_position, is_current) 
        VALUES($1, $2, $3, $4)`

	for _, j := range p.JobInfos {
		_, err := tx.ExecContext(ctx, stmt, p.ID, j.Company, j.Position, j.Current)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *PersonRepository) insertSettings(ctx context.Context, tx *sql.Tx, p *Person) error {
	const stmt = `INSERT INTO person_settings (person_id, birthday_notify)
        VALUES($1, $2)`

	_, err := tx.ExecContext(ctx, stmt, p.ID, p.Settings.BirthdayNotify)
	if err != nil {
		return err
	}

	return nil
}
