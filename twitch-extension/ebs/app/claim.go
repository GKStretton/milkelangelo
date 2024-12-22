package app

import (
	"context"
	"fmt"
	"time"

	"github.com/gkstretton/study-of-light/twitch-ebs/entities"
)

const (
	controlClaimTimeout = 8 * time.Second
)

func (a *App) CanUserControl(ctx context.Context) error {
	user, err := getUserFromContext(ctx)
	if err != nil {
		return err
	}

	a.lock.Lock()
	defer a.lock.Unlock()

	if a.ConnectedUser == nil {
		return fmt.Errorf("no user connected")
	}

	if a.ConnectedUser.OUID != user.OUID {
		return fmt.Errorf("user not connected")
	}

	return nil
}

// Claim initally claims control or renews control for a user.
func (a *App) Claim(ctx context.Context) error {
	user, err := getUserFromContext(ctx)
	if err != nil {
		return err
	}

	a.lock.Lock()
	defer a.lock.Unlock()

	if a.ConnectedUser != nil && a.ConnectedUser.OUID != user.OUID {
		return fmt.Errorf("another user has already claimed control")
	}

	l.Debugf("user %s made / renewed claim", user.OUID)

	if a.ConnectedUser == nil {
		a.ConnectedUser = user
		a.expiryTimer = time.AfterFunc(controlClaimTimeout, func() {
			a.lock.Lock()
			defer a.lock.Unlock()

			l.Debugf("control claim expired for user %s", user.OUID)
			if a.ConnectedUser != nil && a.ConnectedUser.OUID == user.OUID {
				a.ConnectedUser = nil
			}
		})
		go a.broadcast()
	} else {
		a.expiryTimer.Reset(controlClaimTimeout)
	}

	go a.reportEbsState()

	return nil
}

func (a *App) Unclaim(ctx context.Context) error {
	user, err := getUserFromContext(ctx)
	if err != nil {
		return err
	}

	a.lock.Lock()
	defer a.lock.Unlock()

	if a.ConnectedUser == nil {
		return fmt.Errorf("no user connected")
	}

	if a.ConnectedUser.OUID != user.OUID {
		return fmt.Errorf("user not connected")
	}

	l.Debugf("user %s released claim", user.OUID)

	a.ConnectedUser = nil
	a.expiryTimer.Stop()

	go a.reportEbsState()
	go a.broadcast()

	return nil
}

func getUserFromContext(ctx context.Context) (*entities.User, error) {
	u := ctx.Value(entities.ContextKeyUser)
	if u == nil {
		return nil, fmt.Errorf("user key not found in context")
	}

	user, ok := u.(*entities.User)
	if !ok {
		return nil, fmt.Errorf("user value not of type *entities.User")
	}

	return user, nil
}
