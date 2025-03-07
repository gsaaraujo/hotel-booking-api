package handlers_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gsaaraujo/hotel-booking-api/internal/application/gateways"
	"github.com/gsaaraujo/hotel-booking-api/internal/infra/handlers"
	webhttp "github.com/gsaaraujo/hotel-booking-api/internal/infra/web-http"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type GetRoomsHandlerSuite struct {
	suite.Suite
	conn               *pgx.Conn
	postgresContainer  testcontainers.Container
	fakeSecretsGateway gateways.FakeSecretsGateway
	getRoomsHandler    handlers.GetRoomsHandler
}

func (g *GetRoomsHandlerSuite) SetupSuite() {
	os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
	ctx := context.Background()
	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		Started: true,
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:17.2-alpine3.21",
			ExposedPorts: []string{"5432/tcp"},
			WaitingFor: wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(10 * time.Second),
			Env: map[string]string{
				"POSTGRES_DB":       "postgres",
				"POSTGRES_USER":     "postgres",
				"POSTGRES_PASSWORD": "postgres",
			},
		},
	})
	g.Require().NoError(err)

	g.postgresContainer = postgresContainer

	host, err := postgresContainer.Host(ctx)
	g.Require().NoError(err)

	port, err := postgresContainer.MappedPort(ctx, "5432/tcp")
	g.Require().NoError(err)

	if _, ok := os.LookupEnv("ACT"); ok {
		host = "host.docker.internal"
	}

	conn, err := pgx.Connect(context.Background(), fmt.Sprintf("postgres://postgres:postgres@%s:%s/postgres", host, port.Port()))
	g.Require().NoError(err)

	g.conn = conn
	httpLogger := webhttp.NewHttpLogger()
	g.fakeSecretsGateway = gateways.FakeSecretsGateway{}
	httpAuthorization := webhttp.HttpAuthorization{
		SecretsGateway: &g.fakeSecretsGateway,
	}
	g.getRoomsHandler = handlers.GetRoomsHandler{
		Conn:              conn,
		HttpLogger:        httpLogger,
		HttpAuthorization: httpAuthorization,
	}

	os.Setenv("PGUSER", "postgres")
	os.Setenv("PGPASSWORD", "postgres")
	os.Setenv("PGHOST", host)
	os.Setenv("PGPORT", port.Port())
	os.Setenv("PGDATABASE", "postgres")

	cmd := exec.Command("tern", "migrate", "-m", "../../../migrations")
	_, err = cmd.CombinedOutput()
	g.Require().NoError(err)
}

func (g *GetRoomsHandlerSuite) SetupTest() {
	ctx := context.Background()
	_, err := g.conn.Exec(ctx, "TRUNCATE TABLE rooms")
	g.Require().NoError(err)
}

func (g *GetRoomsHandlerSuite) TearDownSuite() {
	ctx := context.Background()

	err := g.postgresContainer.Terminate(ctx)
	g.Require().NoError(err)

	err = g.conn.Close(ctx)
	g.Require().NoError(err)
}

func (g *GetRoomsHandlerSuite) TestHandle_OnAuthorizationTokenIsMissing_ReturnsError() {
	request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`
	{
		"number": "101",
		"type": "SUITE",
		"capacity": 2,
		"price": 250
	}
`))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(request, recorder)

	err := g.getRoomsHandler.Handle(c)
	g.Require().NoError(err)

	g.Equal(401, recorder.Code)
	g.JSONEq(`
	{
		"statusCode": 401,
		"statusText": "UNAUTHORIZED",
		"error": "missing or invalid authorization token"
	}
`, recorder.Body.String())
}

func (g *GetRoomsHandlerSuite) TestHandle_OnNoPermissonToAccessResource_ReturnsError() {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role": "ANY",
	})
	signedToken, err := token.SignedString([]byte("6b45b2cb79974f989447f1d850d139f1"))
	g.Require().NoError(err)
	g.fakeSecretsGateway.Secrets = map[string]string{"JWT_SIGNING_ACCESS_TOKEN": "6b45b2cb79974f989447f1d850d139f1"}
	request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`
	{
		"number": "101",
		"type": "SUITE",
		"capacity": 2,
		"price": 250
	}
`))
	request.Header.Set("Authorization", signedToken)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(request, recorder)

	err = g.getRoomsHandler.Handle(c)
	g.Require().NoError(err)

	g.Equal(403, recorder.Code)
	g.JSONEq(`
	{
		"statusCode": 403,
		"statusText": "FORBIDDEN",
		"error": "you do not have permission to access this resource"
	}
`, recorder.Body.String())
}

func (g *GetRoomsHandlerSuite) TestHandle_OnNoErrorsAndThereAreRooms_ReturnsOk() {
	_, err := g.conn.Exec(context.Background(), "INSERT INTO rooms (id, number, type, capacity, price) VALUES ($1, $2, $3, $4, $5)",
		"849702fc-aad3-478f-9dd7-9963b4ca33ca", "101", "SUITE", 2, 250)
	g.Require().NoError(err)
	_, err = g.conn.Exec(context.Background(), "INSERT INTO rooms (id, number, type, capacity, price) VALUES ($1, $2, $3, $4, $5)",
		"57dba1c3-0421-4f24-a7c3-2a0b6c13063d", "204", "SINGLE", 8, 122)
	g.Require().NoError(err)
	_, err = g.conn.Exec(context.Background(), "INSERT INTO rooms (id, number, type, capacity, price) VALUES ($1, $2, $3, $4, $5)",
		"0dc94e80-3df8-40c9-8a79-9e9e555abbde", "132", "DOUBLE", 3, 990)
	g.Require().NoError(err)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role": "ADMIN",
	})
	signedToken, err := token.SignedString([]byte("6b45b2cb79974f989447f1d850d139f1"))
	g.Require().NoError(err)
	g.fakeSecretsGateway.Secrets = map[string]string{"JWT_SIGNING_ACCESS_TOKEN": "6b45b2cb79974f989447f1d850d139f1"}
	request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`
		{
			"number": "101",
			"type": "SUITE",
			"capacity": 2,
			"price": 250
		}
	`))
	request.Header.Set("Authorization", signedToken)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(request, recorder)

	err = g.getRoomsHandler.Handle(c)
	g.Require().NoError(err)

	g.Equal(200, recorder.Code)
	g.JSONEq(`
		{
			"statusCode": 200,
			"statusText": "OK",
			"data": [
				{
					"id": "849702fc-aad3-478f-9dd7-9963b4ca33ca",
					"type": "SUITE",
					"number": "101",
					"capacity": 2,
					"price": 250
				},
				{
					"id": "57dba1c3-0421-4f24-a7c3-2a0b6c13063d",
					"type": "SINGLE",
					"number": "204",
					"capacity": 8,
					"price": 122
				},
				{
					"id": "0dc94e80-3df8-40c9-8a79-9e9e555abbde",
					"type": "DOUBLE",
					"number": "132",
					"capacity": 3,
					"price": 990
				}
			]
		}
	`, recorder.Body.String())
}

func (g *GetRoomsHandlerSuite) TestHandle_OnNoErrorsAndThereAreNoRooms_ReturnsOk() {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role": "CUSTOMER",
	})
	signedToken, err := token.SignedString([]byte("6b45b2cb79974f989447f1d850d139f1"))
	g.Require().NoError(err)
	g.fakeSecretsGateway.Secrets = map[string]string{"JWT_SIGNING_ACCESS_TOKEN": "6b45b2cb79974f989447f1d850d139f1"}
	request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`
		{
			"number": "101",
			"type": "SUITE",
			"capacity": 2,
			"price": 250
		}
	`))
	request.Header.Set("Authorization", signedToken)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(request, recorder)

	err = g.getRoomsHandler.Handle(c)
	g.Require().NoError(err)

	g.Equal(200, recorder.Code)
	g.JSONEq(`
		{
			"statusCode": 200,
			"statusText": "OK",
			"data": []
		}
	`, recorder.Body.String())
}

func TestGetRoomsHandler(t *testing.T) {
	suite.Run(t, new(GetRoomsHandlerSuite))
}
