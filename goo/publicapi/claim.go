package publicapi

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"time"
)

const controlClaimTimeout = 8 * time.Second

var (
	errClaimedByOther = errors.New("someone else has control")
	errNotClaimant    = errors.New("you don't have control, claim it first")
)

// claim gives control to the caller if nobody has it, returning a token the
// caller sends with each command. Calling it with the current token renews
// the claim. A claim expires if not renewed within controlClaimTimeout.
func (a *Api) claim(token string) (string, error) {
	a.lock.Lock()
	defer a.lock.Unlock()

	if a.claimToken != "" {
		if !tokensEqual(token, a.claimToken) {
			return "", errClaimedByOther
		}
		a.expiryTimer.Reset(controlClaimTimeout)
		return a.claimToken, nil
	}

	newToken, err := randomToken()
	if err != nil {
		return "", err
	}

	l.Println("control claimed")
	a.claimToken = newToken
	a.expiryTimer = time.AfterFunc(controlClaimTimeout, func() {
		a.lock.Lock()
		defer a.lock.Unlock()

		if a.claimToken == newToken {
			l.Println("control claim expired")
			a.claimToken = ""
			a.notify()
		}
	})
	a.notify()

	return newToken, nil
}

func (a *Api) unclaim(token string) error {
	a.lock.Lock()
	defer a.lock.Unlock()

	if a.claimToken == "" || !tokensEqual(token, a.claimToken) {
		return errNotClaimant
	}

	l.Println("control released")
	a.claimToken = ""
	a.expiryTimer.Stop()
	a.notify()

	return nil
}

// canControl returns nil if token is the current claim token.
func (a *Api) canControl(token string) error {
	a.lock.Lock()
	defer a.lock.Unlock()

	if a.claimToken == "" || !tokensEqual(token, a.claimToken) {
		return errNotClaimant
	}
	return nil
}

func tokensEqual(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

func randomToken() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
