package handlers_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/gsaaraujo/hotel-booking-api/internal/application/usecases"
	"github.com/gsaaraujo/hotel-booking-api/internal/infra/handlers"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type LoginWithEmailAndPasswordMock struct {
	mock.Mock
}

func (l *LoginWithEmailAndPasswordMock) Execute(input usecases.LoginWithEmailAndPasswordInput) (usecases.LoginWithEmailAndPasswordOutput, error) {
	args := l.Called(input)
	return args.Get(0).(usecases.LoginWithEmailAndPasswordOutput), args.Error(1)
}

type LoginWithEmailAndPasswordHandlerSuite struct {
	suite.Suite
	loginWithEmailAndPasswordMock    LoginWithEmailAndPasswordMock
	loginWithEmailAndPasswordHandler handlers.LoginWithEmailAndPasswordHandler
}

func (l *LoginWithEmailAndPasswordHandlerSuite) SetupTest() {
	l.loginWithEmailAndPasswordMock = LoginWithEmailAndPasswordMock{}
	l.loginWithEmailAndPasswordHandler = handlers.LoginWithEmailAndPasswordHandler{
		LoginWithEmailAndPassword: &l.loginWithEmailAndPasswordMock,
	}
}

func (l *LoginWithEmailAndPasswordHandlerSuite) TestHandle_OnNoErrors_ReturnsOk() {
	customerId, err := uuid.Parse("5579e9a0-2596-4a30-9741-c0d4005a0327")
	l.Require().NoError(err)
	loginWithEmailAndPasswordInput := usecases.LoginWithEmailAndPasswordInput{
		Email:         "joe.doe@gmail.com",
		PlainPassword: "123456",
	}
	loginWithEmailAndPasswordOutput := usecases.LoginWithEmailAndPasswordOutput{
		CustomerId:   customerId,
		CustomerName: "Joe Doe",
		AccessToken:  "any_access_token",
	}
	l.loginWithEmailAndPasswordMock.On("Execute", loginWithEmailAndPasswordInput).Return(loginWithEmailAndPasswordOutput, nil)
	request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`
	{
		"email": "joe.doe@gmail.com",
		"password": "123456"
		}
		`))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(request, recorder)

	err = l.loginWithEmailAndPasswordHandler.Handle(c)

	l.Require().NoError(err)
	l.Equal(200, recorder.Code)
	l.JSONEq(`
		{
			"statusCode": 200,
			"statusText": "OK",
			"data": {
				"customerId": "5579e9a0-2596-4a30-9741-c0d4005a0327",
				"customerName": "Joe Doe",
				"accessToken": "any_access_token"
			}
		}
	`, recorder.Body.String())
}

func (l *LoginWithEmailAndPasswordHandlerSuite) TestHandle_OnEmailOrPasswordIsIncorrectError_ReturnsUnauthorized() {
	l.loginWithEmailAndPasswordMock.On("Execute", usecases.LoginWithEmailAndPasswordInput{
		Email:         "joe.doe@gmail.com",
		PlainPassword: "123456",
	}).Return(usecases.LoginWithEmailAndPasswordOutput{}, errors.New("email or password is incorrect"))
	request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`
		{
			"email": "joe.doe@gmail.com",
			"password": "123456"
		}
	`))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(request, recorder)

	err := l.loginWithEmailAndPasswordHandler.Handle(c)

	l.Require().NoError(err)
	l.Equal(401, recorder.Code)
	l.JSONEq(`
		{
			"statusCode": 401,
			"statusText": "UNAUTHORIZED",
			"error": "email or password is incorrect"
		}
	`, recorder.Body.String())
}

func (l *LoginWithEmailAndPasswordHandlerSuite) TestHandle_OnInvalidBody_ReturnsBadRequest() {
	bodiesAndErrors := []map[string]string{
		{
			"body":   `abc`,
			"errors": `["content-type must be application/json"]`,
		},
		{
			"body":   `{}`,
			"errors": `["email is required", "password is required"]`,
		},
		{
			"body": `{
				"email": null,
				"password": null
			}`,
			"errors": `["email is required", "password is required"]`,
		},
		{
			"body": `{
				"email": "",
				"password": ""
			}`,
			"errors": `["email must not be empty", "password must not be empty"]`,
		},
		{
			"body": `{
				"email": " ",
				"password": " "
			}`,
			"errors": `["email must not be empty", "password must not be empty"]`,
		},
		{
			"body": `{
				"email": 1,
				"password": 2
			}`,
			"errors": `["email must be string", "password must be string"]`,
		},
		{
			"body": `{
				"email": "",
				"password": -2
			}`,
			"errors": `["email must not be empty", "password must be string"]`,
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

		err := l.loginWithEmailAndPasswordHandler.Handle(c)
		l.Require().NoError(err)

		l.Equal(400, recorder.Code)
		l.JSONEq(fmt.Sprintf(`
		{
			"statusCode": 400,
			"statusText": "BAD_REQUEST",
			"errors": %s
		}
		`, errorMessage), recorder.Body.String())
	}
}

func TestLoginWithEmailAndPasswordHandler(t *testing.T) {
	suite.Run(t, new(LoginWithEmailAndPasswordHandlerSuite))
}
