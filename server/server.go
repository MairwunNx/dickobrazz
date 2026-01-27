package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"dickobrazz/server/config"
	"dickobrazz/server/handlers"
	"dickobrazz/server/logging"
	"dickobrazz/server/random"
	"dickobrazz/server/repository"
	"dickobrazz/server/storage"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	cfg        config.Config
	log        *slog.Logger
	httpServer *http.Server
	mongo      *mongo.Client
	redis      *redis.Client
	rnd        *random.Random
}

func New() (*Server, error) {
	log := logging.New()
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	mongoClient, err := storage.ConnectMongo(ctx, log, cfg.MongoURI)
	if err != nil {
		return nil, err
	}

	redisClient, _, err := storage.ConnectRedis(ctx, log, cfg.RedisAddress, cfg.RedisPassword)
	if err != nil {
		return nil, err
	}

	repos := repository.Repositories{
		Mongo: repository.NewMongoRepository(mongoClient, cfg.MongoDatabase),
		Redis: repository.NewRedisRepository(redisClient),
	}

	rnd := random.New(log, cfg.RandomOrgToken)
	handler := handlers.New(log, cfg, repos)
	router := newRouter(cfg, handler)

	httpServer := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	return &Server{
		cfg:        cfg,
		log:        log,
		httpServer: httpServer,
		mongo:      mongoClient,
		redis:      redisClient,
		rnd:        rnd,
	}, nil
}

func (s *Server) Start() {
	s.log.Info("HTTP server is starting", "addr", s.cfg.Address)
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.log.Error("HTTP server failed", "error", err)
		}
	}()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.log.Info("HTTP server is shutting down")
	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.log.Error("HTTP shutdown error", "error", err)
	}

	if s.mongo != nil {
		if err := s.mongo.Disconnect(ctx); err != nil {
			s.log.Error("Mongo disconnect error", "error", err)
		}
	}
	if s.redis != nil {
		if err := s.redis.Close(); err != nil {
			s.log.Error("Redis close error", "error", err)
		}
	}
	return nil
}
