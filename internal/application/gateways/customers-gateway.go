package gateways

import "github.com/google/uuid"

type CustomerDTO struct {
	Id             uuid.UUID
	Name           string
	HashedPassword string
}

type ICustomersGateway interface {
	FindOneByEmail(email string) (*CustomerDTO, error)
}
