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

	if time.Since(a.ConnectedUserTimestamp) > controlClaimTimeout {
		return fmt.Errorf("control claim expired")
	}

	return nil
}

// ClaimControl initally claims control or renews control for a user.
func (a *App) ClaimControl(ctx context.Context) error {
	user, err := getUserFromContext(ctx)
	if err != nil {
		return err
	}

	a.lock.Lock()
	defer a.lock.Unlock()

	if time.Since(a.ConnectedUserTimestamp) > controlClaimTimeout {
		a.ConnectedUser = nil
	}

	if a.ConnectedUser != nil && a.ConnectedUser.OUID != user.OUID {
		return fmt.Errorf("another user has already claimed control")
	}

	l.Debugf("user %s claimed / renewed control", user.OUID)

	a.ConnectedUser = user
	a.ConnectedUserTimestamp = time.Now()

	go a.reportEbsState()

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
