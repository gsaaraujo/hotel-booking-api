package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/gsaaraujo/hotel-booking-api/internal/domain/models/room"
	"github.com/jackc/pgx/v5"
)

type RoomsRepository struct {
	Conn *pgx.Conn
}

func (r *RoomsRepository) Create(room room.Room) error {
	_, err := r.Conn.Exec(context.Background(), "INSERT INTO rooms (id, number, type, capacity, price) VALUES ($1, $2, $3, $4, $5)",
		room.Id.String(), room.Number, room.Type, room.Capacity, room.Price)

	if err != nil {
		return err
	}

	return nil
}

func (r *RoomsRepository) ExistsByRoomNumber(roomNumber string) (bool, error) {
	var roomId uuid.UUID
	err := r.Conn.QueryRow(context.Background(), "SELECT id FROM rooms WHERE number = $1", roomNumber).Scan(&roomId)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
