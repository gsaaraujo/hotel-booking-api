package repositories_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gsaaraujo/hotel-booking-api/internal/domain/models/room"
	"github.com/gsaaraujo/hotel-booking-api/internal/infra/repositories"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type RoomsRepositorySuite struct {
	suite.Suite
	conn              *pgx.Conn
	postgresContainer testcontainers.Container
	roomsRepository   repositories.RoomsRepository
}

func (r *RoomsRepositorySuite) SetupSuite() {
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
	r.Require().NoError(err)

	r.postgresContainer = postgresContainer

	host, err := postgresContainer.Host(ctx)
	r.Require().NoError(err)

	port, err := postgresContainer.MappedPort(ctx, "5432/tcp")
	r.Require().NoError(err)

	if _, ok := os.LookupEnv("ACT"); ok {
		host = "host.docker.internal"
	}

	conn, err := pgx.Connect(context.Background(), fmt.Sprintf("postgres://postgres:postgres@%s:%s/postgres", host, port.Port()))
	r.Require().NoError(err)

	r.conn = conn
	r.roomsRepository = repositories.RoomsRepository{
		Conn: conn,
	}

	os.Setenv("PGUSER", "postgres")
	os.Setenv("PGPASSWORD", "postgres")
	os.Setenv("PGHOST", host)
	os.Setenv("PGPORT", port.Port())
	os.Setenv("PGDATABASE", "postgres")

	cmd := exec.Command("tern", "migrate", "-m", "../../../migrations")
	_, err = cmd.CombinedOutput()
	r.Require().NoError(err)
}

func (r *RoomsRepositorySuite) SetupTest() {
	ctx := context.Background()
	_, err := r.conn.Exec(ctx, "TRUNCATE TABLE rooms")
	r.Require().NoError(err)
}

func (r *RoomsRepositorySuite) TearDownSuite() {
	ctx := context.Background()

	err := r.postgresContainer.Terminate(ctx)
	r.Require().NoError(err)

	err = r.conn.Close(ctx)
	r.Require().NoError(err)
}

func (r *RoomsRepositorySuite) TestCreate_OnNoErrors_ReturnsNil() {
	type RoomSchema struct {
		Id       uuid.UUID
		Number   string
		Type     string
		Capacity uint8
		Price    uint64
	}
	roomId, err := uuid.Parse("849702fc-aad3-478f-9dd7-9963b4ca33ca")
	r.Require().NoError(err)
	newRoom := room.Room{
		Id:       roomId,
		Number:   "101",
		Type:     "SUITE",
		Price:    uint64(250),
		Capacity: uint8(2),
	}

	err = r.roomsRepository.Create(newRoom)
	r.NoError(err)

	var roomSchema RoomSchema
	err = r.conn.QueryRow(context.Background(), "SELECT id, number, type, capacity, price FROM rooms WHERE id = $1", roomId).
		Scan(&roomSchema.Id, &roomSchema.Number, &roomSchema.Type, &roomSchema.Capacity, &roomSchema.Price)
	r.NoError(err)
	r.Equal("849702fc-aad3-478f-9dd7-9963b4ca33ca", roomSchema.Id.String())
	r.Equal("101", roomSchema.Number)
	r.Equal("SUITE", roomSchema.Type)
	r.Equal(uint8(2), roomSchema.Capacity)
	r.Equal(uint64(250), roomSchema.Price)
}

func (r *RoomsRepositorySuite) TestExistsByRoomNumber_OnExists_ReturnsTrue() {
	_, err := r.conn.Exec(context.Background(), "INSERT INTO rooms (id, number, type, capacity, price) VALUES ($1, $2, $3, $4, $5)",
		"849702fc-aad3-478f-9dd7-9963b4ca33ca", "101", "SUITE", 2, 250)
	r.Require().NoError(err)

	exists, err := r.roomsRepository.ExistsByRoomNumber("101")
	r.NoError(err)

	r.True(exists)
}

func (r *RoomsRepositorySuite) TestExistsByRoomNumber_OnNotExists_ReturnsFalse() {
	exists, err := r.roomsRepository.ExistsByRoomNumber("101")
	r.NoError(err)

	r.False(exists)
}

func TestRoomsRepository(t *testing.T) {
	suite.Run(t, new(RoomsRepositorySuite))
}
