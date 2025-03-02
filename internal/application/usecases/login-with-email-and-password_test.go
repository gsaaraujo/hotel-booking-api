package usecases_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/gsaaraujo/hotel-booking-api/internal/application/gateways"
	"github.com/gsaaraujo/hotel-booking-api/internal/application/usecases"
	"github.com/gsaaraujo/hotel-booking-api/internal/application/usecases/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
)

type SecretsGatewayMock struct {
	mock.Mock
}

func (m *SecretsGatewayMock) Get(key string) (string, error) {
	args := m.Called(key)
	return args.String(0), args.Error(1)
}

type LoginWithEmailAndPasswordSuite struct {
	suite.Suite
	secretsGatewayMock        SecretsGatewayMock
	customersGatewayMock      mocks.CustomersGatewayMock
	loginWithEmailAndPassword usecases.LoginWithEmailAndPassword
}

func (l *LoginWithEmailAndPasswordSuite) SetupTest() {
	l.secretsGatewayMock = SecretsGatewayMock{}
	l.customersGatewayMock = mocks.CustomersGatewayMock{}
	l.loginWithEmailAndPassword = usecases.LoginWithEmailAndPassword{
		SecretsGateway:   &l.secretsGatewayMock,
		CustomersGateway: &l.customersGatewayMock,
	}
}

func (l *LoginWithEmailAndPasswordSuite) TestExecute_OnCorrectEmailAndPassword_ReturnsOutput() {
	customerId, err := uuid.Parse("aa473b65-90a8-48ad-ab7d-5bd50a806d38")
	customerName := "John Doe"
	l.Require().NoError(err)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("123456"), 12)
	l.Require().NoError(err)
	customerDTO := gateways.CustomerDTO{
		Id:             customerId,
		Name:           customerName,
		HashedPassword: string(hashedPassword),
	}
	l.secretsGatewayMock.On("Get", "JWT_SIGNING_ACCESS_TOKEN").Return("secret", nil)
	l.customersGatewayMock.On("FindOneByEmail", "john.doe@gmail.com").Return(&customerDTO, nil)
	input := usecases.LoginWithEmailAndPasswordInput{
		Email:         "john.doe@gmail.com",
		PlainPassword: "123456",
	}

	output, err := l.loginWithEmailAndPassword.Execute(input)
	l.Require().NoError(err)

	l.Require().NoError(err)
	l.Equal(customerId, output.CustomerId)
	l.Equal(customerName, output.CustomerName)
}

func (l *LoginWithEmailAndPasswordSuite) TestExecute_OnCorrectEmailButIncorrectPassword_ReturnsError() {
	customerId, err := uuid.Parse("aa473b65-90a8-48ad-ab7d-5bd50a806d38")
	l.Require().NoError(err)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("123456"), 12)
	l.Require().NoError(err)
	customerDTO := gateways.CustomerDTO{
		Id:             customerId,
		HashedPassword: string(hashedPassword),
	}
	l.customersGatewayMock.On("FindOneByEmail", "john.doe@gmail.com").Return(&customerDTO, nil)
	input := usecases.LoginWithEmailAndPasswordInput{
		Email:         "john.doe@gmail.com",
		PlainPassword: "abcdef",
	}

	_, err = l.loginWithEmailAndPassword.Execute(input)

	l.EqualError(err, "email or password is incorrect")
}

func (l *LoginWithEmailAndPasswordSuite) TestExecute_OnUnregisteredEmail_ReturnsError() {
	l.customersGatewayMock.On("FindOneByEmail", "john.doe@gmail.com").Return(nil, nil)
	input := usecases.LoginWithEmailAndPasswordInput{
		Email:         "john.doe@gmail.com",
		PlainPassword: "123456",
	}

	_, err := l.loginWithEmailAndPassword.Execute(input)

	l.EqualError(err, "email or password is incorrect")
}

func TestLoginWithEmailAndPassword(t *testing.T) {
	suite.Run(t, new(LoginWithEmailAndPasswordSuite))
}
