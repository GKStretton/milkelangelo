package gooapi

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

func (g *connectedGooApi) verifyInternalRequest(r *http.Request) (*jwt.StandardClaims, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, errors.New("authorization header is missing")
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, errors.New("invalid token format")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := g.verifyInternalTokenString(tokenString)
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

func (g *connectedGooApi) verifyInternalTokenString(tokenString string) (*jwt.StandardClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return g.internalSecret, nil
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
