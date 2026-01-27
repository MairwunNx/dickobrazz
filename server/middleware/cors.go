package middleware

import (
	"net/http"
	"strings"

	"dickobrazz/server/config"
	"github.com/gin-gonic/gin"
)

func CORS(cfg config.Config) gin.HandlerFunc {
	allowed := normalizeOrigins(cfg.CorsOrigins)
	allowAll := len(allowed) == 0

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" {
			if allowAll || allowed[origin] {
				c.Header("Access-Control-Allow-Origin", origin)
				c.Header("Vary", "Origin")
			}
		}

		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Authorization,Content-Type,X-Telegram-Init-Data,X-Internal-Token,X-Request-Id")
		c.Header("Access-Control-Expose-Headers", "X-Request-Id")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func normalizeOrigins(origins []string) map[string]bool {
	result := make(map[string]bool, len(origins))
	for _, value := range origins {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		result[value] = true
	}
	return result
}
