package gateways

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type AwsSecretsGateway struct {
	SecretsClient *secretsmanager.Client
}

func (a *AwsSecretsGateway) Get(key string) (string, error) {
	if _, ok := os.LookupEnv("AWS_SECRET_NAME"); !ok {
		return "", errors.New("environment variable AWS_SECRET_NAME not set")
	}

	secretValue, err := a.SecretsClient.GetSecretValue(context.TODO(), &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(os.Getenv("AWS_SECRET_NAME")),
	})
	if err != nil {
		return "", err
	}

	var secret map[string]any

	err = json.Unmarshal([]byte(*secretValue.SecretString), &secret)
	if err != nil {
		return "", err
	}

	value, exists := secret[key]

	if !exists {
		return "", fmt.Errorf("key %s not found in secret", key)
	}

	switch v := value.(type) {
	case string:
		return v, nil
	case float64:
		return fmt.Sprintf("%g", v), nil
	case bool:
		return fmt.Sprintf("%t", v), nil
	default:
		return "", fmt.Errorf("secret value must be string, float or bool")
	}
}
