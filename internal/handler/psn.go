package handler

import (
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

func (h *PSNHandler) GetUserTitles(c *gin.Context) {
	resp, err := h.psnService.GetUserTitles("me", nil)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, resp)
}
