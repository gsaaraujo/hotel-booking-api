package gateways

type FakeSecretsGateway struct {
	Secrets map[string]string
}

func (f *FakeSecretsGateway) Get(key string) (string, error) {
	return f.Secrets[key], nil

}
