package main

import (
	"log"

	"github.com/camtrik/ebbilogue-backend/internal/config"
	"github.com/camtrik/ebbilogue-backend/internal/handler"
	"github.com/camtrik/ebbilogue-backend/internal/service/psn"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Config struct {
	PSNRefreshToken string
}

func (c *Config) GetPSNRefreshToken() string {
	return c.PSNRefreshToken
}

// func testTrophyTitle() {
// 	cfg := config.Load()
// 	psnService := psn.NewPSNService(cfg)
// 	resp, err := psnService.GetUserTitles("me", nil)
// 	if err != nil {
// 		log.Fatalf("Failed to get user titles: %v", err)
// 	}

// 	fmt.Printf("Total games: %d\n", resp.TotalItemCount)
// 	for _, title := range resp.TrophyTitles {
// 		fmt.Printf("\nGame: %s\n", title.TrophyTitleName)
// 		fmt.Printf("Platform: %s\n", title.TrophyTitlePlatform)
// 		fmt.Printf("Progress: %d%%\n", title.Progress)
// 		fmt.Printf("Trophies - Platinum: %d, Gold: %d, Silver: %d, Bronze: %d\n",
// 			title.EarnedTrophies.Platinum,
// 			title.EarnedTrophies.Gold,
// 			title.EarnedTrophies.Silver,
// 			title.EarnedTrophies.Bronze)
// 	}
// }

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Application panicked: %v", r)
		}
	}()
	// testTrophyTitle()
	log.Printf("Application starting...")
	cfg := config.Load()
	log.Printf("Config loaded: PSN_ACCOUNT_ID exists: %v", cfg.PSNAccountID != "")

	psnService := psn.NewPSNService(
		cfg.PSNAccountID,
		cfg.PSNRefreshToken,
	)

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
	r.GET("/api/psn/me/trophyTitles", psnHandler.GetMyTitles)
	r.GET("/api/psn/me/trophyTitles/filtered", psnHandler.GetMyFilteredTitles)

	r.GET("/api/psn/:accountId/trophyTitles", psnHandler.GetUserTitles)
	r.GET("/api/psn/:accountId/trophyTitles/filtered", psnHandler.GetUserFilteredTitles)

	r.Run(":6061")

	log.Printf("Server starting on port 6061...")
	if err := r.Run(":6061"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
