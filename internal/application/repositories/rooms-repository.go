package repositories

import "github.com/gsaaraujo/hotel-booking-api/internal/domain/models/room"

type IRoomsRepository interface {
	Create(room room.Room) error
	ExistsByRoomNumber(roomNumber string) (bool, error)
}
