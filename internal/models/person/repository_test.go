package person

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/lincentpega/personal-crm/internal/common/txcontext"
	database "github.com/lincentpega/personal-crm/internal/db"
	"github.com/lincentpega/personal-crm/internal/log"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type personRepoTestSuite struct {
	suite.Suite
	db          *sql.DB
	pgContainer *postgres.PostgresContainer
	ctx         context.Context
	repo        *PersonRepository
	log         *log.Logger
	tx          *sql.Tx
}

func (suite *personRepoTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	suite.log = log.New()

	pgContainer, err := postgres.Run(suite.ctx,
		"postgres:16.4-alpine",
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		suite.log.ErrorLog.Fatal(err)
	}

	suite.pgContainer = pgContainer

	connStr, err := pgContainer.ConnectionString(suite.ctx, "sslmode=disable")
	if err != nil {
		suite.log.ErrorLog.Fatal(err)
	}

	db, err := database.Connect(connStr)
	if err != nil {
		suite.log.ErrorLog.Fatal(err)
	}

	suite.db = db

	if err := database.ExecMigrations(suite.db, suite.log); err != nil {
		suite.log.ErrorLog.Fatal(err)
	}

	repo := NewRepository(db)

	suite.repo = repo
}

func (suite *personRepoTestSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		suite.log.ErrorLog.Fatalf("error terminating postgres container: %s", err)
	}
}

func (suite *personRepoTestSuite) SetupTest() {
	var err error
	suite.tx, err = suite.db.BeginTx(suite.ctx, nil)
	suite.Require().NoError(err)
}

func (suite *personRepoTestSuite) TearDownTest() {
	err := suite.tx.Rollback()
	suite.Require().NoError(err)
}

func (suite *personRepoTestSuite) TestGet() {
	firstName := "John"
	lastName := "Smith"
	secondName := "James"
	birthDate := "2002-07-19"

	ctx := txcontext.WithTx(suite.ctx, suite.tx)

	var personID int
	stmt := `INSERT INTO persons (first_name, last_name, second_name, birth_date) 
	VALUES ($1, $2, $3, $4) RETURNING id`
	err := suite.tx.QueryRowContext(ctx, stmt, firstName, lastName, secondName, birthDate).Scan(&personID)
	suite.NoError(err)

	methodName := "email"
	contactData := "john.smith@example.com"
	stmt = `INSERT INTO contact_infos (person_id, method_name, contact_data) 
	VALUES ($1, $2, $3)`
	_, err = suite.tx.ExecContext(ctx, stmt, personID, methodName, contactData)
	suite.NoError(err)

	company := "Meta"
	position := "Platform SE L3"
	current := true
	stmt = `INSERT INTO job_infos (person_id, company, position, current) 
	VALUES ($1, $2, $3, $4)`
	_, err = suite.tx.ExecContext(ctx, stmt, personID, company, position, current)
	suite.NoError(err)

	birthdayNotify := true
	stmt = `INSERT INTO person_settings (person_id, birthday_notify) 
	VALUES ($1, $2)`
	_, err = suite.tx.ExecContext(ctx, stmt, personID, birthdayNotify)
	suite.NoError(err)

	person, err := suite.repo.Get(ctx, personID)
	suite.NoError(err)
	suite.Equal(firstName, person.FirstName)
	suite.Equal(lastName, *person.LastName)
	suite.Equal(secondName, *person.SecondName)

	expectedBirthDate, err := time.Parse("2006-01-02", birthDate)
	suite.NoError(err)
	expectedBirthDateUTC := expectedBirthDate.UTC()
	suite.Equal(expectedBirthDateUTC, person.BirthDate.UTC())

	suite.Equal(1, len(person.ContactInfos))
	suite.Equal(methodName, person.ContactInfos[0].Method)
	suite.Equal(contactData, person.ContactInfos[0].Data)
	suite.Equal(1, len(person.JobInfos))
	suite.Equal(company, person.JobInfos[0].Company)
	suite.Equal(position, person.JobInfos[0].Position)
	suite.Equal(current, person.JobInfos[0].Current)
	suite.Equal(birthdayNotify, person.Settings.BirthdayNotify)
}

func (suite *personRepoTestSuite) TestInsert() {
	firstName := "John"
	lastName := "Smith"
	secondName := "James"
	birthDate := time.Date(2002, time.July, 19, 0, 0, 0, 0, time.UTC)

	method1 := "telegram"
	method2 := "phone"
	data1 := "@paveldurov"
	data2 := "+795554433"
	ci1 := ContactInfo{Method: method1, Data: data1}
	ci2 := ContactInfo{Method: method2, Data: data2}
	contactInfos := []ContactInfo{ci1, ci2}

	company1 := "Meta"
	position1 := "Platform SE L3"
	current1 := true
	company2 := "Google"
	position2 := "Junior SE"
	current2 := false
	ji1 := JobInfo{Company: company1, Position: position1, Current: current1}
	ji2 := JobInfo{Company: company2, Position: position2, Current: current2}
	jobInfos := []JobInfo{ji1, ji2}

	birthdayNotify := true
	settings := Settings{BirthdayNotify: birthdayNotify}

	person := Person{
		FirstName:    firstName,
		LastName:     &lastName,
		SecondName:   &secondName,
		BirthDate:    &birthDate,
		ContactInfos: contactInfos,
		JobInfos:     jobInfos,
		Settings:     settings,
	}

	suite.NoError(suite.repo.Insert(suite.ctx, &person))

	insertedPerson, err := suite.repo.Get(suite.ctx, person.ID)
	suite.NoError(err)
	suite.Equal(firstName, insertedPerson.FirstName)
	suite.Equal(lastName, *insertedPerson.LastName)
	suite.Equal(secondName, *insertedPerson.SecondName)
	suite.Equal(birthDate.UTC(), insertedPerson.BirthDate.UTC())
	suite.Equal(contactInfos, insertedPerson.ContactInfos)
	suite.Equal(jobInfos, insertedPerson.JobInfos)
	suite.Equal(settings, insertedPerson.Settings)
}

func TestPersonRepoTestSuite(t *testing.T) {
	suite.Run(t, new(personRepoTestSuite))
}
