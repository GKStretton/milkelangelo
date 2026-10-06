package publicapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gkstretton/dark/services/goo/control"
)

type collectBody struct {
	Id *int `json:"id"`
}

type positionBody struct {
	X *float32 `json:"x"`
	Y *float32 `json:"y"`
}

func (a *Api) getState(w http.ResponseWriter, r *http.Request) {
	data, err := a.marshalState()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(data)
}

func (a *Api) handleClaim(w http.ResponseWriter, r *http.Request) {
	token, err := a.claim(r.Header.Get(claimHeader))
	if err != nil {
		writeError(w, http.StatusForbidden, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (a *Api) handleUnclaim(w http.ResponseWriter, r *http.Request) {
	if err := a.unclaim(r.Header.Get(claimHeader)); err != nil {
		writeError(w, http.StatusForbidden, err)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (a *Api) handleCollect(w http.ResponseWriter, r *http.Request) {
	var body collectBody
	if !a.authorise(w, r, &body) {
		return
	}
	if body.Id == nil {
		writeError(w, http.StatusBadRequest, errors.New("missing id"))
		return
	}

	l.Printf("collect request: vial %d\n", *body.Id)
	respond(w, control.Collect(*body.Id))
}

func (a *Api) handleDispense(w http.ResponseWriter, r *http.Request) {
	var body positionBody
	if !a.authorise(w, r, &body) {
		return
	}
	if body.X == nil || body.Y == nil {
		writeError(w, http.StatusBadRequest, errors.New("missing x or y"))
		return
	}

	l.Printf("dispense request: %.3f, %.3f\n", *body.X, *body.Y)
	respond(w, control.Dispense(*body.X, *body.Y))
}

func (a *Api) handleGoTo(w http.ResponseWriter, r *http.Request) {
	var body positionBody
	if !a.authorise(w, r, &body) {
		return
	}
	if body.X == nil || body.Y == nil {
		writeError(w, http.StatusBadRequest, errors.New("missing x or y"))
		return
	}

	respond(w, control.GoTo(*body.X, *body.Y))
}

// authorise checks the caller holds the claim and decodes the body. It writes
// the error response and returns false on failure.
func (a *Api) authorise(w http.ResponseWriter, r *http.Request, body any) bool {
	if err := a.canControl(r.Header.Get(claimHeader)); err != nil {
		writeError(w, http.StatusForbidden, err)
		return false
	}
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Errorf("invalid body: %w", err))
		return false
	}
	return true
}

// respond reports whether control accepted the command. Accepted commands
// carry on running after the response is sent.
func respond(w http.ResponseWriter, err error) {
	if err != nil {
		writeError(w, http.StatusConflict, err)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}
