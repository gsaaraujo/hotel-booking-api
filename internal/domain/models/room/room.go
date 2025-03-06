package room

import (
	"errors"

	"github.com/google/uuid"
)

type Room struct {
	Id       uuid.UUID
	Number   string
	Type     string
	Capacity uint8
	Price    uint64
}

func NewRoom(number string, roomType string, capacity uint8, price uint64) (Room, error) {
	if !isRoomNumberValid(number) {
		return Room{}, errors.New("invalid room number format. Please enter a three-digit room number (e.g. 101) where the first digit indicates the floor number, and the last two digits represent the room number on that floor")
	}

	if roomType != "SINGLE" && roomType != "DOUBLE" && roomType != "TWIN" && roomType != "SUITE" {
		return Room{}, errors.New("room type must be SINGLE, DOUBLE, TWIN or SUITE")
	}

	if capacity <= 0 {
		return Room{}, errors.New("invalid room capacity. Please enter a capacity of at least one to accommodate guests")
	}

	if price <= 0 {
		return Room{}, errors.New("invalid room price. Please enter a value greater than zero to ensure proper pricing")
	}

	return Room{
		Id:       uuid.New(),
		Number:   number,
		Type:     roomType,
		Capacity: capacity,
		Price:    price,
	}, nil
}

func isRoomNumberValid(number string) bool {
	if len(number) != 3 {
		return false
	}

	if string(number[0]) != "1" {
		return false
	}

	return true
}
