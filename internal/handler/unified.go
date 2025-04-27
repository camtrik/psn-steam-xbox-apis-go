package handler

import (
	"net/http"

	"github.com/camtrik/psn-steam-api/internal/global"
	unified "github.com/camtrik/psn-steam-api/internal/service"
	"github.com/gin-gonic/gin"
)

type UnifiedGameHandler struct {
	UnifiedService *unified.UnifiedGameService
	cfg            *global.Config
}

func NewUnifiedGameHandler(unifiedService *unified.UnifiedGameService, cfg *global.Config) *UnifiedGameHandler {
	return &UnifiedGameHandler{
		UnifiedService: unifiedService,
		cfg:            cfg,
	}
}

func (h *UnifiedGameHandler) GetRecentlyPlayedGames(c *gin.Context) {
	psnAccountId := c.Query("psnAccountId")
	steamId := c.Query("steamId")
	timeRangeStr := c.Query("timeRange")

	if psnAccountId == "" {
		psnAccountId = "me"
	}

	if steamId == "" {
		steamId = h.cfg.SteamId
	}
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

	resp, err := h.UnifiedService.GetRecentlyPlayedGames(c.Request.Context(), psnAccountId, steamId, timeRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
