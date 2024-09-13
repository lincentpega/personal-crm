package person

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lincentpega/personal-crm/internal/common/txcontext"
	"github.com/lincentpega/personal-crm/internal/models"
)

type PersonRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *PersonRepository {
	return &PersonRepository{db: db}
}

type DB interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func (m *PersonRepository) getDB(ctx context.Context) DB {
	if tx, ok := txcontext.GetTx(ctx); ok {
		return tx
	}
	return m.db
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
	if err := m.insertPerson(ctx, p); err != nil {
		return err
	}

	if err := m.insertContactInfos(ctx, p); err != nil {
		return err
	}

	if err := m.insertJobInfos(ctx, p); err != nil {
		return err
	}

	if err := m.insertSettings(ctx, p); err != nil {
		return err
	}

	return nil
}

func (m *PersonRepository) fetchPerson(ctx context.Context, id int, p *Person) error {
	const stmt = `SELECT id, first_name, last_name, second_name, birth_date 
        FROM persons 
        WHERE id = $1`

	err := m.getDB(ctx).QueryRowContext(ctx, stmt, id).Scan(&p.ID, &p.FirstName, &p.LastName, &p.SecondName, &p.BirthDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.ErrRecordNotFound
		}
		return err
	}

	return nil
}

func (m *PersonRepository) fetchContactInfos(ctx context.Context, id int, p *Person) error {
	const stmt = `SELECT method_name, contact_data 
        FROM contact_infos 
        WHERE person_id = $1`

	rows, err := m.getDB(ctx).QueryContext(ctx, stmt, id)
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
	const stmt = `SELECT company, position, current 
        FROM job_infos 
        WHERE person_id = $1`

	rows, err := m.getDB(ctx).QueryContext(ctx, stmt, id)
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

	err := m.getDB(ctx).QueryRowContext(ctx, stmt, id).Scan(&p.Settings.BirthdayNotify)
	if err != nil {
		return err
	}

	return nil
}

func (m *PersonRepository) insertPerson(ctx context.Context, p *Person) error {
	const stmt = `INSERT INTO persons (first_name, last_name, second_name, birth_date)
        VALUES($1, $2, $3, $4) RETURNING id`

	m.getDB(ctx).QueryRowContext(ctx, stmt, p.FirstName, p.LastName, p.SecondName, p.BirthDate).Scan(&p.ID)

	return nil
}

func (m *PersonRepository) insertContactInfos(ctx context.Context, p *Person) error {
	const stmt = `INSERT INTO contact_infos (person_id, method_name, contact_data) 
        VALUES($1, $2, $3)`

	for _, c := range p.ContactInfos {
		_, err := m.getDB(ctx).ExecContext(ctx, stmt, p.ID, c.Method, c.Data)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *PersonRepository) insertJobInfos(ctx context.Context, p *Person) error {
	const stmt = `INSERT INTO job_infos (person_id, company, position, current) 
        VALUES($1, $2, $3, $4)`

	for _, j := range p.JobInfos {
		_, err := m.getDB(ctx).ExecContext(ctx, stmt, p.ID, j.Company, j.Position, j.Current)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *PersonRepository) insertSettings(ctx context.Context, p *Person) error {
	const stmt = `INSERT INTO person_settings (person_id, birthday_notify)
        VALUES($1, $2)`

	_, err := m.getDB(ctx).ExecContext(ctx, stmt, p.ID, p.Settings.BirthdayNotify)
	if err != nil {
		return err
	}

	return nil
}
