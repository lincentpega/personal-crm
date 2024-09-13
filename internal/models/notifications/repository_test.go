package notifications

import (
	"database/sql"
	"testing"
	"time"

	"github.com/lincentpega/personal-crm/internal/common/txcontext"
	"github.com/lincentpega/personal-crm/internal/models/person"
	"github.com/lincentpega/personal-crm/internal/test"
	"github.com/stretchr/testify/suite"
)

type notificationRepoTestSuite struct {
	test.TestSuite
	notifRepo  *NotificationRepository
	personRepo *person.PersonRepository
	tx         *sql.Tx
}

func (suite *notificationRepoTestSuite) SetupSuite() {
	suite.TestSuite.SetupSuite()

	suite.notifRepo = NewRepository(suite.DB)
}

func (suite *notificationRepoTestSuite) SetupTest() {
	var err error
	suite.tx, err = suite.DB.BeginTx(suite.Ctx, nil)
	suite.Require().NoError(err)
}

func (suite *notificationRepoTestSuite) TearDownTest() {
	err := suite.tx.Rollback()
	suite.Require().NoError(err)
}

func (suite *notificationRepoTestSuite) TestGet() {
	personFirstName := "John"
	personLastName := "Smith"
	personBirthDate := "2000-01-01"

	ctx := txcontext.WithTx(suite.Ctx, suite.tx)

	var personID int
	stmt := `INSERT INTO persons (first_name, last_name, birth_date) 
	VALUES ($1, $2, $3) RETURNING id`
	err := suite.tx.QueryRowContext(ctx, stmt, personFirstName, personLastName, personBirthDate).Scan(&personID)
	suite.NoError(err)

	notifType := KeepInTouch
	notifStatus := Pending
	notifTime := time.Now().Add(24 * time.Hour).UTC()
	notifDescription := "Keep in touch with John Smith"

	var notifID int
	stmt = `INSERT INTO notifications (person_id, type, status, notification_time, description) 
	VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err = suite.tx.QueryRowContext(ctx, stmt, personID, notifType, notifStatus, notifTime, notifDescription).Scan(&notifID)
	suite.NoError(err)

	notification, err := suite.notifRepo.Get(ctx, notifID)
	suite.NoError(err)
	suite.Equal(notifID, notification.ID)
	suite.Equal(personID, notification.PersonID)
	suite.Equal(notifType, notification.Type)
	suite.Equal(notifStatus, notification.Status)
	suite.Equal(notifTime.UTC(), notification.NotificationTime.UTC())
	suite.Equal(notifDescription, notification.Description)
}

func (suite *notificationRepoTestSuite) TestInsert() {
	ctx := txcontext.WithTx(suite.Ctx, suite.tx)
	pLastName := "Smith"
	pBirthDate := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	person := person.Person{
		FirstName: "John",
		LastName:  &pLastName,
		BirthDate: &pBirthDate,
	}

	err := suite.personRepo.Insert(ctx, &person)
	suite.Require().NoError(err)
	suite.Require().NotZero(person.ID)

	notifTime := time.Now().Add(time.Hour * 24).UTC()
	notifType := KeepInTouch
	notifStatus := Pending
	notifDescription := "Keep in touch with John Smith"
	notif := Notification{
		PersonID:         person.ID,
		Type:             notifType,
		Status:           notifStatus,
		NotificationTime: &notifTime,
		Description:      notifDescription,
	}

	err = suite.notifRepo.Insert(ctx, &notif)
	suite.Require().NoError(err)
	suite.Require().NotZero(notif.ID)

	insertedNotif, err := suite.notifRepo.Get(ctx, notif.ID)
	suite.Require().NoError(err)
	suite.Require().Equal(person.ID, insertedNotif.PersonID)
	suite.Require().Equal(notifType, insertedNotif.Type)
	suite.Require().Equal(notifStatus, insertedNotif.Status)
	suite.Require().Equal(notifTime.UTC(), insertedNotif.NotificationTime.UTC())
	suite.Require().Equal(notifDescription, insertedNotif.Description)
}

func TestNotificationRepoTestSuite(t *testing.T) {
	suite.Run(t, new(notificationRepoTestSuite))
}
