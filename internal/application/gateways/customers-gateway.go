package gateways

import "github.com/google/uuid"

type CustomerDTO struct {
	Id             uuid.UUID
	Name           string
	Email          string
	HashedPassword string
}

type ICustomersGateway interface {
	Create(customerDTO CustomerDTO) error
	FindOneByEmail(email string) (*CustomerDTO, error)
	ExistsByEmail(email string) (bool, error)
}
