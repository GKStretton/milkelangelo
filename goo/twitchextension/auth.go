package twitchextension

import (
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/gkstretton/dark/services/goo/keyvalue"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

// returns token used to send to extensions pubsub through Helix api.
func getBroadcastToken(dur time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":        time.Now().Add(dur).UTC().Unix(),
		"user_id":    channelId,
		"role":       "external",
		"channel_id": channelId,
		"pubsub_perms": map[string]interface{}{
			"send": []string{"broadcast"},
		},
	})
	secret, err := getSharedSecret()
	if err != nil {
		return "", err
	}
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

// returns custom token for listening from the EBS
func getEBSListeningToken(dur time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		IssuedAt:  time.Now().UTC().Unix(),
		ExpiresAt: time.Now().Add(dur).UTC().Unix(),
		NotBefore: time.Now().Add(-time.Minute).UTC().Unix(),
		Audience:  "StudyOfLightTwitchEBS",
		Id:        uuid.New().String(),
		Issuer:    "depth/goo",
		Subject:   "depth/goo",
	})
	secret, err := getSharedSecret()
	if err != nil {
		return "", err
	}
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func getSharedSecret() (string, error) {
	secret := keyvalue.Get("TWITCH_EXTENSION_SECRET")
	if secret == nil {
		return "", errors.New("no value for TWITCH_EXTENSION_SECRET in kv")
	}
	decodedBytes, err := base64.StdEncoding.DecodeString(string(secret))
	if err != nil {
		return "", fmt.Errorf("failed to decode shared secret: %v", err)
	}
	return string(decodedBytes), nil
}
