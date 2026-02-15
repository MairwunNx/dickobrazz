package main

import (
	"dickobrazz/src/app"
	"dickobrazz/src/features"
	"dickobrazz/src/shared"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		shared.Module,
		features.Module,
		app.Module,
		fx.Invoke(func(a *app.Application) {
			defer a.Shutdown()
			a.Run()
		}),
	).Run()
}
