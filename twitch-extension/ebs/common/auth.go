package common

import (
	"encoding/base64"
	"fmt"
	"os"
)

func GetSecret(path string) ([]byte, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	decodedBytes, err := base64.StdEncoding.DecodeString(string(b))
	if err != nil {
		return nil, fmt.Errorf("failed to decode secret: %v", err)
	}
	return decodedBytes, nil
}
