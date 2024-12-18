package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gkstretton/study-of-light/twitch-ebs/server/openapi"
)

func (s *server) CollectFromVial(c *gin.Context) {
	var collectionRequest *openapi.CollectFromVialJSONRequestBody
	err := c.Bind(&collectionRequest)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("could not bind: %w", err))
		return
	}

	if collectionRequest == nil || collectionRequest.Id == nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("could not get id"))
		return
	}

	l.Infof("received collection request from vial %d", *collectionRequest.Id)

	err = s.app.CollectFromVial(c, *collectionRequest.Id)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusAccepted)
}

func (s *server) Dispense(c *gin.Context) {
	var dispenseRequest *openapi.DispenseJSONRequestBody
	err := c.Bind(&dispenseRequest)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("could not bind: %w", err))
		return
	}

	if dispenseRequest == nil || dispenseRequest.X == nil || dispenseRequest.Y == nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("could not get x and y fields from request"))
		return
	}

	x, y := *dispenseRequest.X, *dispenseRequest.Y

	l.Infof("received dispense request x: %f, y: %f", x, y)

	err = s.app.Dispense(c, x, y)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusAccepted)
}

func (s *server) GoToPosition(c *gin.Context) {
	var goToRequest *openapi.GoToPositionJSONRequestBody
	err := c.Bind(&goToRequest)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("could not bind: %w", err))
		return
	}

	if goToRequest == nil || goToRequest.X == nil || goToRequest.Y == nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("could not get x and y fields from request"))
		return
	}

	x, y := *goToRequest.X, *goToRequest.Y

	l.Infof("received goTo request x: %f, y: %f", x, y)

	err = s.app.GoToPosition(c, x, y)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusAccepted)
}

func (s *server) ClaimControl(c *gin.Context) {
	err := s.app.ClaimControl(c)
	if err != nil {
		c.AbortWithError(http.StatusForbidden, err)
		return
	}
	c.Status(http.StatusAccepted)
}
