package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
)

var sharedSecretCache string

type twitchClaims struct {
	OpaqueUserID string `json:"opaque_user_id"`
	UserID       string `json:"user_id"`
	ChannelID    string `json:"channel_id"`
	// broadcaster or external
	Role       string `json:"role"`
	IsUnlinked bool   `json:"is_unlinked"`
	jwt.StandardClaims
}

func verifyTwitchTokenString(tokenString string) (*twitchClaims, error) {
	sharedSecret, err := getSharedSecret()
	if err != nil {
		return nil, fmt.Errorf("failed to get shared secret: %v", err)
	}

	token, err := jwt.ParseWithClaims(tokenString, &twitchClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(sharedSecret), nil
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

// The verifyUserRequest function can now use verifyTokenString internally.
func verifyUserRequest(r *http.Request) (*twitchClaims, error) {
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
	return verifyTwitchTokenString(tokenString)
}

func verifyInternalTokenString(tokenString string) (*jwt.StandardClaims, error) {
	sharedSecret, err := getSharedSecret()
	if err != nil {
		return nil, fmt.Errorf("failed to get shared secret: %v", err)
	}

	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(sharedSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New("invalid token")
	}
}

func verifyInternalRequest(r *http.Request) (*jwt.StandardClaims, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, errors.New("authorization header is missing")
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, errors.New("invalid token format")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := verifyInternalTokenString(tokenString)
	if err != nil {
		return nil, err
	}

	// ensure it comes from goo, not twitch
	if claims.Issuer != "depth/goo" {
		return nil, fmt.Errorf("invalid token issuer: %s", claims.Issuer)
	}
	if claims.Audience != "StudyOfLightTwitchEBS" {
		return nil, fmt.Errorf("invalid audience: %s", claims.Audience)
	}

	return claims, nil
}

func getSharedSecret() (string, error) {
	if sharedSecretCache != "" {
		return sharedSecretCache, nil
	}
	b, err := os.ReadFile(".shared-secret")
	if err != nil {
		return "", err
	}
	decodedBytes, err := base64.StdEncoding.DecodeString(string(b))
	if err != nil {
		return "", fmt.Errorf("failed to decode shared secret: %v", err)
	}
	sharedSecretCache = string(decodedBytes)
	return sharedSecretCache, nil
}
