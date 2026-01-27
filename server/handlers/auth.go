package handlers

import (
	"net/http"

	"dickobrazz/server/auth"
	"dickobrazz/server/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) Login(c *gin.Context) {
	var payload models.TelegramAuthPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		respondError(c, http.StatusBadRequest, "invalid_payload")
		return
	}

	if err := auth.ValidateTelegramAuth(payload, h.cfg.TelegramToken); err != nil {
		respondError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	user := models.UserProfile{
		ID:        payload.ID,
		Username:  payload.Username,
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		PhotoURL:  payload.PhotoURL,
	}

	respondOK(c, models.AuthResponse{User: user})
}
