package gateways

import (
	"context"

	"github.com/google/uuid"
	"github.com/gsaaraujo/hotel-booking-api/internal/application/gateways"
	"github.com/jackc/pgx/v5"
)

type CustomersGateway struct {
	Conn *pgx.Conn
}

func (c *CustomersGateway) FindOneByEmail(email string) (*gateways.CustomerDTO, error) {
	schema := struct {
		Id       uuid.UUID
		Name     string
		Password string
	}{}

	err := c.Conn.QueryRow(context.Background(), "SELECT id, name, password FROM customers WHERE email = $1", email).
		Scan(&schema.Id, &schema.Name, &schema.Password)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, nil
		}

		return &gateways.CustomerDTO{}, err
	}

	customerDTO := gateways.CustomerDTO{
		Id:             schema.Id,
		Name:           schema.Name,
		HashedPassword: schema.Password,
	}

	return &customerDTO, nil
}
