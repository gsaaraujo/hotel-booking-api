package room_test

import (
	"testing"

	"github.com/gsaaraujo/hotel-booking-api/internal/domain/models/room"
	"github.com/stretchr/testify/suite"
)

type RoomSuite struct {
	suite.Suite
}

func (r *RoomSuite) TestNewRoom_OnNoErrors_ReturnsRoom() {
	newRoom, err := room.NewRoom("101", "SINGLE", 2, 250)
	r.NoError(err)

	r.Equal("101", newRoom.Number)
	r.Equal("SINGLE", newRoom.Type)
	r.Equal(uint8(2), newRoom.Capacity)
	r.Equal(uint64(250), newRoom.Price)
}

func (r *RoomSuite) TestNewRoom_OnInvalidNumber_ReturnsError() {
	roomNumbers := []string{"", " ", "0", "01", "10", "000", "001", "0000", "0001", "1010", "abc"}

	for _, roomNumber := range roomNumbers {
		_, err := room.NewRoom(roomNumber, "SINGLE", 2, 250)
		r.EqualError(err, "invalid room number format. Please enter a three-digit room number (e.g. 101) where the first digit indicates the floor number, and the last two digits represent the room number on that floor")
	}
}

func (r *RoomSuite) TestNewRoom_OnInvalidType_ReturnsError() {
	_, err := room.NewRoom("101", "", 2, 250)

	r.EqualError(err, "room type must be SINGLE, DOUBLE, TWIN or SUITE")
}

func (r *RoomSuite) TestNewRoom_OnInvalidCapacity_ReturnsError() {
	_, err := room.NewRoom("101", "SINGLE", 0, 250)

	r.EqualError(err, "invalid room capacity. Please enter a capacity of at least one to accommodate guests")
}

func (r *RoomSuite) TestNewRoom_OnInvalidPrice_ReturnsError() {
	_, err := room.NewRoom("101", "SINGLE", 2, 0)

	r.EqualError(err, "invalid room price. Please enter a value greater than zero to ensure proper pricing")
}

func TestRoom(t *testing.T) {
	suite.Run(t, new(RoomSuite))
}
