package usecases

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/gsaaraujo/hotel-booking-api/internal/application/gateways"
	"golang.org/x/crypto/bcrypt"
)

type JwtClaims struct {
	CustomerId uuid.UUID `json:"customerId"`
	jwt.RegisteredClaims
}

type LoginWithEmailAndPasswordInput struct {
	Email         string
	PlainPassword string
}

type LoginWithEmailAndPasswordOutput struct {
	CustomerId   uuid.UUID
	CustomerName string
	AccessToken  string
}

type ILoginWithEmailAndPassword interface {
	Execute(input LoginWithEmailAndPasswordInput) (LoginWithEmailAndPasswordOutput, error)
}

type LoginWithEmailAndPassword struct {
	SecretsGateway   gateways.ISecretsGateway
	CustomersGateway gateways.ICustomersGateway
}

func (l *LoginWithEmailAndPassword) Execute(input LoginWithEmailAndPasswordInput) (LoginWithEmailAndPasswordOutput, error) {
	customerDTO, err := l.CustomersGateway.FindOneByEmail(input.Email)
	if err != nil {
		return LoginWithEmailAndPasswordOutput{}, err
	}

	if customerDTO == nil {
		return LoginWithEmailAndPasswordOutput{}, errors.New("email or password is incorrect")
	}

	err = bcrypt.CompareHashAndPassword([]byte(customerDTO.HashedPassword), []byte(input.PlainPassword))
	if err != nil {
		return LoginWithEmailAndPasswordOutput{}, errors.New("email or password is incorrect")
	}

	ONE_MONTH := time.Now().Add(30 * 24 * time.Hour)

	claims := &JwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(ONE_MONTH),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtSigningAccessToken, err := l.SecretsGateway.Get("JWT_SIGNING_ACCESS_TOKEN")
	if err != nil {
		return LoginWithEmailAndPasswordOutput{}, err
	}

	signedToken, err := token.SignedString([]byte(jwtSigningAccessToken))
	if err != nil {
		return LoginWithEmailAndPasswordOutput{}, err
	}

	return LoginWithEmailAndPasswordOutput{
		CustomerId:   customerDTO.Id,
		CustomerName: customerDTO.Name,
		AccessToken:  signedToken,
	}, nil
}
