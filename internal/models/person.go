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

type PersonModel struct {
	DB *sql.DB
}

func NewPersonModel(db *sql.DB) *PersonModel {
    return &PersonModel{DB: db}
}

func (m *PersonModel) Get(ctx context.Context, id int) (*Person, error) {
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

	return &p, nil
}

func (m *PersonModel) fetchPerson(ctx context.Context, id int, p *Person) error {
	const personStmt = `SELECT id, first_name, last_name, second_name, birth_date
                        FROM persons
                        WHERE id = $1`

	err := m.DB.QueryRowContext(ctx, personStmt, id).Scan(&p.ID, &p.FirstName, &p.LastName, &p.SecondName, &p.BirthDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrRecordNotFound
		}
		return err
	}

	return nil
}

func (m *PersonModel) fetchContactInfos(ctx context.Context, id int, p *Person) error {
	const contactStmt = `SELECT method_name, contact_data
                         FROM contact_infos
                         WHERE person_id = $1`

	rows, err := m.DB.QueryContext(ctx, contactStmt, id)
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

func (m *PersonModel) fetchJobInfos(ctx context.Context, id int, p *Person) error {
	const jobStmt = `SELECT company, job_position, is_current
                     FROM job_infos
                     WHERE person_id = $1`

	rows, err := m.DB.QueryContext(ctx, jobStmt, id)
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

func (m *PersonModel) Insert(ctx context.Context, p *Person) error {
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

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (m *PersonModel) insertPerson(ctx context.Context, tx *sql.Tx, p *Person) error {
	const personStmt = `INSERT INTO persons (id, first_name, last_name, second_name, birth_date)
                        VALUES($1, $2, $3, $4, $5)`

	_, err := tx.ExecContext(ctx, personStmt, p.ID, p.FirstName, p.LastName, p.SecondName, p.BirthDate)
	if err != nil {
		return err
	}

	return nil
}

func (m *PersonModel) insertContactInfos(ctx context.Context, tx *sql.Tx, p *Person) error {
	const contactStmt = `INSERT INTO contact_infos (person_id, method_name, contact_data)
                         VALUES($1, $2, $3)`

	for _, c := range p.ContactInfos {
		_, err := tx.ExecContext(ctx, contactStmt, p.ID, c.Method, c.Data)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *PersonModel) insertJobInfos(ctx context.Context, tx *sql.Tx, p *Person) error {
	const jobStmt = `INSERT INTO job_infos (person_id, company, job_position, is_current)
                     VALUES($1, $2, $3, $4)`

	for _, j := range p.JobInfos {
		_, err := tx.ExecContext(ctx, jobStmt, p.ID, j.Company, j.Position, j.Current)
		if err != nil {
			return err
		}
	}

	return nil
}
