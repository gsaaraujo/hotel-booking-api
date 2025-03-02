package usecases_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/gsaaraujo/hotel-booking-api/internal/application/gateways"
	"github.com/gsaaraujo/hotel-booking-api/internal/application/usecases"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
)

type SignUpSuite struct {
	suite.Suite
	customersGatewayFake gateways.FakeCustomersGateway
	signUp               usecases.SignUp
}

func (s *SignUpSuite) SetupTest() {
	s.customersGatewayFake = gateways.FakeCustomersGateway{}
	s.signUp = usecases.SignUp{
		CustomersGateway: &s.customersGatewayFake,
	}
}

func (s *SignUpSuite) TestExecute_OnValidNameEmailPassword_ReturnsNil() {
	err := s.signUp.Execute(usecases.SignUpInput{
		Name:     "John Doe",
		Email:    "john.doe@gmail.com",
		Password: "123456",
	})
	s.NoError(err)

	createdCustomerDTO := s.customersGatewayFake.CustomersDTO[0]
	s.Equal("John Doe", createdCustomerDTO.Name)
	s.Equal("john.doe@gmail.com", createdCustomerDTO.Email)
	err = bcrypt.CompareHashAndPassword([]byte(createdCustomerDTO.HashedPassword), []byte("123456"))
	s.NoError(err)
}

func (s *SignUpSuite) TestExecute_OnInvalidName_ReturnsError() {
	err := s.signUp.Execute(usecases.SignUpInput{
		Name:     "J",
		Email:    "john.doe@gmail.com",
		Password: "123456",
	})

	s.EqualError(err, "name must be at least 3 characters long")
}

func (s *SignUpSuite) TestExecute_OnInvalidEmail_ReturnsError() {
	err := s.signUp.Execute(usecases.SignUpInput{
		Name:     "John Doe",
		Email:    "j.gmail.com",
		Password: "123456",
	})

	s.EqualError(err, "email is invalid")
}

func (s *SignUpSuite) TestExecute_OnInvalidPassword_ReturnsError() {
	err := s.signUp.Execute(usecases.SignUpInput{
		Name:     "John Doe",
		Email:    "john.doe@gmail.com",
		Password: "123",
	})

	s.EqualError(err, "password must be at least 6 characters long")
}

func (s *SignUpSuite) TestExecute_OnEmailAddressAlreadyAssociatedWithAnotherAccount_ReturnsError() {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("123456"), 12)
	s.Require().NoError(err)
	s.customersGatewayFake.CustomersDTO = []gateways.CustomerDTO{
		{
			Id:             uuid.New(),
			Name:           "John Doe",
			Email:          "john.doe@gmail.com",
			HashedPassword: string(hashedPassword),
		},
	}

	err = s.signUp.Execute(usecases.SignUpInput{
		Name:     "John Doe",
		Email:    "john.doe@gmail.com",
		Password: "123456",
	})

	s.EqualError(err, "email address is already associated with another account")
}

func TestSignUp(t *testing.T) {
	suite.Run(t, new(SignUpSuite))
}
