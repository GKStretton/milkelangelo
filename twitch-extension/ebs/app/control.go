package app

import (
	"context"
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

func (a *App) GetState(ctx context.Context) (ebsState, error) {
	return a.buildStateResponse(), nil
}
