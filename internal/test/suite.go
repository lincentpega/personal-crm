package test

import (
	"context"
	"database/sql"
	"time"

	database "github.com/lincentpega/personal-crm/internal/db"
	"github.com/lincentpega/personal-crm/internal/log"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestSuite struct {
	suite.Suite
	Ctx         context.Context
	Log         *log.Logger
	PgContainer *postgres.PostgresContainer
	DB          *sql.DB
}

func (suite *TestSuite) SetupSuite() {
	suite.Ctx = context.Background()
	suite.Log = log.New()

	pgContainer, err := postgres.Run(suite.Ctx,
		"postgres:16.4-alpine",
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		suite.Log.ErrorLog.Fatal(err)
	}

	suite.PgContainer = pgContainer

	connStr, err := pgContainer.ConnectionString(suite.Ctx, "sslmode=disable")
	if err != nil {
		suite.Log.ErrorLog.Fatal(err)
	}

	db, err := database.Connect(connStr)
	if err != nil {
		suite.Log.ErrorLog.Fatal(err)
	}

	suite.DB = db

	if err := database.ExecMigrations(suite.DB, suite.Log); err != nil {
		suite.Log.ErrorLog.Fatal(err)
	}
}

func (suite *TestSuite) TearDownSuite() {
	if err := suite.PgContainer.Terminate(suite.Ctx); err != nil {
		suite.Log.ErrorLog.Fatalf("error terminating postgres container: %s", err)
	}
}
