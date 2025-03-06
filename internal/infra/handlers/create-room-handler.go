package handlers

import (
	"fmt"

	"github.com/gsaaraujo/hotel-booking-api/internal/application/usecases"
	webhttp "github.com/gsaaraujo/hotel-booking-api/internal/infra/web-http"
	"github.com/labstack/echo/v4"
)

type CreateRoomHandlerInput struct {
	Number   any `validate:"required,string,notEmpty,lt=256"`
	Type     any `validate:"required,string,notEmpty,lt=256"`
	Capacity any `validate:"required,integer,positive,lt=1000"`
	Price    any `validate:"required,integer,positive,lt=1000000000"`
}

type CreateRoomHandler struct {
	HttpLogger        webhttp.HttpLogger
	HttpAuthorization webhttp.HttpAuthorization
	HttpValidator     webhttp.HttpValidator
	CreateRoom        usecases.ICreateRoom
}

func (cr *CreateRoomHandler) Handle(c echo.Context) error {
	authorizationToken := c.Request().Header.Get("Authorization")

	if authorizationToken == "" {
		return webhttp.NewUnauthorized(c, "missing or invalid authorization token")
	}

	if !cr.HttpAuthorization.IsAdmin(authorizationToken) {
		return webhttp.NewForbidden(c, "you do not have permission to access this resource")
	}

	var input CreateRoomHandlerInput

	if err := c.Bind(&input); err != nil {
		return webhttp.NewBadRequestValidation(c, []string{"content-type must be application/json"})
	}

	if len(cr.HttpValidator.Validate(input)) > 0 {
		return webhttp.NewBadRequestValidation(c, cr.HttpValidator.Validate(input))
	}

	err := cr.CreateRoom.Execute(usecases.CreateRoomInput{
		Number:   input.Number.(string),
		Type:     input.Type.(string),
		Capacity: uint8(input.Capacity.(float64)),
		Price:    uint64(input.Price.(float64)),
	})

	if err != nil {
		if err.Error() == fmt.Sprintf("the room number '%s' is already in use. Please assign another room number", input.Number) {
			return webhttp.NewConflict(c, err.Error())
		}

		if err.Error() == "invalid room number format. Please enter a three-digit room number (e.g. 101) where the first digit indicates the floor number, and the last two digits represent the room number on that floor" {
			return webhttp.NewConflict(c, err.Error())
		}

		if err.Error() == "room type must be SINGLE, DOUBLE, TWIN or SUITE" {
			return webhttp.NewConflict(c, err.Error())
		}

		if err.Error() == "invalid room capacity. Please enter a capacity of at least one to accommodate guests" {
			return webhttp.NewConflict(c, err.Error())
		}

		if err.Error() == "invalid room price. Please enter a value greater than zero to ensure proper pricing" {
			return webhttp.NewConflict(c, err.Error())
		}

		cr.HttpLogger.Log(c, err)
		return webhttp.NewInternalServerError(c)
	}

	return webhttp.NewCreated(c, nil)
}
