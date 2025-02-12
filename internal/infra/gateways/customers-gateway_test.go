package gateways_test

import (
	"context"
	"fmt"
	"os"
	"testing"

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

func (c *CustomersGatewaySuite) SetupTest() {
	os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
	ctx := context.Background()
	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		Started: true,
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:17.2-alpine3.21",
			ExposedPorts: []string{"5432/tcp"},
			WaitingFor:   wait.ForListeningPort("5432/tcp"),
			Env: map[string]string{
				"POSTGRES_DB":       "postgres",
				"POSTGRES_USER":     "postgres",
				"POSTGRES_PASSWORD": "postgres",
			},
		},
	})
	c.Require().NoError(err)

	host, err := postgresContainer.Host(ctx)
	c.Require().NoError(err)

	port, err := postgresContainer.MappedPort(ctx, "5432/tcp")
	c.Require().NoError(err)

	conn, err := pgx.Connect(context.Background(), fmt.Sprintf("postgres://postgres:postgres@%s:%s/postgres", host, port.Port()))
	c.Require().NoError(err)

	conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS customers (
			id UUID PRIMARY KEY,
			name VARCHAR(50) NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			password TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)

	c.conn = conn
	c.postgresContainer = postgresContainer
	c.customersGateway = gateways.CustomersGateway{
		Conn: conn,
	}
}

func (c *CustomersGatewaySuite) TearDownTest() {
	ctx := context.Background()
	c.postgresContainer.Terminate(ctx)
	c.conn.Close(ctx)
}

func (c *CustomersGatewaySuite) TestFindOneByEmail_OnFinding_ReturnsCustomer() {
	ctx := context.Background()
	c.conn.Exec(ctx, `INSERT INTO customers (id, name, email, password) 
		VALUES ('620d8a0f-abc2-4f80-a1bc-407a037bd920', 'John Doe', 'john.doe@gmail.com', '$2a$12$zkX5/W4LHciSZLR4YRLxHetVwAdppboUHJ6JnNhfSrKqVaSJk5hzu')`)

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
