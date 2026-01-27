package handlers

import (
	"net/http"

	"dickobrazz/server/middleware"
	"dickobrazz/server/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) Me(c *gin.Context) {
	authContext, ok := middleware.GetAuthContext(c)
	if !ok || authContext.User == nil {
		respondError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	user := models.UserProfile{
		ID:        authContext.User.ID,
		Username:  authContext.User.Username,
		FirstName: authContext.User.FirstName,
		LastName:  authContext.User.LastName,
		PhotoURL:  authContext.User.PhotoURL,
	}

	respondOK(c, user)
}
