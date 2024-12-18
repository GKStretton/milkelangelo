package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gkstretton/study-of-light/twitch-ebs/app"
	"github.com/gkstretton/study-of-light/twitch-ebs/common"
	"github.com/gkstretton/study-of-light/twitch-ebs/gooapi"
	"github.com/gkstretton/study-of-light/twitch-ebs/server/openapi"
	"github.com/op/go-logging"
)

var l = logging.MustGetLogger("server")

type server struct {
	r   *gin.Engine
	goo gooapi.GooApi
	app *app.App

	// address to listen on
	addr         string
	sharedSecret []byte
}

func (s *server) Run() {
	l.Infof("listening for twitch requests on %s...", s.addr)
	s.r.Run(s.addr)
}

func NewServer(addr string, sharedSecretPath string, goo gooapi.GooApi, app *app.App, disableAuthentication bool) (*server, error) {
	r := gin.Default()

	sharedSecret, err := common.GetSecret(sharedSecretPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get shared secret: %w", err)
	}

	s := &server{
		r:            r,
		goo:          goo,
		app:          app,
		addr:         addr,
		sharedSecret: sharedSecret,
	}

	s.r.Use(corsMiddleware)
	// todo: s.r.Use(rateLimiterMiddleware)
	if disableAuthentication {
		s.r.Use(s.localAuthMiddleware)
	} else {
		s.r.Use(s.twitchAuthMiddleware)
	}

	openapi.RegisterHandlers(r, s)

	return s, nil
}

func corsMiddleware(c *gin.Context) {
	// vscode openapi extension for calling endpoints uses a strange origin,
	// so setting this to *
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, Authorization, Cache-Control, X-Twitch-Extension-Client-Id")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}

	c.Next()
}
