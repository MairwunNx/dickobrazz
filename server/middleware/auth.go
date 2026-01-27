package middleware

import (
	"net/http"
	"strings"

	"dickobrazz/server/auth"
	"dickobrazz/server/config"
	"dickobrazz/server/models"

	"github.com/gin-gonic/gin"
)

const authContextKey = "auth_context"

type AuthKind string

const (
	AuthKindTelegram AuthKind = "telegram"
	AuthKindInternal AuthKind = "internal"
)

type AuthContext struct {
	Kind AuthKind
	User *models.TelegramAuthPayload
}

func RequireAuth(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if internalContext := authorizeInternal(c, cfg.InternalToken); internalContext != nil {
			c.Set(authContextKey, *internalContext)
			c.Next()
			return
		}

		initData := c.GetHeader("X-Telegram-Init-Data")
		if initData == "" {
			authHeader := c.GetHeader("Authorization")
			if strings.HasPrefix(authHeader, "Telegram ") {
				initData = strings.TrimPrefix(authHeader, "Telegram ")
			}
		}

		if initData == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.APIResponse{
				Error: &models.APIError{Message: "unauthorized"},
			})
			return
		}

		payload, err := auth.ParseInitData(initData)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.APIResponse{
				Error: &models.APIError{Message: "unauthorized"},
			})
			return
		}

		if err := auth.ValidateTelegramAuth(payload, cfg.TelegramToken); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.APIResponse{
				Error: &models.APIError{Message: "unauthorized"},
			})
			return
		}

		c.Set(authContextKey, AuthContext{
			Kind: AuthKindTelegram,
			User: &payload,
		})
		c.Next()
	}
}

func GetAuthContext(c *gin.Context) (AuthContext, bool) {
	value, exists := c.Get(authContextKey)
	if !exists {
		return AuthContext{}, false
	}
	context, ok := value.(AuthContext)
	return context, ok
}

func authorizeInternal(c *gin.Context, token string) *AuthContext {
	if token == "" {
		return nil
	}
	if c.GetHeader("X-Internal-Token") != token {
		return nil
	}
	return &AuthContext{Kind: AuthKindInternal}
}
