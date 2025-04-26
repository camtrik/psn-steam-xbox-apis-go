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
