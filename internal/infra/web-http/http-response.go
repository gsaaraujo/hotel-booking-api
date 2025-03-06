package webhttp

import "github.com/labstack/echo/v4"

type HttpResponseSuccess struct {
	StatusCode uint16 `json:"statusCode"`
	StatusText string `json:"statusText"`
	Data       any    `json:"data"`
}

type HttpResponseError struct {
	StatusCode   uint16 `json:"statusCode"`
	StatusText   string `json:"statusText"`
	ErrorMessage string `json:"error"`
}

type HttpResponseErrors struct {
	StatusCode    uint16   `json:"statusCode"`
	StatusText    string   `json:"statusText"`
	ErrorMessages []string `json:"errors"`
}

func NewOk(c echo.Context, data any) error {
	return c.JSON(200, HttpResponseSuccess{
		StatusCode: 200,
		StatusText: "OK",
		Data:       data,
	})
}

func NewCreated(c echo.Context, data any) error {
	return c.JSON(201, HttpResponseSuccess{
		StatusCode: 201,
		StatusText: "CREATED",
		Data:       data,
	})
}

func NewBadRequestValidation(c echo.Context, errorMessages []string) error {
	return c.JSON(400, HttpResponseErrors{
		StatusCode:    400,
		StatusText:    "BAD_REQUEST",
		ErrorMessages: errorMessages,
	})
}

func NewBadRequest(c echo.Context, errorMessage string) error {
	return c.JSON(400, HttpResponseError{
		StatusCode:   400,
		StatusText:   "BAD_REQUEST",
		ErrorMessage: errorMessage,
	})
}

func NewUnauthorized(c echo.Context, errorMessage string) error {
	return c.JSON(401, HttpResponseError{
		StatusCode:   401,
		StatusText:   "UNAUTHORIZED",
		ErrorMessage: errorMessage,
	})
}

func NewForbidden(c echo.Context, errorMessage string) error {
	return c.JSON(403, HttpResponseError{
		StatusCode:   403,
		StatusText:   "FORBIDDEN",
		ErrorMessage: errorMessage,
	})
}

func NewNotFound(c echo.Context, errorMessage string) error {
	return c.JSON(404, HttpResponseError{
		StatusCode:   404,
		StatusText:   "NOT_FOUND",
		ErrorMessage: errorMessage,
	})
}

func NewConflict(c echo.Context, errorMessage string) error {
	return c.JSON(409, HttpResponseError{
		StatusCode:   409,
		StatusText:   "CONFLICT",
		ErrorMessage: errorMessage,
	})
}

func NewInternalServerError(c echo.Context) error {
	return c.JSON(500, HttpResponseError{
		StatusCode:   500,
		StatusText:   "INTERNAL_SERVER_ERROR",
		ErrorMessage: "something went wrong. Please try again later",
	})
}
