package handler

import (
	"net/http"

	"github.com/camtrik/psn-steam-api/internal/service/xbox"
	"github.com/gin-gonic/gin"
)

type XboxHandler struct {
	xboxService *xbox.XboxService
}

func NewXboxHandler(xboxService *xbox.XboxService) *XboxHandler {
	return &XboxHandler{xboxService: xboxService}
}

func (h *XboxHandler) GetPlayerAchievements(c *gin.Context) {
	resp, err := h.xboxService.GetPlayerAchievements(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *XboxHandler) GetGameStats(c *gin.Context) {
	resp, err := h.xboxService.GetGameStats(c.Request.Context(), c.Param("titleId"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *XboxHandler) GetRecentlyPlayedGames(c *gin.Context) {
	timeRangeStr := c.Query("timeRange")

	var timeRange int64
	switch timeRangeStr {
	case "two_weeks":
		timeRange = 14 * 24 * 60 * 60
	case "one_month":
		timeRange = 30 * 24 * 60 * 60
	case "three_months":
		timeRange = 90 * 24 * 60 * 60
	default:
		timeRange = 30 * 24 * 60 * 60
	}

	resp, err := h.xboxService.GetRecentlyPlayedGames(c.Request.Context(), timeRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
