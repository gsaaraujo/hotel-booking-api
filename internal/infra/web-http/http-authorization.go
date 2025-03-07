package webhttp

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/gsaaraujo/hotel-booking-api/internal/application/gateways"
)

type HttpAuthorization struct {
	SecretsGateway gateways.ISecretsGateway
}

func (h *HttpAuthorization) IsAdmin(authorizationToken string) bool {
	token := h.isTokenValid(authorizationToken)

	if token == nil {
		return false
	}

	claims := token.Claims.(jwt.MapClaims)
	return claims["role"] == "ADMIN"
}

func (h *HttpAuthorization) IsCustomer(authorizationToken string) bool {
	token := h.isTokenValid(authorizationToken)

	if token == nil {
		return false
	}

	claims := token.Claims.(jwt.MapClaims)
	return claims["role"] == "ADMIN" || claims["role"] == "CUSTOMER"
}

func (h *HttpAuthorization) isTokenValid(authorizationToken string) *jwt.Token {
	token, err := jwt.Parse(authorizationToken, func(token *jwt.Token) (any, error) {
		jwtSigningAccessToken, err := h.SecretsGateway.Get("JWT_SIGNING_ACCESS_TOKEN")

		if err != nil {
			return nil, err
		}

		return []byte(jwtSigningAccessToken), nil
	})

	if err != nil {
		return nil
	}

	return token
}
