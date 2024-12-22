package app

import (
	"context"
	"fmt"
	"time"

	"github.com/gkstretton/study-of-light/twitch-ebs/entities"
)

const (
	controlClaimTimeout = 15 * time.Second
)

// thin layer for user-specific auth (e.g. whitelist)

func (a *App) CollectFromVial(ctx context.Context, vial int) error {
	if err := a.CanUserControl(ctx); err != nil {
		return err
	}

	return a.goo.CollectFromVial(vial)
}

func (a *App) Dispense(ctx context.Context, x, y float32) error {
	if err := a.CanUserControl(ctx); err != nil {
		return err
	}

	return a.goo.Dispense(x, y)
}

func (a *App) GoToPosition(ctx context.Context, x, y float32) error {
	if err := a.CanUserControl(ctx); err != nil {
		return err
	}

	return a.goo.GoToPosition(x, y)
}

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

	if time.Now().After(a.ConnectedUserExpiryTimestamp) {
		return fmt.Errorf("control claim expired")
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

	a.ConnectedUser = user
	a.ConnectedUserExpiryTimestamp = time.Now().Add(controlClaimTimeout)

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
	return nil
}

func (a *App) GetState(ctx context.Context) (ebsState, error) {
	return a.buildStateResponse(), nil
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
