package main

import (
	"context"
	"dickobrazz/application"
	"dickobrazz/server"
	"log"
	"time"
)

func main() {
	app := application.NewApplication()
	srv, err := server.New()
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("Server shutdown error: %v", err)
		}
		app.Shutdown()
	}()

	srv.Start()
	app.Run()
}
