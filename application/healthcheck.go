package application

import (
	"context"
	"dickobrazz/application/logging"
	"net/http"
	"sync"
	"time"
)

type HealthcheckServer struct {
	server *http.Server
	log    *logging.Logger
	wg     *sync.WaitGroup
}

func InitializeHealthcheckServer(log *logging.Logger, wg *sync.WaitGroup) *HealthcheckServer {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthCheckHandler)
	server := &http.Server{
		Addr:         ":80",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	return &HealthcheckServer{server: server, log: log, wg: wg}
}

func (hs *HealthcheckServer) Start() {
	hs.wg.Add(1)
	go func() {
		defer hs.wg.Done()
		hs.log.I("Starting HTTP healthcheck server on :80")
		if err := hs.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			hs.log.E("HTTP healthcheck server failed", logging.InnerError, err)
		}
	}()
}

func (hs *HealthcheckServer) Shutdown(ctx context.Context) error {
	hs.log.I("Shutting down HTTP healthcheck server")
	return hs.server.Shutdown(ctx)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok","service":"dickobrazz"}`))
}
