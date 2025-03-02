package handlers

import (
	"github.com/gsaaraujo/hotel-booking-api/internal/application/usecases"
	webhttp "github.com/gsaaraujo/hotel-booking-api/internal/infra/web-http"
	"github.com/labstack/echo/v4"
)

type SignUpHandlerInput struct {
	Name     any `validate:"required,string,notEmpty,lt=256"`
	Email    any `validate:"required,string,notEmpty,lt=256"`
	Password any `validate:"required,string,notEmpty,lt=256"`
}

type SignUpHandler struct {
	HttpLogger webhttp.HttpLogger
	SignUp     usecases.ISignUp
}

func (s *SignUpHandler) Handle(c echo.Context) error {
	var requestBody SignUpHandlerInput

	if err := c.Bind(&requestBody); err != nil {
		return webhttp.NewBadRequestValidation(c, []string{"content-type must be application/json"})
	}

	errorMessages := []string{}
	validator, err := webhttp.NewHttpValidator()

	if err != nil {
		s.HttpLogger.Log(c, err)
		return webhttp.NewInternalServerError(c)
	}

	errorMessages = append(errorMessages, validator.Validate(requestBody)...)

	if len(errorMessages) > 0 {
		return webhttp.NewBadRequestValidation(c, errorMessages)
	}

	err = s.SignUp.Execute(usecases.SignUpInput{
		Name:     requestBody.Name.(string),
		Email:    requestBody.Email.(string),
		Password: requestBody.Password.(string),
	})

	if err != nil {
		if err.Error() == "name must be at least 3 characters long" {
			return webhttp.NewBadRequest(c, err.Error())
		}

		if err.Error() == "password must be at least 6 characters long" {
			return webhttp.NewBadRequest(c, err.Error())
		}

		if err.Error() == "email is invalid" {
			return webhttp.NewBadRequest(c, "email address is invalid. Please enter a valid email address")
		}

		if err.Error() == "email address is already associated with another account" {
			return webhttp.NewConflict(c, "this email address is already in use. Please use a different email or login to your existing account")
		}

		s.HttpLogger.Log(c, err)
		return webhttp.NewInternalServerError(c)
	}

	return webhttp.NewCreated(c, "customer sign up successfully")

}
