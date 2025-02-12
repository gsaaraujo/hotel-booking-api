package main

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/gsaaraujo/hotel-booking-api/internal/application/usecases"
	"github.com/gsaaraujo/hotel-booking-api/internal/infra/gateways"
	"github.com/gsaaraujo/hotel-booking-api/internal/infra/handlers"
	webhttp "github.com/gsaaraujo/hotel-booking-api/internal/infra/web-http"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

func main() {

	if _, ok := os.LookupEnv("AWS_REGION"); !ok {
		panic("environment variable AWS_REGION not set")
	}

	defaultConfig, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		panic(err)
	}

	secretsClient := secretsmanager.NewFromConfig(defaultConfig)

	awsSecretsGateway := gateways.AwsSecretsGateway{
		SecretsClient: secretsClient,
	}

	postgresUrl, err := awsSecretsGateway.Get("POSTGRES_URL")
	if err != nil {
		panic(err)
	}

	conn, err := pgx.Connect(context.Background(), postgresUrl)
	if err != nil {
		panic(err)
	}

	defer conn.Close(context.Background())

	httpLogger := webhttp.NewHttpLogger()

	customersGateway := gateways.CustomersGateway{
		Conn: conn,
	}

	loginWithEmailAndPassword := usecases.LoginWithEmailAndPassword{
		SecretsGateway:   &awsSecretsGateway,
		CustomersGateway: &customersGateway,
	}

	loginWithEmailAndPasswordHandler := handlers.LoginWithEmailAndPasswordHandler{
		HttpLogger:                httpLogger,
		LoginWithEmailAndPassword: &loginWithEmailAndPassword,
	}

	e := echo.New()
	api := e.Group("/api")

	api.POST("/login-with-email-and-password", func(c echo.Context) error {
		return loginWithEmailAndPasswordHandler.Handle(c)
	})

	e.Start(":8080")
}
