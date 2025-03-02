package handlers_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gsaaraujo/hotel-booking-api/internal/application/usecases"
	"github.com/gsaaraujo/hotel-booking-api/internal/infra/handlers"
	webhttp "github.com/gsaaraujo/hotel-booking-api/internal/infra/web-http"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type SignUpMock struct {
	mock.Mock
}

func (s *SignUpMock) Execute(input usecases.SignUpInput) error {
	args := s.Called(input)
	return args.Error(0)
}

type SignUpHandlerSuite struct {
	suite.Suite
	signUpMock    SignUpMock
	signUpHandler handlers.SignUpHandler
}

func (s *SignUpHandlerSuite) SetupTest() {
	httpLogger := webhttp.NewHttpLogger()
	s.signUpMock = SignUpMock{}
	s.signUpHandler = handlers.SignUpHandler{
		SignUp:     &s.signUpMock,
		HttpLogger: httpLogger,
	}
}

func (s *SignUpHandlerSuite) TestHandle_OnNoErrors_ReturnsCreated() {
	s.signUpMock.On("Execute", usecases.SignUpInput{
		Name:     "John Doe",
		Email:    "john.doe@gmail.com",
		Password: "123456",
	}).Return(nil)
	request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`
		{
			"name": "John Doe",
			"email": "john.doe@gmail.com",
			"password": "123456"
		}
	`))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(request, recorder)

	err := s.signUpHandler.Handle(c)
	s.Require().NoError(err)

	s.Require().NoError(err)
	s.Equal(201, recorder.Code)
	s.JSONEq(`
		{
			"statusCode": 201,
			"statusText": "CREATED",
			"data": "customer sign up successfully"
		}
	`, recorder.Body.String())
}

func (s *SignUpHandlerSuite) TestHandle_OnInvalidName_ReturnsBadRequest() {
	s.signUpMock.On("Execute", usecases.SignUpInput{
		Name:     "J",
		Email:    "john.doe@gmail.com",
		Password: "123456",
	}).Return(errors.New("name must be at least 3 characters long"))
	request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`
		{
			"name": "J",
			"email": "john.doe@gmail.com",
			"password": "123456"
		}
	`))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(request, recorder)

	err := s.signUpHandler.Handle(c)
	s.Require().NoError(err)

	s.Equal(400, recorder.Code)
	s.JSONEq(`
		{
			"statusCode": 400,
			"statusText": "BAD_REQUEST",
			"error": "name must be at least 3 characters long"
		}
	`, recorder.Body.String())
}

func (s *SignUpHandlerSuite) TestHandle_OnInvalidEmail_ReturnsBadRequest() {
	s.signUpMock.On("Execute", usecases.SignUpInput{
		Name:     "John Doe",
		Email:    "jjgmail.com",
		Password: "123456",
	}).Return(errors.New("email is invalid"))
	request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`
		{
			"name": "John Doe",
			"email": "jjgmail.com",
			"password": "123456"
		}
	`))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(request, recorder)

	err := s.signUpHandler.Handle(c)
	s.Require().NoError(err)

	s.Equal(400, recorder.Code)
	s.JSONEq(`
		{
			"statusCode": 400,
			"statusText": "BAD_REQUEST",
			"error": "email address is invalid. Please enter a valid email address"
		}
	`, recorder.Body.String())
}

func (s *SignUpHandlerSuite) TestHandle_OnInvalidPassword_ReturnsBadRequest() {
	s.signUpMock.On("Execute", usecases.SignUpInput{
		Name:     "John Doe",
		Email:    "john.doe@gmail.com",
		Password: "123",
	}).Return(errors.New("password must be at least 6 characters long"))
	request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`
		{
			"name": "John Doe",
			"email": "john.doe@gmail.com",
			"password": "123"
		}
	`))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(request, recorder)

	err := s.signUpHandler.Handle(c)
	s.Require().NoError(err)

	s.Equal(400, recorder.Code)
	s.JSONEq(`
		{
			"statusCode": 400,
			"statusText": "BAD_REQUEST",
			"error": "password must be at least 6 characters long"
		}
	`, recorder.Body.String())
}

func (s *SignUpHandlerSuite) TestHandle_OnEmailAddressAlreadyAssociatedWithAnotherAccount_ReturnsConflict() {
	s.signUpMock.On("Execute", usecases.SignUpInput{
		Name:     "John Doe",
		Email:    "john.doe@gmail.com",
		Password: "123456",
	}).Return(errors.New("email address is already associated with another account"))
	request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`
		{
			"name": "John Doe",
			"email": "john.doe@gmail.com",
			"password": "123456"
		}
	`))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(request, recorder)

	err := s.signUpHandler.Handle(c)
	s.Require().NoError(err)

	s.Equal(409, recorder.Code)
	s.JSONEq(`
		{
			"statusCode": 409,
			"statusText": "CONFLICT",
			"error": "this email address is already in use. Please use a different email or login to your existing account"
		}
	`, recorder.Body.String())
}

func (s *SignUpHandlerSuite) TestHandle_OnAnyUnexpectedError_ReturnsInternalServerError() {
	s.signUpMock.On("Execute", usecases.SignUpInput{
		Name:     "John Doe",
		Email:    "john.doe@gmail.com",
		Password: "123",
	}).Return(errors.New("any unexpected error"))
	request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`
		{
			"name": "John Doe",
			"email": "john.doe@gmail.com",
			"password": "123"
		}
	`))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(request, recorder)

	err := s.signUpHandler.Handle(c)
	s.Require().NoError(err)

	s.Equal(500, recorder.Code)
	s.JSONEq(`
		{
			"statusCode": 500,
			"statusText": "INTERNAL_SERVER_ERROR",
			"error": "something went wrong. Please try again later"
		}
	`, recorder.Body.String())
}

func (s *SignUpHandlerSuite) TestHandle_OnInvalidBody_ReturnsBadRequest() {
	bodiesAndErrors := []map[string]string{
		{
			"body":   `abc`,
			"errors": `["content-type must be application/json"]`,
		},
		{
			"body":   `{}`,
			"errors": `["name is required", "email is required", "password is required"]`,
		},
		{
			"body": `{
				"name": null,
				"email": null,
				"password": null
			}`,
			"errors": `["name is required", "email is required", "password is required"]`,
		},
		{
			"body": `{
				"name": "",
				"email": "",
				"password": ""
			}`,
			"errors": `["name must not be empty", "email must not be empty", "password must not be empty"]`,
		},
		{
			"body": `{
				"name": " ",
				"email": " ",
				"password": " "
			}`,
			"errors": `["name must not be empty", "email must not be empty", "password must not be empty"]`,
		},
		{
			"body": `{
				"name": 1,
				"email": 1,
				"password": 2
			}`,
			"errors": `["name must be string", "email must be string", "password must be string"]`,
		},
		{
			"body": `{
				"name": "",
				"email": " ",
				"password": -2
			}`,
			"errors": `["name must not be empty", "email must not be empty", "password must be string"]`,
		},
		{
			"body": `{
				"name": [],
				"email": 0.5,
				"password": {}
			}`,
			"errors": `["name must be string", "email must be string", "password must be string"]`,
		},
		{
			"body": `{
				"name": "joieoplkkuyrmmnchskkloaokaweokfopsdjgsdijjoieoplkkuyrmmnchskkloaokaweokfopsdjgsdijjoieoplkkuyrmmnchskkloaokaweokfopsdjgsdijjoieoplkkuyrmmnchskkloaokaweokfopsdjgsdijjoieoplkkuyrmmnchskkloaokaweokfopsdjgsdijjoieoplkkuyrmmnchskkloaokaweokfopsdjgsdijjoieoplkkuyuio",
				"email": "joieoplkkuyrmmnchskkloaokaweokfopsdjgsdijjoieoplkkuyrmmnchskkloaokaweokfopsdjgsdijjoieoplkkuyrmmnchskkloaokaweokfopsdjgsdijjoieoplkkuyrmmnchskkloaokaweokfopsdjgsdijjoieoplkkuyrmmnchskkloaokaweokfopsdjgsdijjoieoplkkuyrmmnchskkloaokaweokfopsdjgsdijjoieoplkkuyuio",
				"password": "joieoplkkuyrmmnchskkloaokaweokfopsdjgsdijjoieoplkkuyrmmnchskkloaokaweokfopsdjgsdijjoieoplkkuyrmmnchskkloaokaweokfopsdjgsdijjoieoplkkuyrmmnchskkloaokaweokfopsdjgsdijjoieoplkkuyrmmnchskkloaokaweokfopsdjgsdijjoieoplkkuyrmmnchskkloaokaweokfopsdjgsdijjoieoplkkuyuio"
			}`,
			"errors": `["name must be less than 256", "email must be less than 256", "password must be less than 256"]`,
		},
	}

	for _, inputAndError := range bodiesAndErrors {
		body := inputAndError["body"]
		errorMessage := inputAndError["errors"]
		request := httptest.NewRequest("POST", "/", strings.NewReader(body))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		recorder := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(request, recorder)

		err := s.signUpHandler.Handle(c)
		s.Require().NoError(err)

		s.Equal(400, recorder.Code)
		s.JSONEq(fmt.Sprintf(`
		{
			"statusCode": 400,
			"statusText": "BAD_REQUEST",
			"errors": %s
		}
		`, errorMessage), recorder.Body.String())
	}
}

func TestSignUp(t *testing.T) {
	suite.Run(t, new(SignUpHandlerSuite))
}
