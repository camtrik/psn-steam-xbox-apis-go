// internal/handler/psn.go
package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/camtrik/psn-steam-api/internal/models"
	"github.com/camtrik/psn-steam-api/internal/service/psn"
	"github.com/gin-gonic/gin"
)

type PSNHandler struct {
	psnService *psn.PSNService
}

func NewPSNHandler(psnService *psn.PSNService) *PSNHandler {
	return &PSNHandler{
		psnService: psnService,
	}
}

func (h *PSNHandler) GetTokensFromNPSSO(c *gin.Context) {
	var req psn.NPSSOAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	if len(req.NPSSO) != 64 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid NPSSO format. NPSSO should be a 64-character string",
		})
		return
	}

	tokens, err := h.psnService.GetTokensFromNPSSO(req.NPSSO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get tokens from NPSSO",
			"details": err.Error(),
		})
		return
	}

	h.psnService.UpdateRefreshToken(tokens.RefreshToken)

	fmt.Printf("New refresh token: %s\n", tokens.RefreshToken)

	c.JSON(http.StatusOK, gin.H{
		"message":      "Tokens retrieved successfully",
		"data":         tokens,
		"instructions": "Save the refresh_token securely. The access_token will expire soon, but the refresh_token can be used to get new access tokens.",
	})
}

func (h *PSNHandler) parseOptions(c *gin.Context) *models.GetUserTitlesOptions {
	var options *models.GetUserTitlesOptions

	if limitStr := c.Query("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err == nil && limit > 0 {
			if options == nil {
				options = &models.GetUserTitlesOptions{}
			}
			options.Limit = &limit
		}
	}

	if offSetStr := c.Query("offset"); offSetStr != "" {
		offset, err := strconv.Atoi(offSetStr)
		if err == nil && offset > 0 {
			if options == nil {
				options = &models.GetUserTitlesOptions{}
			}
			options.Offset = &offset
		}
	}
	return options
}

func (h *PSNHandler) parseFilter(c *gin.Context) models.TrophyFilter {
	filter := models.TrophyFilter{
		MinProgress: 0,
		Platform:    c.Query("platform"),
		SortBy:      c.Query("sortBy"),
	}

	if progressStr := c.Query("minProgress"); progressStr != "" {
		progress, err := strconv.Atoi(progressStr)
		if err == nil && progress > 0 {
			filter.MinProgress = progress
		}
	}
	return filter
}

// Get User Trophies by accountId
func (h *PSNHandler) GetUserTitles(c *gin.Context) {
	filter := h.parseFilter(c)
	accountId := c.Param("accountId")
	options := h.parseOptions(c)
	resp, err := h.psnService.GetUserTitles(c.Request.Context(), accountId, options, filter)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, resp)
}

func (h *PSNHandler) GetRecentlyPlayedGames(c *gin.Context) {
	accountId := c.Param("accountId")
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

	resp, err := h.psnService.GetRecentlyPlayedGames(c.Request.Context(), accountId, timeRange)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, resp)

}
