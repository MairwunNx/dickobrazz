package main

import (
	"dickobrazz/application"
)

func main() {
	app := application.NewApplication()
	defer app.Shutdown()
	app.Run()
}
