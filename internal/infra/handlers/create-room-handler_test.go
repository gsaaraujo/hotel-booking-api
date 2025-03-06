package handlers_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gsaaraujo/hotel-booking-api/internal/application/gateways"
	"github.com/gsaaraujo/hotel-booking-api/internal/application/usecases"
	"github.com/gsaaraujo/hotel-booking-api/internal/infra/handlers"
	webhttp "github.com/gsaaraujo/hotel-booking-api/internal/infra/web-http"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockCreateRoom struct {
	mock.Mock
}

func (m *MockCreateRoom) Execute(input usecases.CreateRoomInput) error {
	args := m.Called(input)
	return args.Error(0)
}

type CreateRoomHandlerSuite struct {
	suite.Suite
	mockCreateRoom     MockCreateRoom
	fakeSecretsGateway gateways.FakeSecretsGateway
	createRoomHandler  handlers.CreateRoomHandler
}

func (cr *CreateRoomHandlerSuite) SetupTest() {
	httpValidator, err := webhttp.NewHttpValidator()
	cr.Require().NoError(err)

	cr.mockCreateRoom = MockCreateRoom{}
	cr.fakeSecretsGateway = gateways.FakeSecretsGateway{}
	httpAuthorization := webhttp.HttpAuthorization{
		SecretsGateway: &cr.fakeSecretsGateway,
	}
	cr.createRoomHandler = handlers.CreateRoomHandler{
		HttpLogger:        webhttp.NewHttpLogger(),
		HttpValidator:     httpValidator,
		HttpAuthorization: httpAuthorization,
		CreateRoom:        &cr.mockCreateRoom,
	}
}

func (cr *CreateRoomHandlerSuite) TestHandle_OnNoErrors_ReturnsCreated() {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role": "ADMIN",
	})
	signedToken, err := token.SignedString([]byte("6b45b2cb79974f989447f1d850d139f1"))
	cr.Require().NoError(err)
	cr.fakeSecretsGateway.Secrets = map[string]string{"JWT_SIGNING_ACCESS_TOKEN": "6b45b2cb79974f989447f1d850d139f1"}
	cr.mockCreateRoom.On("Execute", usecases.CreateRoomInput{
		Number:   "101",
		Type:     "SUITE",
		Capacity: 2,
		Price:    250,
	}).Return(nil)
	request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`
		{
			"number": "101",
			"type": "SUITE",
			"capacity": 2,
			"price": 250
		}
	`))
	request.Header.Set("Authorization", signedToken)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(request, recorder)

	err = cr.createRoomHandler.Handle(c)
	cr.Require().NoError(err)

	cr.Equal(201, recorder.Code)
	cr.JSONEq(`
		{
			"statusCode": 201,
			"statusText": "CREATED",
			"data": null
		}
	`, recorder.Body.String())
}

func (cr *CreateRoomHandlerSuite) TestHandle_OnAuthorizationTokenIsMissing_ReturnsError() {
	cr.mockCreateRoom.On("Execute", usecases.CreateRoomInput{
		Number:   "101",
		Type:     "SUITE",
		Capacity: 2,
		Price:    250,
	}).Return(nil)
	request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`
		{
			"number": "101",
			"type": "SUITE",
			"capacity": 2,
			"price": 250
		}
	`))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(request, recorder)

	err := cr.createRoomHandler.Handle(c)
	cr.Require().NoError(err)

	cr.Equal(401, recorder.Code)
	cr.JSONEq(`
	{
		"statusCode": 401,
		"statusText": "UNAUTHORIZED",
		"error": "missing or invalid authorization token"
	}
`, recorder.Body.String())

}

func (cr *CreateRoomHandlerSuite) TestHandle_OnNoPermissonToAccessResource_ReturnsError() {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role": "CUSTOMER",
	})
	signedToken, err := token.SignedString([]byte("6b45b2cb79974f989447f1d850d139f1"))
	cr.Require().NoError(err)
	cr.fakeSecretsGateway.Secrets = map[string]string{"JWT_SIGNING_ACCESS_TOKEN": "6b45b2cb79974f989447f1d850d139f1"}
	cr.mockCreateRoom.On("Execute", usecases.CreateRoomInput{
		Number:   "101",
		Type:     "SUITE",
		Capacity: 2,
		Price:    250,
	}).Return(nil)
	request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`
		{
			"number": "101",
			"type": "SUITE",
			"capacity": 2,
			"price": 250
		}
	`))
	request.Header.Set("Authorization", signedToken)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(request, recorder)

	err = cr.createRoomHandler.Handle(c)
	cr.Require().NoError(err)

	cr.Equal(403, recorder.Code)
	cr.JSONEq(`
	{
		"statusCode": 403,
		"statusText": "FORBIDDEN",
		"error": "you do not have permission to access this resource"
	}
`, recorder.Body.String())

}

func (cr *CreateRoomHandlerSuite) TestHandle_OnDuplicateRoomNumberError_ReturnsConflict() {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role": "ADMIN",
	})
	signedToken, err := token.SignedString([]byte("6b45b2cb79974f989447f1d850d139f1"))
	cr.Require().NoError(err)
	cr.fakeSecretsGateway.Secrets = map[string]string{"JWT_SIGNING_ACCESS_TOKEN": "6b45b2cb79974f989447f1d850d139f1"}

	cr.mockCreateRoom.On("Execute", usecases.CreateRoomInput{
		Number:   "101",
		Type:     "SUITE",
		Capacity: 2,
		Price:    250,
	}).Return(errors.New("the room number '101' is already in use. Please assign another room number"))
	request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`
		{
			"number": "101",
			"type": "SUITE",
			"capacity": 2,
			"price": 250
		}
	`))
	request.Header.Set("Authorization", signedToken)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(request, recorder)

	err = cr.createRoomHandler.Handle(c)
	cr.Require().NoError(err)

	cr.Equal(409, recorder.Code)
	cr.JSONEq(`
		{
			"statusCode": 409,
			"statusText": "CONFLICT",
			"error": "the room number '101' is already in use. Please assign another room number"
		}
	`, recorder.Body.String())
}

func (cr *CreateRoomHandlerSuite) TestHandle_OnInvalidNumberError_ReturnsConflict() {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role": "ADMIN",
	})
	signedToken, err := token.SignedString([]byte("6b45b2cb79974f989447f1d850d139f1"))
	cr.Require().NoError(err)
	cr.fakeSecretsGateway.Secrets = map[string]string{"JWT_SIGNING_ACCESS_TOKEN": "6b45b2cb79974f989447f1d850d139f1"}
	cr.mockCreateRoom.On("Execute", usecases.CreateRoomInput{
		Number:   "1",
		Type:     "SUITE",
		Capacity: 2,
		Price:    250,
	}).Return(errors.New("invalid room number format. Please enter a three-digit room number (e.g. 101) where the first digit indicates the floor number, and the last two digits represent the room number on that floor"))
	request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`
		{
			"number": "1",
			"type": "SUITE",
			"capacity": 2,
			"price": 250
		}
	`))
	request.Header.Set("Authorization", signedToken)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(request, recorder)

	err = cr.createRoomHandler.Handle(c)
	cr.Require().NoError(err)

	cr.Equal(409, recorder.Code)
	cr.JSONEq(`
		{
			"statusCode": 409,
			"statusText": "CONFLICT",
			"error": "invalid room number format. Please enter a three-digit room number (e.g. 101) where the first digit indicates the floor number, and the last two digits represent the room number on that floor"
		}
	`, recorder.Body.String())
}

func (cr *CreateRoomHandlerSuite) TestHandle_OnInvalidTypeError_ReturnsConflict() {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role": "ADMIN",
	})
	signedToken, err := token.SignedString([]byte("6b45b2cb79974f989447f1d850d139f1"))
	cr.Require().NoError(err)
	cr.fakeSecretsGateway.Secrets = map[string]string{"JWT_SIGNING_ACCESS_TOKEN": "6b45b2cb79974f989447f1d850d139f1"}
	cr.mockCreateRoom.On("Execute", usecases.CreateRoomInput{
		Number:   "101",
		Type:     "abc",
		Capacity: 2,
		Price:    250,
	}).Return(errors.New("room type must be SINGLE, DOUBLE, TWIN or SUITE"))
	request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`
		{
			"number": "101",
			"type": "abc",
			"capacity": 2,
			"price": 250
		}
	`))
	request.Header.Set("Authorization", signedToken)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(request, recorder)

	err = cr.createRoomHandler.Handle(c)
	cr.Require().NoError(err)

	cr.Equal(409, recorder.Code)
	cr.JSONEq(`
		{
			"statusCode": 409,
			"statusText": "CONFLICT",
			"error": "room type must be SINGLE, DOUBLE, TWIN or SUITE"
		}
	`, recorder.Body.String())
}

func (cr *CreateRoomHandlerSuite) TestHandle_OnInvalidCapacityError_ReturnsConflict() {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role": "ADMIN",
	})
	signedToken, err := token.SignedString([]byte("6b45b2cb79974f989447f1d850d139f1"))
	cr.Require().NoError(err)
	cr.fakeSecretsGateway.Secrets = map[string]string{"JWT_SIGNING_ACCESS_TOKEN": "6b45b2cb79974f989447f1d850d139f1"}
	cr.mockCreateRoom.On("Execute", usecases.CreateRoomInput{
		Number:   "101",
		Type:     "SUITE",
		Capacity: 0,
		Price:    250,
	}).Return(errors.New("invalid room capacity. Please enter a capacity of at least one to accommodate guests"))
	request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`
		{
			"number": "101",
			"type": "SUITE",
			"capacity": 0,
			"price": 250
		}
	`))
	request.Header.Set("Authorization", signedToken)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(request, recorder)

	err = cr.createRoomHandler.Handle(c)
	cr.Require().NoError(err)

	cr.Equal(409, recorder.Code)
	cr.JSONEq(`
		{
			"statusCode": 409,
			"statusText": "CONFLICT",
			"error": "invalid room capacity. Please enter a capacity of at least one to accommodate guests"
		}
	`, recorder.Body.String())
}

func (cr *CreateRoomHandlerSuite) TestHandle_OnInvalidPriceError_ReturnsConflict() {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role": "ADMIN",
	})
	signedToken, err := token.SignedString([]byte("6b45b2cb79974f989447f1d850d139f1"))
	cr.Require().NoError(err)
	cr.fakeSecretsGateway.Secrets = map[string]string{"JWT_SIGNING_ACCESS_TOKEN": "6b45b2cb79974f989447f1d850d139f1"}
	cr.mockCreateRoom.On("Execute", usecases.CreateRoomInput{
		Number:   "101",
		Type:     "SUITE",
		Capacity: 2,
		Price:    0,
	}).Return(errors.New("invalid room price. Please enter a value greater than zero to ensure proper pricing"))
	request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`
		{
			"number": "101",
			"type": "SUITE",
			"capacity": 2,
			"price": 0
		}
	`))
	request.Header.Set("Authorization", signedToken)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(request, recorder)

	err = cr.createRoomHandler.Handle(c)
	cr.Require().NoError(err)

	cr.Equal(409, recorder.Code)
	cr.JSONEq(`
		{
			"statusCode": 409,
			"statusText": "CONFLICT",
			"error": "invalid room price. Please enter a value greater than zero to ensure proper pricing"
		}
	`, recorder.Body.String())
}

func (cr *CreateRoomHandlerSuite) TestHandle_OnAnyUnexpectedError_ReturnsInternalServerError() {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role": "ADMIN",
	})
	signedToken, err := token.SignedString([]byte("6b45b2cb79974f989447f1d850d139f1"))
	cr.Require().NoError(err)
	cr.fakeSecretsGateway.Secrets = map[string]string{"JWT_SIGNING_ACCESS_TOKEN": "6b45b2cb79974f989447f1d850d139f1"}
	cr.mockCreateRoom.On("Execute", usecases.CreateRoomInput{
		Number:   "101",
		Type:     "SUITE",
		Capacity: 2,
		Price:    250,
	}).Return(errors.New("any unexpected error"))
	request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`
		{
			"number": "101",
			"type": "SUITE",
			"capacity": 2,
			"price": 250
		}
	`))
	request.Header.Set("Authorization", signedToken)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(request, recorder)

	err = cr.createRoomHandler.Handle(c)
	cr.Require().NoError(err)

	cr.Equal(500, recorder.Code)
	cr.JSONEq(`
		{
			"statusCode": 500,
			"statusText": "INTERNAL_SERVER_ERROR",
			"error": "something went wrong. Please try again later"
		}
	`, recorder.Body.String())
}

func (cr *CreateRoomHandlerSuite) TestHandle_OnInvalidBody_ReturnsBadRequest() {
	bodiesAndErrors := []map[string]string{
		{
			"body":   `abc`,
			"errors": `["content-type must be application/json"]`,
		},
		{
			"body":   `{}`,
			"errors": `["number is required", "type is required", "capacity is required", "price is required"]`,
		},
		{
			"body": `{
				"number": null,
				"type": null,
				"capacity": null,
				"price": null
			}`,
			"errors": `["number is required", "type is required", "capacity is required", "price is required"]`,
		},
		{
			"body": `{
				"number": "",
				"type": "",
				"capacity": null,
				"price": null
			}`,
			"errors": `["number must not be empty", "type must not be empty", "capacity is required", "price is required"]`,
		},
		{
			"body": `{
				"number": " ",
				"type": " ",
				"capacity": null,
				"price": null
			}`,
			"errors": `["number must not be empty", "type must not be empty", "capacity is required", "price is required"]`,
		},
		{
			"body": `{
				"number": 1,
				"type": 1,
				"capacity": "",
				"price": ""
			}`,
			"errors": `["number must be string", "type must be string", "capacity must be integer", "price must be integer"]`,
		},
		{
			"body": `{
				"number": 1,
				"type": 1,
				"capacity": 0.2,
				"price": 0.5
			}`,
			"errors": `["number must be string", "type must be string", "capacity must be integer", "price must be integer"]`,
		},
		{
			"body": `{
				"number": 1,
				"type": 1,
				"capacity": -1,
				"price": -2
			}`,
			"errors": `["number must be string", "type must be string", "capacity must be positive", "price must be positive"]`,
		},
		{
			"body": `{
				"number": "joieoplkkuyrmmnchskkloaokaweokfopsdjgsdijjoieoplkkuyrmmnchskkloaokaweokfopsdjgsdijjoieoplkkuyrmmnchskkloaokaweokfopsdjgsdijjoieoplkkuyrmmnchskkloaokaweokfopsdjgsdijjoieoplkkuyrmmnchskkloaokaweokfopsdjgsdijjoieoplkkuyrmmnchskkloaokaweokfopsdjgsdijjoieoplkkuyuio",
				"type": "joieoplkkuyrmmnchskkloaokaweokfopsdjgsdijjoieoplkkuyrmmnchskkloaokaweokfopsdjgsdijjoieoplkkuyrmmnchskkloaokaweokfopsdjgsdijjoieoplkkuyrmmnchskkloaokaweokfopsdjgsdijjoieoplkkuyrmmnchskkloaokaweokfopsdjgsdijjoieoplkkuyrmmnchskkloaokaweokfopsdjgsdijjoieoplkkuyuio",
				"capacity": 10000,
				"price": 10000000000
			}`,
			"errors": `["number must be less than 256", "type must be less than 256", "capacity must be less than 1000", "price must be less than 1000000000"]`,
		},
	}

	for _, inputAndError := range bodiesAndErrors {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"role": "ADMIN",
		})
		signedToken, err := token.SignedString([]byte("6b45b2cb79974f989447f1d850d139f1"))
		cr.Require().NoError(err)
		cr.fakeSecretsGateway.Secrets = map[string]string{"JWT_SIGNING_ACCESS_TOKEN": "6b45b2cb79974f989447f1d850d139f1"}
		body := inputAndError["body"]
		errorMessage := inputAndError["errors"]
		request := httptest.NewRequest("POST", "/", strings.NewReader(body))
		request.Header.Set("Authorization", signedToken)
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		recorder := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(request, recorder)

		err = cr.createRoomHandler.Handle(c)
		cr.Require().NoError(err)

		cr.Equal(400, recorder.Code)
		cr.JSONEq(fmt.Sprintf(`
		{
			"statusCode": 400,
			"statusText": "BAD_REQUEST",
			"errors": %s
		}
		`, errorMessage), recorder.Body.String())
	}
}

func TestCreateRoomHandler(t *testing.T) {
	suite.Run(t, new(CreateRoomHandlerSuite))
}
