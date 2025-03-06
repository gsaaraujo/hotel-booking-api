package repositories

import "github.com/gsaaraujo/hotel-booking-api/internal/domain/models/room"

type FakeRoomsRepository struct {
	Rooms []room.Room
}

func (f *FakeRoomsRepository) Create(room room.Room) error {
	f.Rooms = append(f.Rooms, room)
	return nil
}

func (f *FakeRoomsRepository) ExistsByRoomNumber(roomNumber string) (bool, error) {
	for _, room := range f.Rooms {
		if room.Number == roomNumber {
			return true, nil
		}
	}

	return false, nil
}
