package person

import (
	"database/sql"
	"testing"
	"time"

	"github.com/lincentpega/personal-crm/internal/common/txcontext"
	"github.com/lincentpega/personal-crm/internal/test"
	"github.com/stretchr/testify/suite"
)

const (
	testFirstName      = "John"
	testLastName       = "Smith"
	testSecondName     = "James"
	testBirthDate      = "2002-07-19"
	testMethodName     = "email"
	testContactData    = "john.smith@example.com"
	testCompany        = "Meta"
	testPosition       = "Platform SE L3"
	testCurrent        = true
	testBirthdayNotify = true
	testMethod1        = "telegram"
	testMethod2        = "phone"
	testData1          = "@paveldurov"
	testData2          = "+795554433"
	testCompany1       = "Meta"
	testPosition1      = "Platform SE L3"
	testCurrent1       = true
	testCompany2       = "Google"
	testPosition2      = "Junior SE"
	testCurrent2       = false
)

type personRepoTestSuite struct {
	test.TestSuite
	repo *PersonRepository
	tx   *sql.Tx
}

func (suite *personRepoTestSuite) SetupSuite() {
	suite.TestSuite.SetupSuite()

	suite.repo = NewRepository(suite.DB)
}

func (suite *personRepoTestSuite) SetupTest() {
	var err error
	suite.tx, err = suite.DB.BeginTx(suite.Ctx, nil)
	suite.Require().NoError(err)
}

func (suite *personRepoTestSuite) TearDownTest() {
	err := suite.tx.Rollback()
	suite.Require().NoError(err)
}

func (suite *personRepoTestSuite) TestGet() {
	ctx := txcontext.WithTx(suite.Ctx, suite.tx)

	var personID int
	stmt := `INSERT INTO persons (first_name, last_name, second_name, birth_date) 
	VALUES ($1, $2, $3, $4) RETURNING id`
	err := suite.tx.QueryRowContext(ctx, stmt, testFirstName, testLastName, testSecondName, testBirthDate).Scan(&personID)
	suite.NoError(err)

	stmt = `INSERT INTO contact_infos (person_id, method_name, contact_data) 
	VALUES ($1, $2, $3)`
	_, err = suite.tx.ExecContext(ctx, stmt, personID, testMethodName, testContactData)
	suite.NoError(err)

	stmt = `INSERT INTO job_infos (person_id, company, position, current) 
	VALUES ($1, $2, $3, $4)`
	_, err = suite.tx.ExecContext(ctx, stmt, personID, testCompany, testPosition, testCurrent)
	suite.NoError(err)

	stmt = `INSERT INTO person_settings (person_id, birthday_notify) 
	VALUES ($1, $2)`
	_, err = suite.tx.ExecContext(ctx, stmt, personID, testBirthdayNotify)
	suite.NoError(err)

	person, err := suite.repo.Get(ctx, personID)
	suite.NoError(err)
	suite.Equal(testFirstName, person.FirstName)
	suite.Equal(testLastName, person.LastName.String)
	suite.Equal(testSecondName, person.SecondName.String)

	expectedBirthDate, err := time.Parse("2006-01-02", testBirthDate)
	suite.NoError(err)
	expectedBirthDateUTC := expectedBirthDate.UTC()
	suite.Equal(expectedBirthDateUTC, person.BirthDate.Time.UTC())

	suite.Equal(1, len(person.ContactInfos))
	suite.Equal(testMethodName, person.ContactInfos[0].Method)
	suite.Equal(testContactData, person.ContactInfos[0].Data)
	suite.Equal(1, len(person.JobInfos))
	suite.Equal(testCompany, person.JobInfos[0].Company)
	suite.Equal(testPosition, person.JobInfos[0].Position)
	suite.Equal(testCurrent, person.JobInfos[0].Current)
	suite.Equal(testBirthdayNotify, person.Settings.BirthdayNotify)
}

func (suite *personRepoTestSuite) TestInsert() {
	ctx := txcontext.WithTx(suite.Ctx, suite.tx)

	birthDate := time.Date(2002, time.July, 19, 0, 0, 0, 0, time.UTC)

	ci1 := ContactInfo{Method: testMethod1, Data: testData1}
	ci2 := ContactInfo{Method: testMethod2, Data: testData2}
	contactInfos := []ContactInfo{ci1, ci2}

	ji1 := JobInfo{Company: testCompany1, Position: testPosition1, Current: testCurrent1}
	ji2 := JobInfo{Company: testCompany2, Position: testPosition2, Current: testCurrent2}
	jobInfos := []JobInfo{ji1, ji2}

	settings := Settings{BirthdayNotify: testBirthdayNotify}

	person := Person{
		FirstName:    testFirstName,
		LastName:     sql.NullString{String: testLastName, Valid: true},
		SecondName:   sql.NullString{String: testSecondName, Valid: true},
		BirthDate:    sql.NullTime{Time: birthDate, Valid: true},
		ContactInfos: contactInfos,
		JobInfos:     jobInfos,
		Settings:     settings,
	}

	suite.NoError(suite.repo.Insert(ctx, &person))

	insertedPerson, err := suite.repo.Get(ctx, person.ID)
	suite.NoError(err)
	suite.Equal(testFirstName, insertedPerson.FirstName)
	suite.Equal(testLastName, insertedPerson.LastName.String)
	suite.Equal(testSecondName, insertedPerson.SecondName.String)
	suite.Equal(birthDate.UTC(), insertedPerson.BirthDate.Time.UTC())
	suite.Equal(contactInfos, insertedPerson.ContactInfos)
	suite.Equal(jobInfos, insertedPerson.JobInfos)
	suite.Equal(settings, insertedPerson.Settings)
}

func TestPersonRepoTestSuite(t *testing.T) {
	suite.Run(t, new(personRepoTestSuite))
}
