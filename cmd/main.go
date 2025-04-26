package main

import (
	"log"
	"os"

	helper "github.com/camtrik/psn-steam-api/internal"
	"github.com/camtrik/psn-steam-api/internal/cache"
	"github.com/camtrik/psn-steam-api/internal/global"
	"github.com/camtrik/psn-steam-api/internal/handler"
	"github.com/camtrik/psn-steam-api/internal/pkg/logger"
	"github.com/camtrik/psn-steam-api/internal/service/psn"
	"github.com/camtrik/psn-steam-api/internal/service/steam"
	"github.com/camtrik/psn-steam-api/internal/service/xbox"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
	zapLogger, err := zap.NewProduction()
	// zapLogger, err := zap.NewDevelopment()

	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
		os.Exit(1)
	}
	defer zapLogger.Sync()

	logger := logger.NewLogger(zapLogger.Sugar())

	// testTrophyTitle()
	cfg := global.Load()

	rdb := helper.InitRedis(cfg)

	httpClient := helper.InitHttpClient()

	psnCache := cache.NewRedisPSNCache(rdb)
	psnService := psn.NewPSNService(httpClient, psnCache, logger, cfg.PSNRefreshToken)
	psnHandler := handler.NewPSNHandler(psnService)

	steamCache := cache.NewRedisSteamCache(rdb)
	steamService := steam.NewSteamService(httpClient, *steamCache, logger, cfg.SteamApiKey)
	steamHandler := handler.NewSteamHandler(steamService)

	xboxService := xbox.NewXboxService(httpClient, logger, cfg.XboxApiKey)
	xboxHandler := handler.NewXboxHandler(xboxService)

	r := gin.Default()

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// psn
	r.GET("/api/psn/:accountId/trophyTitles", psnHandler.GetUserTitles)

	// steam
	r.GET("/api/steam/:steamId/ownedGames", steamHandler.GetOwnedGames)
	r.GET("/api/steam/:steamId/playerAchievements/:appId", steamHandler.GetPlayerAchievements)
	r.GET("/api/steam/:steamId/playerGameDetails", steamHandler.GetPlayerGameDetails)
	r.GET("/api/steam/:steamId/recentlyPlayed", steamHandler.GetRecentlyPlayedGames)

	// xbox
	r.GET("/api/xbox/achievements", xboxHandler.GetPlayerAchievements)

	r.Run(":7071")

	logger.Info("Start server on port 6061")
}
