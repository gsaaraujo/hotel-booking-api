package gateways_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/gsaaraujo/hotel-booking-api/internal/infra/gateways"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type AwsSecretsGatewaySuite struct {
	suite.Suite
	secretsClient     *secretsmanager.Client
	awsContainer      testcontainers.Container
	awsSecretsGateway gateways.AwsSecretsGateway
}

func (a *AwsSecretsGatewaySuite) SetupTest() {
	os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
	ctx := context.Background()
	awsContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		Started: true,
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "localstack/localstack:latest",
			ExposedPorts: []string{"4566/tcp"},
			WaitingFor:   wait.ForLog("Ready.").WithStartupTimeout(10 * time.Second),
			Env: map[string]string{
				"SERVICES": "secretsmanager",
			},
		},
	})
	a.Require().NoError(err)

	a.awsContainer = awsContainer

	host, err := awsContainer.Host(ctx)
	a.Require().NoError(err)

	port, err := awsContainer.MappedPort(ctx, "4566/tcp")
	a.Require().NoError(err)

	if _, ok := os.LookupEnv("ACT"); ok {
		host = "host.docker.internal"
	}

	awsConfig, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("test", "test", "test")),
	)
	a.Require().NoError(err)

	secretsClient := secretsmanager.NewFromConfig(awsConfig, func(o *secretsmanager.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("http://%s:%s", host, port.Port()))
	})

	a.secretsClient = secretsClient
	a.awsSecretsGateway = gateways.AwsSecretsGateway{
		SecretsClient: secretsClient,
	}
}

func (a *AwsSecretsGatewaySuite) TearDownTest() {
	ctx := context.Background()
	err := a.awsContainer.Terminate(ctx)
	a.Require().NoError(err)
}

func (a *AwsSecretsGatewaySuite) TestGet_OnSuccess_ReturnsSecret() {
	type Secret struct {
		AnyKeyString  string `json:"ANY_KEY_STRING"`
		AnyKeyNumber  string `json:"ANY_KEY_NUMBER"`
		AnyKeyBoolean string `json:"ANY_KEY_BOOLEAN"`
	}
	secret := Secret{
		AnyKeyString:  "any_key_string",
		AnyKeyNumber:  "any_key_number",
		AnyKeyBoolean: "any_key_boolean",
	}
	secretEncoding, err := json.Marshal(secret)
	a.Require().NoError(err)
	os.Setenv("AWS_SECRET_NAME", "any_secret_name")
	secretJson := string(secretEncoding)
	_, err = a.secretsClient.CreateSecret(context.TODO(), &secretsmanager.CreateSecretInput{
		Name:         aws.String("any_secret_name"),
		SecretString: &secretJson,
	})
	a.Require().NoError(err)
	inputs := map[string]string{
		"ANY_KEY_STRING":  "any_key_string",
		"ANY_KEY_NUMBER":  "any_key_number",
		"ANY_KEY_BOOLEAN": "any_key_boolean",
	}

	for key, value := range inputs {
		sut, err := a.awsSecretsGateway.Get(key)

		a.Require().NoError(err)
		a.Equal(value, sut)
	}
}

func (a *AwsSecretsGatewaySuite) TestGet_OnEnvironmentVariableNotSet_ReturnsError() {
	_, err := a.awsSecretsGateway.Get("ANY_KEY_STRING")

	a.EqualError(err, "environment variable AWS_SECRET_NAME not set")
}

func (a *AwsSecretsGatewaySuite) TestGet_OnSecretKeyNotFound_ReturnsError() {
	type Secret struct {
		AnyKeyString string `json:"ANY_KEY_STRING"`
	}
	secret := Secret{
		AnyKeyString: "any_key_string",
	}
	secretEncoding, err := json.Marshal(secret)
	a.Require().NoError(err)
	secretJson := string(secretEncoding)
	os.Setenv("AWS_SECRET_NAME", "any_secret_name")
	_, err = a.secretsClient.CreateSecret(context.TODO(), &secretsmanager.CreateSecretInput{
		Name:         aws.String("any_secret_name"),
		SecretString: &secretJson,
	})
	a.Require().NoError(err)

	_, err = a.awsSecretsGateway.Get("POSTGRES_URL")

	a.EqualError(err, "key POSTGRES_URL not found in secret")
}

func (a *AwsSecretsGatewaySuite) TestGet_OnSecretValuesNotBeingStringFloatBool_ReturnsError() {
	type Secret struct {
		AnyKeyArray []string `json:"ANY_KEY_ARRAY"`
	}
	secret := Secret{
		AnyKeyArray: []string{"any_key_array"},
	}
	secretEncoding, err := json.Marshal(secret)
	a.Require().NoError(err)
	os.Setenv("AWS_SECRET_NAME", "any_secret_name")
	secretJson := string(secretEncoding)
	_, err = a.secretsClient.CreateSecret(context.TODO(), &secretsmanager.CreateSecretInput{
		Name:         aws.String("any_secret_name"),
		SecretString: &secretJson,
	})
	a.Require().NoError(err)

	_, err = a.awsSecretsGateway.Get("ANY_KEY_ARRAY")

	a.EqualError(err, "secret value must be string, float or bool")
}

func TestAwsSecretsGateway(t *testing.T) {
	suite.Run(t, new(AwsSecretsGatewaySuite))
}
