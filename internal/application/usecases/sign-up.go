package usecases

import (
	"errors"
	"regexp"

	"github.com/google/uuid"
	"github.com/gsaaraujo/hotel-booking-api/internal/application/gateways"
	"golang.org/x/crypto/bcrypt"
)

type SignUpInput struct {
	Name     string
	Email    string
	Password string
}

type ISignUp interface {
	Execute(input SignUpInput) error
}

type SignUp struct {
	CustomersGateway gateways.ICustomersGateway
}

func (s *SignUp) Execute(input SignUpInput) error {
	if len(input.Name) < 3 {
		return errors.New("name must be at least 3 characters long")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	if !emailRegex.MatchString(input.Email) {
		return errors.New("email is invalid")
	}

	if len(input.Password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}

	exists, err := s.CustomersGateway.ExistsByEmail(input.Email)

	if err != nil {
		return err
	}

	if exists {
		return errors.New("email address is already associated with another account")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), 12)

	if err != nil {
		return err
	}

	err = s.CustomersGateway.Create(gateways.CustomerDTO{
		Id:             uuid.New(),
		Name:           input.Name,
		Email:          input.Email,
		HashedPassword: string(hashedPassword),
	})

	if err != nil {
		return err
	}

	return nil
}
