package ebsinterface

import (
	"errors"
	"time"

	"github.com/gkstretton/dark/services/goo/config"
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
	s := config.SharedSecretEbs()
	if s == "" {
		return "", errors.New("ebs shared secret not set")
	}
	return s, nil
}
