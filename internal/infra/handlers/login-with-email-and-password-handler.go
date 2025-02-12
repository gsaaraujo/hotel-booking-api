package handlers

import (
	"reflect"
	"strings"

	"github.com/google/uuid"
	"github.com/gsaaraujo/hotel-booking-api/internal/application/usecases"
	webhttp "github.com/gsaaraujo/hotel-booking-api/internal/infra/web-http"
	"github.com/labstack/echo/v4"
)

type LoginWithEmailAndPasswordHandlerInput struct {
	Email    *interface{}
	Password *interface{}
}

type LoginWithEmailAndPasswordHandlerOutput struct {
	CustomerId   uuid.UUID `json:"customerId"`
	CustomerName string    `json:"customerName"`
	AccessToken  string    `json:"accessToken"`
}

type LoginWithEmailAndPasswordHandler struct {
	HttpLogger                webhttp.HttpLogger
	LoginWithEmailAndPassword usecases.ILoginWithEmailAndPassword
}

func (l *LoginWithEmailAndPasswordHandler) Handle(c echo.Context) error {
	var requestBody LoginWithEmailAndPasswordHandlerInput

	if err := c.Bind(&requestBody); err != nil {
		return webhttp.NewBadRequestValidation(c, []string{"content-type must be application/json"})
	}

	errorMessages := []string{}

	if requestBody.Email == nil {
		errorMessages = append(errorMessages, "email is required")
	}

	if requestBody.Password == nil {
		errorMessages = append(errorMessages, "password is required")
	}

	if requestBody.Email != nil && reflect.TypeOf(*requestBody.Email).Kind() == reflect.String && strings.TrimSpace((*requestBody.Email).(string)) == "" {
		errorMessages = append(errorMessages, "email must not be empty")
	}

	if requestBody.Password != nil && reflect.TypeOf(*requestBody.Password).Kind() == reflect.String && strings.TrimSpace((*requestBody.Password).(string)) == "" {
		errorMessages = append(errorMessages, "password must not be empty")
	}

	if requestBody.Email != nil && reflect.TypeOf(*requestBody.Email).Kind() != reflect.String {
		errorMessages = append(errorMessages, "email must be string")
	}

	if requestBody.Password != nil && reflect.TypeOf(*requestBody.Password).Kind() != reflect.String {
		errorMessages = append(errorMessages, "password must be string")
	}

	if len(errorMessages) > 0 {
		return webhttp.NewBadRequestValidation(c, errorMessages)
	}

	input := usecases.LoginWithEmailAndPasswordInput{
		Email:         (*requestBody.Email).(string),
		PlainPassword: (*requestBody.Password).(string),
	}

	output, err := l.LoginWithEmailAndPassword.Execute(input)

	if err != nil {
		if err.Error() == "email or password is incorrect" {
			return webhttp.NewUnauthorized(c, err.Error())
		}

		l.HttpLogger.Log(c, err)
		return webhttp.NewInternalServerError(c)
	}

	requestOutput := LoginWithEmailAndPasswordHandlerOutput{
		CustomerId:   output.CustomerId,
		CustomerName: output.CustomerName,
		AccessToken:  output.AccessToken,
	}

	return webhttp.NewOk(c, requestOutput)
}
