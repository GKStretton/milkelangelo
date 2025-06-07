package ebsinterface

import (
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/gkstretton/dark/services/goo/keyvalue"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

// returns custom token for listening from the EBS
func getEBSListeningToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		IssuedAt: time.Now().UTC().Unix(),
		// todo: make token expire after a certain time
		ExpiresAt: time.Now().Add(time.Hour * 24 * 365 * 100).UTC().Unix(),
		NotBefore: time.Now().Add(-time.Minute).UTC().Unix(),
		Audience:  "StudyOfLightTwitchEBS",
		Id:        uuid.New().String(),
		Issuer:    "milkelangelo/goo",
		Subject:   "milkelangelo/goo",
	})
	secret, err := getInternalSecret()
	if err != nil {
		return "", err
	}
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func getInternalSecret() (string, error) {
	secret := keyvalue.Get("EBS_INTERNAL_SECRET")
	if secret == nil {
		return "", errors.New("no value for EBS_INTERNAL_SECRET in kv")
	}
	decodedBytes, err := base64.StdEncoding.DecodeString(string(secret))
	if err != nil {
		return "", fmt.Errorf("failed to decode internal secret: %v", err)
	}
	return string(decodedBytes), nil
}
