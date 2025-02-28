package gateways_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/gsaaraujo/hotel-booking-api/internal/infra/gateways"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type CustomersGatewaySuite struct {
	suite.Suite
	conn              *pgx.Conn
	postgresContainer testcontainers.Container
	customersGateway  gateways.CustomersGateway
}

func (c *CustomersGatewaySuite) SetupSuite() {
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
	c.Require().NoError(err)

	c.postgresContainer = postgresContainer

	host, err := postgresContainer.Host(ctx)
	c.Require().NoError(err)

	port, err := postgresContainer.MappedPort(ctx, "5432/tcp")
	c.Require().NoError(err)

	if _, ok := os.LookupEnv("ACT"); ok {
		host = "host.docker.internal"
	}

	conn, err := pgx.Connect(context.Background(), fmt.Sprintf("postgres://postgres:postgres@%s:%s/postgres", host, port.Port()))
	c.Require().NoError(err)

	c.conn = conn
	c.customersGateway = gateways.CustomersGateway{
		Conn: conn,
	}

	os.Setenv("PGUSER", "postgres")
	os.Setenv("PGPASSWORD", "postgres")
	os.Setenv("PGHOST", host)
	os.Setenv("PGPORT", port.Port())
	os.Setenv("PGDATABASE", "postgres")

	cmd := exec.Command("tern", "migrate", "-m", "../../../migrations")
	_, err = cmd.CombinedOutput()
	c.Require().NoError(err)
}

func (c *CustomersGatewaySuite) SetupTest() {
	ctx := context.Background()
	_, err := c.conn.Exec(ctx, "TRUNCATE TABLE customers")
	c.Require().NoError(err)
}

func (c *CustomersGatewaySuite) TearDownSuite() {
	ctx := context.Background()

	err := c.postgresContainer.Terminate(ctx)
	c.Require().NoError(err)

	err = c.conn.Close(ctx)
	c.Require().NoError(err)
}

func (c *CustomersGatewaySuite) TestFindOneByEmail_OnFinding_ReturnsCustomer() {
	ctx := context.Background()
	_, err := c.conn.Exec(ctx, `INSERT INTO customers (id, name, email, password) 
		VALUES ('620d8a0f-abc2-4f80-a1bc-407a037bd920', 'John Doe', 'john.doe@gmail.com', '$2a$12$zkX5/W4LHciSZLR4YRLxHetVwAdppboUHJ6JnNhfSrKqVaSJk5hzu')`)
	c.Require().NoError(err)

	customerDTO, err := c.customersGateway.FindOneByEmail("john.doe@gmail.com")

	c.Require().NoError(err)
	c.Equal("620d8a0f-abc2-4f80-a1bc-407a037bd920", customerDTO.Id.String())
	c.Equal("John Doe", customerDTO.Name)
	c.Equal("$2a$12$zkX5/W4LHciSZLR4YRLxHetVwAdppboUHJ6JnNhfSrKqVaSJk5hzu", customerDTO.HashedPassword)
}

func (c *CustomersGatewaySuite) TestFindOneByEmail_OnNoFinding_ReturnsNil() {
	customerDTO, err := c.customersGateway.FindOneByEmail("john.doe@gmail.com")

	c.Require().NoError(err)
	c.Nil(customerDTO)
}

func TestCustomersGateway(t *testing.T) {
	suite.Run(t, new(CustomersGatewaySuite))
}
