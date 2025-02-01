package main

import (
	"fmt"
	"log"

	"github.com/camtrik/ebbilogue-backend/internal/config"
	"github.com/camtrik/ebbilogue-backend/internal/handler"
	"github.com/camtrik/ebbilogue-backend/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Config struct {
	PSNRefreshToken string
}

func (c *Config) GetPSNRefreshToken() string {
	return c.PSNRefreshToken
}

func testTrophyTitle() {
	cfg := config.Load()
	psnService := service.NewPSNService(cfg)
	resp, err := psnService.GetUserTitles("me", nil)
	if err != nil {
		log.Fatalf("Failed to get user titles: %v", err)
	}

	fmt.Printf("Total games: %d\n", resp.TotalItemCount)
	for _, title := range resp.TrophyTitles {
		fmt.Printf("\nGame: %s\n", title.TrophyTitleName)
		fmt.Printf("Platform: %s\n", title.TrophyTitlePlatform)
		fmt.Printf("Progress: %d%%\n", title.Progress)
		fmt.Printf("Trophies - Platinum: %d, Gold: %d, Silver: %d, Bronze: %d\n",
			title.EarnedTrophies.Platinum,
			title.EarnedTrophies.Gold,
			title.EarnedTrophies.Silver,
			title.EarnedTrophies.Bronze)
	}
}

func main() {
	// testTrophyTitle()
	cfg := config.Load()
	psnService := service.NewPSNService(cfg)
	psnHandler := handler.NewPSNHandler(psnService)

	r := gin.Default()

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Routes
	r.GET("/api/psn/trophyTitles", psnHandler.GetUserTitles)

	r.Run(":6061")
}
