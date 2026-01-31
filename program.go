package main

import (
	"dickobrazz/application"
)

func main() {
	app := application.NewApplication()

	defer func() {
		app.Shutdown()
	}()

	app.Run()
}
