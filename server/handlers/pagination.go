package handlers

import (
	"strconv"

	"dickobrazz/server/models"
	"github.com/gin-gonic/gin"
)

func paginationFromQuery(c *gin.Context) models.PageMeta {
	meta := models.PageMeta{}

	if value := c.Query("limit"); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil && parsed > 0 {
			meta.Limit = parsed
		}
	}

	if value := c.Query("page"); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil && parsed > 0 {
			meta.Page = parsed
		}
	}

	if value := c.Query("cursor"); value != "" {
		meta.NextCursor = value
	}

	if value := c.Query("prev_cursor"); value != "" {
		meta.PrevCursor = value
	}

	return meta
}
