package mocks

import (
	"github.com/gsaaraujo/hotel-booking-api/internal/application/gateways"
	"github.com/stretchr/testify/mock"
)

type CustomersGatewayMock struct {
	mock.Mock
}

func (m *CustomersGatewayMock) Create(customerDTO gateways.CustomerDTO) error {
	args := m.Called(customerDTO)
	return args.Error(1)
}

func (m *CustomersGatewayMock) FindOneByEmail(email string) (*gateways.CustomerDTO, error) {
	args := m.Called(email)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*gateways.CustomerDTO), args.Error(1)
}

func (m *CustomersGatewayMock) ExistsByEmail(email string) (bool, error) {
	args := m.Called(email)
	return args.Bool(0), args.Error(1)
}
