package handlers

import (
	"log/slog"
	"net/http"

	"dickobrazz/server/config"
	"dickobrazz/server/models"
	"dickobrazz/server/repository"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	log   *slog.Logger
	cfg   config.Config
	repos repository.Repositories
}

func New(log *slog.Logger, cfg config.Config, repos repository.Repositories) *Handler {
	return &Handler{log: log, cfg: cfg, repos: repos}
}

func respondError(c *gin.Context, status int, message string) {
	c.AbortWithStatusJSON(status, models.APIResponse{
		Error:     &models.APIError{Message: message},
		RequestID: requestIDFromContext(c),
	})
}

func respondOK(c *gin.Context, payload any) {
	c.JSON(http.StatusOK, models.APIResponse{
		Data:      payload,
		RequestID: requestIDFromContext(c),
	})
}

func requestIDFromContext(c *gin.Context) string {
	if value, ok := c.Get("request_id"); ok {
		if id, ok := value.(string); ok {
			return id
		}
	}
	return ""
}
