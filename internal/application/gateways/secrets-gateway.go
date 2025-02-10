package gateways

type ISecretsGateway interface {
	Get(key string) (string, error)
}
