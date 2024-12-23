package server

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gkstretton/study-of-light/twitch-ebs/entities"
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

const extensionClientID = "ihiyqlxtem517wq76f4hn8pvo9is30"

func (s *server) localAuthMiddleware(c *gin.Context) {
	// add user object to context
	c.Set(entities.ContextKeyUser, &entities.User{
		OUID: "local",
	})

	l.Debug("added local user to context")

	c.Next()
}

func (s *server) twitchAuthMiddleware(c *gin.Context) {
	if c.GetHeader("X-Twitch-Extension-Client-Id") != extensionClientID {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid extension client id"))
		return
	}

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.AbortWithError(http.StatusUnauthorized, errors.New("authorization header is missing"))
		return
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		c.AbortWithError(http.StatusUnauthorized, errors.New("invalid token format, no bearer prefix"))
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := s.verifyTwitchTokenString(tokenString)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("token verification failed: %v", err))
		return
	}

	l.Debugf("verified user request raddr %s, uid %s, ouid %s", c.Request.RemoteAddr, claims.UserID, claims.OpaqueUserID)

	// add user object to context
	c.Set(entities.ContextKeyUser, &entities.User{
		OUID: claims.OpaqueUserID,
	})

	c.Next()
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
