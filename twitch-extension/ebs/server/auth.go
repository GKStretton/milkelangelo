package server

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
)

type twitchClaims struct {
	OpaqueUserID string `json:"opaque_user_id"`
	UserID       string `json:"user_id"`
	ChannelID    string `json:"channel_id"`
	// broadcaster or external
	Role       string `json:"role"`
	IsUnlinked bool   `json:"is_unlinked"`
	jwt.StandardClaims
}

// The verifyUserRequest function can now use verifyTokenString internally.
func (s *server) verifyUserRequest(r *http.Request) (*twitchClaims, error) {
	if r.Header.Get("X-Twitch-Extension-Client-Id") != "ihiyqlxtem517wq76f4hn8pvo9is30" {
		return nil, fmt.Errorf("invalid extension client id")
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, errors.New("authorization header is missing")
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, errors.New("invalid token format")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	return s.verifyTwitchTokenString(tokenString)
}

func (s *server) verifyTwitchTokenString(tokenString string) (*twitchClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &twitchClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.sharedSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*twitchClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New("invalid token")
	}
}

func getSharedSecret(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	decodedBytes, err := base64.StdEncoding.DecodeString(string(b))
	if err != nil {
		return "", fmt.Errorf("failed to decode shared secret: %v", err)
	}
	return string(decodedBytes), nil
}
