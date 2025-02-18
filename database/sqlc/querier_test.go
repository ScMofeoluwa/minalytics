package database

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type DatabaseSuite struct {
	suite.Suite
	container testcontainers.Container
	querier   Querier
	ctx       context.Context
}

func setupQuerier(suite *DatabaseSuite, dsn string) *Queries {
	connPool, err := pgxpool.New(suite.ctx, dsn)
	if err != nil {
		suite.FailNowf("unable to create connection pool", err.Error())
	}

	return New(connPool)
}

func migrateDB(suite *DatabaseSuite, dsn string) {
	m, err := migrate.New("file://../migrations", dsn)
	if err != nil {
		suite.FailNowf("failed create migrate instance", err.Error())
	}

	if err := m.Up(); err != nil {
		suite.FailNowf("failed to migrate db", err.Error())
	}
}

func (suite *DatabaseSuite) SetupSuite() {
	suite.ctx = context.Background()

	// setup test container
	req := testcontainers.ContainerRequest{
		Image:        "timescale/timescaledb:latest-pg17",
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor: wait.Strategy(
			wait.ForAll(
				wait.ForListeningPort("5432/tcp"),
				wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(5*time.Second),
			),
		),
		Env: map[string]string{
			"POSTGRES_DB":       "test_db",
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "postgres",
		},
	}
	container, err := testcontainers.GenericContainer(suite.ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	suite.NoError(err)
	suite.container = container

	port, err := container.MappedPort(suite.ctx, "5432/tcp")
	suite.NoError(err)

	dsn := fmt.Sprintf("postgres://postgres:postgres@localhost:%s/test_db?sslmode=disable", port.Port())

	// migrate
	migrateDB(suite, dsn)

	// setup querier
	suite.querier = setupQuerier(suite, dsn)
}

func (suite *DatabaseSuite) TearDownSuite() {
	err := suite.container.Terminate(suite.ctx)
	suite.NoError(err)
}

// Helper functions
func (suite *DatabaseSuite) createTestUser() uuid.UUID {
	email := faker.Email()
	userID, err := suite.querier.GetOrCreateUser(suite.ctx, email)
	suite.NoError(err)
	return userID
}

func (suite *DatabaseSuite) createTestApp(userID uuid.UUID) App {
	app, err := suite.querier.CreateApp(suite.ctx, CreateAppParams{
		Name:   faker.Word(),
		UserID: userID,
	})
	suite.NoError(err)
	return app
}

func (suite *DatabaseSuite) createTestEvent(trackingID uuid.UUID) {
	err := suite.querier.CreateEvent(suite.ctx, CreateEventParams{
		VisitorID:       faker.Word(),
		TrackingID:      trackingID,
		EventType:       "pageview",
		Url:             stringPtr(faker.URL()),
		Referrer:        stringPtr(faker.URL()),
		Country:         faker.GetCountryInfo().Name,
		Browser:         "Safari",
		Device:          "iPhone",
		OperatingSystem: "iOS",
		Details:         EventDetails{},
	})
	suite.NoError(err)
}

func (suite *DatabaseSuite) TestGetOrCreateUser() {
	t := suite.T()

	testCases := []struct {
		name     string
		email    string
		hasError bool
	}{
		{
			name:     "create new user",
			email:    "user@gmail.com",
			hasError: false,
		},
		{
			name:     "returns existing user",
			email:    "user@gmail.com",
			hasError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := suite.querier.GetOrCreateUser(suite.ctx, tc.email)
			if tc.hasError {
				suite.Error(err)
				return
			}
			suite.NoError(err)
		})
	}
}

func (suite *DatabaseSuite) TestCreateApp() {
	t := suite.T()

	userID, _ := suite.querier.GetOrCreateUser(suite.ctx, faker.Email())

	testCases := []struct {
		name     string
		actual   string
		arg      CreateAppParams
		hasError bool
	}{
		{
			name: "create new app",
			arg: CreateAppParams{
				Name:   "testapp",
				UserID: userID,
			},
			hasError: false,
		},
		{
			name: "duplicate app name for same user",
			arg: CreateAppParams{
				Name:   "testapp",
				UserID: userID,
			},
			hasError: true,
		},
		{
			name: "invalid userID",
			arg: CreateAppParams{
				Name:   faker.Name(),
				UserID: uuid.New(),
			},
			hasError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := suite.querier.CreateApp(suite.ctx, tc.arg)
			if tc.hasError {
				suite.Error(err)
				return
			}
			suite.NoError(err)
		})
	}
}

func (suite *DatabaseSuite) TestCheckAppExists() {
	t := suite.T()

	userID := suite.createTestUser()
	app := suite.createTestApp(userID)

	testCases := []struct {
		name        string
		arg         CheckAppExistsParams
		expectedErr error
	}{
		{
			name: "app exists",
			arg: CheckAppExistsParams{
				Name:   app.Name,
				UserID: userID,
			},
			expectedErr: nil,
		},
		{
			name: "app doesn't exist",
			arg: CheckAppExistsParams{
				Name:   faker.Name(),
				UserID: userID,
			},
			expectedErr: pgx.ErrNoRows,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := suite.querier.CheckAppExists(suite.ctx, tc.arg)
			if tc.expectedErr != nil {
				suite.Error(err)
				suite.Equal(tc.expectedErr.Error(), err.Error())
				return
			}
			suite.NoError(err)
		})
	}
}

func (suite *DatabaseSuite) TestCreateEvent() {
	userID := suite.createTestUser()
	app := suite.createTestApp(userID)

	err := suite.querier.CreateEvent(suite.ctx, CreateEventParams{
		VisitorID:       faker.Word(),
		TrackingID:      app.TrackingID,
		EventType:       "pageview",
		Url:             stringPtr(faker.URL()),
		Referrer:        stringPtr(faker.URL()),
		Country:         faker.GetCountryInfo().Name,
		Browser:         "Safari",
		Device:          "iPhone",
		OperatingSystem: "iOS",
		Details:         EventDetails{},
	})
	suite.NoError(err)
}

func (suite *DatabaseSuite) TestGetAppByTrackingID() {
	userID := suite.createTestUser()
	app := suite.createTestApp(userID)

	app_, err := suite.querier.GetAppByTrackingID(suite.ctx, app.TrackingID)
	suite.NoError(err)
	suite.Equal(app.ID, app_.ID)
	suite.Equal(app.Name, app_.Name)
}

func (suite *DatabaseSuite) TestGetApps() {
	userID := suite.createTestUser()
	suite.createTestApp(userID)

	apps, err := suite.querier.GetApps(suite.ctx, userID)
	suite.NoError(err)
	suite.Greater(len(apps), 0)
}

func (suite *DatabaseSuite) TestGetVisitors() {
	userID := suite.createTestUser()
	app := suite.createTestApp(userID)
	suite.createTestEvent(app.TrackingID)

	visitors, err := suite.querier.GetVisitors(suite.ctx, GetVisitorsParams{
		TrackingID: app.TrackingID,
		Column2:    sql.NullTime{},
		Column3:    sql.NullTime{},
		TimeBucket: "1 hour",
	})
	suite.NoError(err)
	suite.Greater(len(visitors), 0)
}

func (suite *DatabaseSuite) TestGetPageViews() {
	userID := suite.createTestUser()
	app := suite.createTestApp(userID)
	suite.createTestEvent(app.TrackingID)

	pageViews, err := suite.querier.GetPageViews(suite.ctx, GetPageViewsParams{
		TrackingID: app.TrackingID,
		Column2:    sql.NullTime{},
		Column3:    sql.NullTime{},
		TimeBucket: "1 hour",
	})
	suite.NoError(err)
	suite.Greater(len(pageViews), 0)
}

func (suite *DatabaseSuite) TestGetReferrals() {
	userID := suite.createTestUser()
	app := suite.createTestApp(userID)
	suite.createTestEvent(app.TrackingID)

	referrals, err := suite.querier.GetReferrals(suite.ctx, GetReferralsParams{
		TrackingID: app.TrackingID,
		Column2:    sql.NullTime{},
		Column3:    sql.NullTime{},
	})
	suite.NoError(err)
	suite.Greater(len(referrals), 0)
}

func (suite *DatabaseSuite) TestGetPages() {
	userID := suite.createTestUser()
	app := suite.createTestApp(userID)
	suite.createTestEvent(app.TrackingID)

	pages, err := suite.querier.GetPages(suite.ctx, GetPagesParams{
		TrackingID: app.TrackingID,
		Column2:    sql.NullTime{},
		Column3:    sql.NullTime{},
	})
	suite.NoError(err)
	suite.Greater(len(pages), 0)
}

func (suite *DatabaseSuite) TestGetCountries() {
	userID := suite.createTestUser()
	app := suite.createTestApp(userID)
	suite.createTestEvent(app.TrackingID)

	countries, err := suite.querier.GetCountries(suite.ctx, GetCountriesParams{
		TrackingID: app.TrackingID,
		Column2:    sql.NullTime{},
		Column3:    sql.NullTime{},
	})
	suite.NoError(err)
	suite.Greater(len(countries), 0)
}

func (suite *DatabaseSuite) TestGetBrowsers() {
	userID := suite.createTestUser()
	app := suite.createTestApp(userID)
	suite.createTestEvent(app.TrackingID)

	browsers, err := suite.querier.GetBrowsers(suite.ctx, GetBrowsersParams{
		TrackingID: app.TrackingID,
		Column2:    sql.NullTime{},
		Column3:    sql.NullTime{},
	})
	suite.NoError(err)
	suite.Greater(len(browsers), 0)
}

func (suite *DatabaseSuite) TestGetDevices() {
	userID := suite.createTestUser()
	app := suite.createTestApp(userID)
	suite.createTestEvent(app.TrackingID)

	devices, err := suite.querier.GetDevices(suite.ctx, GetDevicesParams{
		TrackingID: app.TrackingID,
		Column2:    sql.NullTime{},
		Column3:    sql.NullTime{},
	})
	suite.NoError(err)
	suite.Greater(len(devices), 0)
}

func (suite *DatabaseSuite) TestGetOS() {
	userID := suite.createTestUser()
	app := suite.createTestApp(userID)
	suite.createTestEvent(app.TrackingID)

	os, err := suite.querier.GetOS(suite.ctx, GetOSParams{
		TrackingID: app.TrackingID,
		Column2:    sql.NullTime{},
		Column3:    sql.NullTime{},
	})
	suite.NoError(err)
	suite.Greater(len(os), 0)
}

func stringPtr(s string) *string {
	return &s
}

func TestDatabaseSuite(t *testing.T) {
	suite.Run(t, new(DatabaseSuite))
}
