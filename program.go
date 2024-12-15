package main

import (
	"dickobot/application"
)

func main() {
	app := application.NewApplication()
	defer app.Shutdown()
	app.Run()
}
