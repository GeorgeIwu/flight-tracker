package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"service-catalog/ent"
	_service "service-catalog/service"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"service-catalog/service/mocks"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

//Unit Tests
type UnitTestSuite struct {
	suite.Suite
	repo *mocks.ServiceRepoInterface
	usecase _service.ServiceUsecase
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, &UnitTestSuite{})
}

func (uts *UnitTestSuite) SetupTest() {
	repoMock := mocks.ServiceRepoInterface{}
	usecase := _service.NewServiceUsecase(&repoMock)

	uts.repo = &repoMock
	uts.usecase = *usecase
}

func (uts *UnitTestSuite) TestFetch() {
	newService := ent.Service{ Title: "Test", ID: 54, Description: "testing"}
	services := []*ent.Service{}
	services = append(services, &newService)

	uts.repo.On("Get", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(services, nil)
	ctx := context.TODO()
	searchBy := ""
	sortBy := ""
	pageNumber := 1
	itemsPerPage := 12

	actual, err := uts.usecase.Fetch(ctx, searchBy, sortBy, pageNumber, itemsPerPage)

	uts.Equal(services[0].ID, actual[0].ID)
	uts.Equal(nil, err)
}

func (uts *UnitTestSuite) TestCreate() {
	service := ent.Service{ Title: "Test", ID: 55, Description: "testing"}

	uts.repo.On("Create", mock.Anything, mock.Anything).Return(&service, nil)
	ctx := context.TODO()
	request := _service.ServiceRequest{}

	actual, err := uts.usecase.Create(ctx, request)

	uts.Equal(service.ID, actual.ID)
	uts.Equal(nil, err)
}

//Integration Tests
type IntTestSuite struct {
	suite.Suite
	client         *ent.Client
	usecase 		_service.ServiceUsecase
}

func TestIntTestSuite(t *testing.T) {
	suite.Run(t, &IntTestSuite{})
}

func (its *IntTestSuite) SetupSuite() {
	const (
		dbHost     = "localhost"
		dbUser     = "kong"
		dbPassword = "password"
		dbType    = "mysql"
		dbProtocol = "tcp"
		dbName = "testcatalog"
	)

	sqlInfo := fmt.Sprintf("%s:%s@%s(%s)/%s?parseTime=True", dbUser, dbPassword, dbProtocol, dbHost, dbName)
	client, err := ent.Open(dbType, sqlInfo)
	if err != nil {
			log.Fatalf("failed opening connection to sqlite: %v", err)
	}

	setupDatabase(its, client)

	repo := _service.NewServiceRepo(client, time.Duration(60) * time.Second)
	its.usecase = *_service.NewServiceUsecase(repo)
	its.client = client
	require.NoError(its.T(), err)
}

func (its *IntTestSuite) TearDownSuite() {
	tearDownDatabase(its)
}

func (its *IntTestSuite) TestFetch() {


	req, err := http.NewRequest(echo.GET, "/services", strings.NewReader(""))
	assert.NoError(its.T(), err)

	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)
	ctx := c.Request().Context()

	cleanTable(its, ctx)
	seedTestTable(its, ctx)

	handler := _service.NewServiceHandler(e, &its.usecase)
	err = handler.Fetch(c)
	require.NoError(its.T(), err)
	assert.Equal(its.T(), http.StatusOK, rec.Code)
}

func setupDatabase(its *IntTestSuite, client *ent.Client) {
	its.T().Log("setting up database")

	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
			its.FailNowf("failed creating schema resources:", err.Error())
	}

}

func tearDownDatabase(its *IntTestSuite) {
	its.T().Log("tearing down database")

	cErr := its.client.Close()
	if cErr != nil {
		its.FailNowf("unable to close database", cErr.Error())
	}
}

func seedTestTable(its *IntTestSuite, ctx context.Context) {
	its.T().Log("seeding test table")

	for i := 1; i <= 2; i++ {
		version, err := its.client.Version.Create().SetName("TestVersion").Save(ctx)
		if err != nil {
			its.FailNowf("unable to create version", err.Error())
		}
		// Create a new Service.
		_, nerr := its.client.Service.
			Create().
			SetTitle("TestTitle"+strconv.Itoa(i)).
			SetDescription("TestDescription").
			SetVersionCount(1).
			SetUpdatedAt(time.Now()).
			SetCreatedAt(time.Now()).
			AddVersions(version).
			Save(ctx)
			if nerr != nil {
				its.FailNowf("unable to create version", nerr.Error())
			}
	}
}

func cleanTable(its *IntTestSuite, ctx context.Context) {
	its.T().Log("cleaning database")

	_, err := its.client.Version.Delete().Exec(ctx)
	if err != nil {
		its.FailNowf("unable to clean version", err.Error())
	}

	_, serr := its.client.Service.Delete().Exec(ctx)
	if serr != nil {
		its.FailNowf("unable to clean service", err.Error())
	}
}
