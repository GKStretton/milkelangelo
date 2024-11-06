package common

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
)

func GetSecret(path string) ([]byte, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	trimmed := strings.TrimSpace(string(b))
	decodedBytes, err := base64.StdEncoding.DecodeString(trimmed)
	if err != nil {
		return nil, fmt.Errorf("failed to decode secret: %v", err)
	}
	return decodedBytes, nil
}
