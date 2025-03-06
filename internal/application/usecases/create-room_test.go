package usecases_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/gsaaraujo/hotel-booking-api/internal/application/repositories"
	"github.com/gsaaraujo/hotel-booking-api/internal/application/usecases"
	"github.com/gsaaraujo/hotel-booking-api/internal/domain/models/room"
	"github.com/stretchr/testify/suite"
)

type CreateRoomSuite struct {
	suite.Suite
	createRoom          usecases.CreateRoom
	fakeRoomsRepository repositories.FakeRoomsRepository
}

func (c *CreateRoomSuite) SetupTest() {
	c.fakeRoomsRepository = repositories.FakeRoomsRepository{}
	c.createRoom = usecases.CreateRoom{
		RoomsRepository: &c.fakeRoomsRepository,
	}
}

func (c *CreateRoomSuite) TestExecute_OnNoErrors_ReturnsNil() {
	err := c.createRoom.Execute(usecases.CreateRoomInput{
		Number:   "101",
		Type:     "SUITE",
		Capacity: uint8(2),
		Price:    uint64(250),
	})
	c.NoError(err)

	createdRoom := c.fakeRoomsRepository.Rooms[0]
	c.Equal("101", createdRoom.Number)
	c.Equal("SUITE", createdRoom.Type)
	c.Equal(uint8(2), createdRoom.Capacity)
	c.Equal(uint64(250), createdRoom.Price)
}

func (c *CreateRoomSuite) TestExecute_OnDuplicateRoomNumber_ReturnsError() {
	c.fakeRoomsRepository.Rooms = []room.Room{
		{
			Id:       uuid.New(),
			Number:   "101",
			Type:     "SUITE",
			Price:    uint64(250),
			Capacity: uint8(2),
		},
	}

	err := c.createRoom.Execute(usecases.CreateRoomInput{
		Number:   "101",
		Type:     "SUITE",
		Capacity: uint8(2),
		Price:    uint64(250),
	})

	c.EqualError(err, "the room number '101' is already in use. Please assign another room number")
}

func TestCreateRoom(t *testing.T) {
	suite.Run(t, new(CreateRoomSuite))
}
