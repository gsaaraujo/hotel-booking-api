package gateways_test

import (
	"os"
	"testing"

	"github.com/gsaaraujo/hotel-booking-api/internal/infra/gateways"
	"github.com/stretchr/testify/suite"
)

type LocalSecretsGatewaySuite struct {
	suite.Suite
	localSecretsGateway gateways.LocalSecretsGateway
}

func (l *LocalSecretsGatewaySuite) SetupSuite() {
	l.localSecretsGateway = gateways.LocalSecretsGateway{
		PathToFile: ".env.test",
	}

}

func (l *LocalSecretsGatewaySuite) TearDownTest() {
	os.Remove(".env.test")
}

func (l *LocalSecretsGatewaySuite) TestGet_OnSecretFound1_ReturnsError() {
	err := os.WriteFile(".env.test", []byte(`
		SECRET1=any_secret_1
		POSTGRES_URL=postgres_url
		SECRET2=any_secret_2
	`), 0644)
	l.Require().NoError(err)

	postgresUrl, err := l.localSecretsGateway.Get("POSTGRES_URL")
	l.Require().NoError(err)

	l.Equal("postgres_url", postgresUrl)
}

func (l *LocalSecretsGatewaySuite) TestGet_OnSecretFound2_ReturnsError() {
	err := os.WriteFile(".env.test", []byte(`
SECRET1=any_secret_1
POSTGRES_URL=postgres_url
SECRET2=any_secret_2
	`), 0644)
	l.Require().NoError(err)

	postgresUrl, err := l.localSecretsGateway.Get("POSTGRES_URL")
	l.Require().NoError(err)

	l.Equal("postgres_url", postgresUrl)
}

func (l *LocalSecretsGatewaySuite) TestGet_OnSecretFound3_ReturnsError() {
	err := os.WriteFile(".env.test", []byte("POSTGRES_URL=postgres_url"), 0644)
	l.Require().NoError(err)

	postgresUrl, err := l.localSecretsGateway.Get("POSTGRES_URL")
	l.Require().NoError(err)

	l.Equal("postgres_url", postgresUrl)
}

func (l *LocalSecretsGatewaySuite) TestGet_OnSecretNotFound_ReturnsError() {
	err := os.WriteFile(".env.test", []byte(""), 0644)
	l.Require().NoError(err)

	_, err = l.localSecretsGateway.Get("POSTGRES_URL")

	l.EqualError(err, "secret POSTGRES_URL not found")
}

func TestLocalSecretsGateway(t *testing.T) {
	suite.Run(t, new(LocalSecretsGatewaySuite))
}
