package notifications

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/lincentpega/personal-crm/internal/common/txcontext"
	"github.com/lincentpega/personal-crm/internal/models/person"
	"github.com/lincentpega/personal-crm/internal/test"
	"github.com/stretchr/testify/suite"
)

const (
	testPersonFirstName  = "John"
	testPersonLastName   = "Smith"
	testPersonBirthDate  = "2000-01-01"
	testNotifType        = KeepInTouch
	testNotifStatus      = Pending
	testNotifDescription = "Keep in touch with John Smith"
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
	ctx := txcontext.WithTx(suite.Ctx, suite.tx)

	person := suite.createTestPerson(ctx)

	notifTime := time.Now().Add(24 * time.Hour).UTC()

	var notifID int
	stmt := `INSERT INTO notifications (person_id, type, status, notification_time, description) 
	VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := suite.tx.QueryRowContext(ctx, stmt, person.ID, testNotifType, testNotifStatus, notifTime, testNotifDescription).Scan(&notifID)
	suite.NoError(err)

	notification, err := suite.notifRepo.Get(ctx, notifID)
	suite.NoError(err)
	suite.Equal(notifID, notification.ID)
	suite.Equal(person.ID, notification.PersonID)
	suite.Equal(testNotifType, notification.Type)
	suite.Equal(testNotifStatus, notification.Status)
	suite.Equal(notifTime.UTC(), notification.NotificationTime.UTC())
	suite.Equal(testNotifDescription, notification.Description)
}

func (suite *notificationRepoTestSuite) TestInsert() {
	ctx := txcontext.WithTx(suite.Ctx, suite.tx)

	person := suite.createTestPerson(ctx)

	notifTime := time.Now().Add(time.Hour * 24).UTC()
	notif := Notification{
		PersonID:         person.ID,
		Type:             testNotifType,
		Status:           testNotifStatus,
		NotificationTime: notifTime,
		Description:      testNotifDescription,
	}

	err := suite.notifRepo.Insert(ctx, &notif)
	suite.Require().NoError(err)
	suite.Require().NotZero(notif.ID)

	insertedNotif, err := suite.notifRepo.Get(ctx, notif.ID)
	suite.Require().NoError(err)
	suite.Require().Equal(person.ID, insertedNotif.PersonID)
	suite.Require().Equal(testNotifType, insertedNotif.Type)
	suite.Require().Equal(testNotifStatus, insertedNotif.Status)
	suite.Require().Equal(notifTime.UTC(), insertedNotif.NotificationTime.UTC())
	suite.Require().Equal(testNotifDescription, insertedNotif.Description)
}

func (suite *notificationRepoTestSuite) TestUpdateNotificationStatus() {
	ctx := txcontext.WithTx(suite.Ctx, suite.tx)

	person := suite.createTestPerson(ctx)

	notifTime := time.Now().Add(time.Hour * 24).UTC()
	notif := Notification{
		PersonID:         person.ID,
		Type:             testNotifType,
		Status:           Pending,
		NotificationTime: notifTime,
		Description:      testNotifDescription,
	}

	err := suite.notifRepo.Insert(ctx, &notif)
	suite.Require().NoError(err)

	err = suite.notifRepo.UpdateNotificationStatus(ctx, notif.ID, Failed)
	suite.Require().NoError(err)

	updatedNotif, err := suite.notifRepo.Get(ctx, notif.ID)
	suite.Require().NoError(err)
	suite.Require().Equal(Failed, updatedNotif.Status)
}

func (suite *notificationRepoTestSuite) TestGetAwaitingSend() {
	ctx := txcontext.WithTx(suite.Ctx, suite.tx)

	person := suite.createTestPerson(ctx)

	notifTime := time.Now().Add(-time.Hour * 24).UTC()
	shouldGetNotif := Notification{
		PersonID:         person.ID,
		Type:             testNotifType,
		Status:           Pending,
		NotificationTime: notifTime,
		Description:      testNotifDescription,
	}

	notifTime = time.Now().Add(time.Hour * 24).UTC()
	earlyNotif := Notification{
		PersonID:         person.ID,
		Type:             testNotifType,
		Status:           Pending,
		NotificationTime: notifTime,
		Description:      testNotifDescription,
	}

	notifTime = time.Now().Add(-time.Hour * 24).UTC()
	failedNotif := Notification{
		PersonID:         person.ID,
		Type:             testNotifType,
		Status:           Failed,
		NotificationTime: notifTime,
		Description:      testNotifDescription,
	}

	err := suite.notifRepo.Insert(ctx, &shouldGetNotif)
	suite.Require().NoError(err)

	err = suite.notifRepo.Insert(ctx, &earlyNotif)
	suite.Require().NoError(err)

	err = suite.notifRepo.Insert(ctx, &failedNotif)
	suite.Require().NoError(err)

	notifs, err := suite.notifRepo.GetAwaitingSend(ctx)
	suite.Require().NoError(err)
	suite.Require().Equal(1, len(notifs))
	suite.Require().Equal(shouldGetNotif.ID, notifs[0].ID)
}

func (suite *notificationRepoTestSuite) createTestPerson(ctx context.Context) *person.Person {
	pBirthDate, err := time.Parse("2006-01-02", testPersonBirthDate)
	suite.Require().NoError(err)

	person := &person.Person{
		FirstName: testPersonFirstName,
		LastName:  sql.NullString{String: testPersonLastName, Valid: true},
		BirthDate: sql.NullTime{Time: pBirthDate, Valid: true},
	}

	err = suite.personRepo.Insert(ctx, person)
	suite.Require().NoError(err)
	suite.Require().NotZero(person.ID)

	return person
}

func TestNotificationRepoTestSuite(t *testing.T) {
	suite.Run(t, new(notificationRepoTestSuite))
}
