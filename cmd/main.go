package main

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	applicationgateway "github.com/gsaaraujo/hotel-booking-api/internal/application/gateways"
	"github.com/gsaaraujo/hotel-booking-api/internal/application/usecases"
	"github.com/gsaaraujo/hotel-booking-api/internal/infra/gateways"
	"github.com/gsaaraujo/hotel-booking-api/internal/infra/handlers"
	"github.com/gsaaraujo/hotel-booking-api/internal/infra/repositories"
	webhttp "github.com/gsaaraujo/hotel-booking-api/internal/infra/web-http"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

func main() {
	defaultConfig, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		panic(err)
	}

	secretsClient := secretsmanager.NewFromConfig(defaultConfig)

	var secretsGateway applicationgateway.ISecretsGateway

	if os.Getenv("API_ENV") == "DEV" {
		secretsGateway = &gateways.LocalSecretsGateway{
			PathToFile: ".env",
		}
	} else {
		secretsGateway = &gateways.AwsSecretsGateway{
			SecretsClient: secretsClient,
		}
	}

	postgresUrl, err := secretsGateway.Get("POSTGRES_URL")
	if err != nil {
		panic(err)
	}

	conn, err := pgx.Connect(context.Background(), postgresUrl)
	if err != nil {
		panic(err)
	}

	defer conn.Close(context.Background())

	httpLogger := webhttp.NewHttpLogger()

	httpValidator, err := webhttp.NewHttpValidator()
	if err != nil {
		panic(err)
	}

	httpAuthorization := webhttp.HttpAuthorization{
		SecretsGateway: secretsGateway,
	}

	customersGateway := gateways.CustomersGateway{
		Conn: conn,
	}

	roomRepository := repositories.RoomsRepository{
		Conn: conn,
	}

	loginWithEmailAndPassword := usecases.LoginWithEmailAndPassword{
		SecretsGateway:   secretsGateway,
		CustomersGateway: &customersGateway,
	}

	signUp := usecases.SignUp{
		CustomersGateway: &customersGateway,
	}

	createRoom := usecases.CreateRoom{
		RoomsRepository: &roomRepository,
	}

	loginWithEmailAndPasswordHandler := handlers.LoginWithEmailAndPasswordHandler{
		HttpLogger:                httpLogger,
		LoginWithEmailAndPassword: &loginWithEmailAndPassword,
	}

	signUpHandler := handlers.SignUpHandler{
		SignUp: &signUp,
	}

	createRoomHandler := handlers.CreateRoomHandler{
		HttpLogger:        httpLogger,
		HttpAuthorization: httpAuthorization,
		HttpValidator:     httpValidator,
		CreateRoom:        &createRoom,
	}

	e := echo.New()
	api := e.Group("/api")

	api.POST("/login-with-email-and-password", func(c echo.Context) error {
		return loginWithEmailAndPasswordHandler.Handle(c)
	})

	api.POST("/sign-up", func(c echo.Context) error {
		return signUpHandler.Handle(c)
	})

	api.POST("/create-room", func(c echo.Context) error {
		return createRoomHandler.Handle(c)
	})

	err = e.Start(":8080")
	if err != nil {
		panic(err)
	}
}
