// internal/handler/psn.go
package handler

import (
	"strconv"

	"github.com/camtrik/ebbilogue-backend/internal/service/psn"
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

// Get My Trophies
func (h *PSNHandler) GetMyTitles(c *gin.Context) {
	resp, err := h.psnService.GetMyTitles(nil)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, resp)
}

// Get User Trophies by accountId
func (h *PSNHandler) GetUserTitles(c *gin.Context) {
	accountId := c.Param("accountId")
	resp, err := h.psnService.GetUserTitles(accountId, nil)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, resp)
}

// Get My Filtered Trophies
func (h *PSNHandler) GetMyFilteredTitles(c *gin.Context) {
	filter := psn.TrophyFilter{
		MinProgress: 0,
		Platform:    c.Query("platform"),
		SortBy:      c.Query("sortBy"),
	}

	if progressStr := c.Query("minProgress"); progressStr != "" {
		if progress, err := strconv.Atoi(progressStr); err == nil {
			filter.MinProgress = progress
		}
	}

	resp, err := h.psnService.GetMyFilteredTitles(nil, filter)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, resp)
}

// Get User Filtered Trophies by accountId
func (h *PSNHandler) GetUserFilteredTitles(c *gin.Context) {
	accountId := c.Param("accountId")

	filter := psn.TrophyFilter{
		MinProgress: 0,
		Platform:    c.Query("platform"),
		SortBy:      c.Query("sortBy"),
	}

	if progressStr := c.Query("minProgress"); progressStr != "" {
		if progress, err := strconv.Atoi(progressStr); err == nil {
			filter.MinProgress = progress
		}
	}

	resp, err := h.psnService.GetFilteredUserTitles(accountId, nil, filter)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, resp)
}
