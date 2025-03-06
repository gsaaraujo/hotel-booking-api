package usecases

import (
	"fmt"

	"github.com/gsaaraujo/hotel-booking-api/internal/application/repositories"
	"github.com/gsaaraujo/hotel-booking-api/internal/domain/models/room"
)

type CreateRoomInput struct {
	Number   string
	Type     string
	Capacity uint8
	Price    uint64
}

type CreateRoomOutput struct{}

type ICreateRoom interface {
	Execute(input CreateRoomInput) error
}

type CreateRoom struct {
	RoomsRepository repositories.IRoomsRepository
}

func (c *CreateRoom) Execute(input CreateRoomInput) error {
	exists, err := c.RoomsRepository.ExistsByRoomNumber(input.Number)

	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("the room number '%s' is already in use. Please assign another room number", input.Number)
	}

	newRoom, err := room.NewRoom(input.Number, input.Type, input.Capacity, input.Price)

	if err != nil {
		return err
	}

	err = c.RoomsRepository.Create(newRoom)

	if err != nil {
		return err
	}

	return nil
}
