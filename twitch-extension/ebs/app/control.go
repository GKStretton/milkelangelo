package app

// thin layer for user-specific auth (e.g. whitelist)

func (app *App) CollectFromVial(vial int) error {
	return app.goo.CollectFromVial(vial)
}

func (app *App) Dispense(x, y float32) error {
	return app.goo.Dispense(x, y)
}

func (app *App) GoToPosition(x, y float32) error {
	return app.goo.GoToPosition(x, y)
}
