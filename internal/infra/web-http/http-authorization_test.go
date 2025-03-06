package webhttp_test

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gsaaraujo/hotel-booking-api/internal/application/gateways"
	webhttp "github.com/gsaaraujo/hotel-booking-api/internal/infra/web-http"
	"github.com/stretchr/testify/suite"
)

type HttpAuthorizationSuite struct {
	suite.Suite
	fakeSecretsGateway gateways.FakeSecretsGateway
	httpAuthorization  webhttp.HttpAuthorization
}

func (h *HttpAuthorizationSuite) SetupTest() {
	h.fakeSecretsGateway = gateways.FakeSecretsGateway{}
	h.httpAuthorization = webhttp.HttpAuthorization{
		SecretsGateway: &h.fakeSecretsGateway,
	}
}

func (h *HttpAuthorizationSuite) TestIsAdmin_OnValidToken_ReturnsTrue() {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role": "ADMIN",
	})
	signedToken, err := token.SignedString([]byte("6b45b2cb79974f989447f1d850d139f1"))
	h.Require().NoError(err)
	h.fakeSecretsGateway.Secrets = map[string]string{"JWT_SIGNING_ACCESS_TOKEN": "6b45b2cb79974f989447f1d850d139f1"}

	isAdmin := h.httpAuthorization.IsAdmin(signedToken)

	h.True(isAdmin)
}

func (h *HttpAuthorizationSuite) TestIsAdmin_OnRoleIsNotAdmin_ReturnsFalse() {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role": "CUSTOMER",
	})
	signedToken, err := token.SignedString([]byte("6b45b2cb79974f989447f1d850d139f1"))
	h.Require().NoError(err)
	h.fakeSecretsGateway.Secrets = map[string]string{"JWT_SIGNING_ACCESS_TOKEN": "6b45b2cb79974f989447f1d850d139f1"}

	isAdmin := h.httpAuthorization.IsAdmin(signedToken)

	h.False(isAdmin)
}

func (h *HttpAuthorizationSuite) TestIsAdmin_OnDifferentJwtSigningAccessToken_ReturnsFalse() {
	token := jwt.New(jwt.SigningMethodHS256)
	signedToken, err := token.SignedString([]byte("7e1408f7ff794cbc85403d2aaeb666c7"))
	h.Require().NoError(err)
	h.fakeSecretsGateway.Secrets = map[string]string{"JWT_SIGNING_ACCESS_TOKEN": "2f9996348690465d8980256d21bc2e43"}

	isAdmin := h.httpAuthorization.IsAdmin(signedToken)

	h.False(isAdmin)
}

func TestHttpAuthorization(t *testing.T) {
	suite.Run(t, new(HttpAuthorizationSuite))
}
