package handler

import (
	"net/http"
	"strconv"

	"github.com/camtrik/psn-steam-api/internal/service/steam"
	"github.com/gin-gonic/gin"
)

type SteamHandler struct {
	steamService *steam.SteamService
}

func NewSteamHandler(steamService *steam.SteamService) *SteamHandler {
	return &SteamHandler{steamService: steamService}
}

func (h *SteamHandler) GetOwnedGames(c *gin.Context) {
	steamId := c.Param("steamId")

	if steamId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "steamId is required"})
		return
	}

	resp, err := h.steamService.GetOwnedGames(c.Request.Context(), steamId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, resp)
}

func (h *SteamHandler) GetPlayerAchievements(c *gin.Context) {
	steamId := c.Param("steamId")
	appIdStr := c.Param("appId")

	if steamId == "" || appIdStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "steamId and appId are required"})
		return
	}

	appId, err := strconv.Atoi(appIdStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid appId"})
		return
	}

	resp, err := h.steamService.GetPlayerAchievements(c.Request.Context(), steamId, appId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *SteamHandler) GetPlayerGameDetails(c *gin.Context) {
	steamId := c.Param("steamId")
	minPlayTimeStr := c.Query("minPlayTime")
	sortByTimeStr := c.Query("sortByTime")

	minPlayTime, _ := strconv.Atoi(minPlayTimeStr)
	sortByTime := sortByTimeStr == "true"

	if steamId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "steamId is required"})
		return
	}

	resp, err := h.steamService.GetPlayerGameDetails(c.Request.Context(), steamId, minPlayTime, sortByTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}
