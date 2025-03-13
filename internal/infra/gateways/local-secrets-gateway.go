package gateways

import (
	"fmt"
	"os"
	"strings"
)

type LocalSecretsGateway struct {
	PathToFile string
}

func (l *LocalSecretsGateway) Get(key string) (string, error) {
	file, err := os.ReadFile(l.PathToFile)

	if err != nil {
		return "", err
	}

	replacer := strings.NewReplacer("\r", "", "\t", "")
	sanitized := replacer.Replace(string(file))
	variables := strings.SplitSeq(string(sanitized), "\n")

	for variable := range variables {
		if variable == "" {
			continue
		}

		parts := strings.Split(variable, "=")

		if parts[0] == key {
			return parts[1], nil
		}
	}

	return "", fmt.Errorf("secret %s not found", key)
}
