package xbox

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/camtrik/psn-steam-api/internal/models"
)

func (s *XboxService) GetRecentlyPlayedGames(ctx context.Context, timeRange int64) ([]models.RecentlyPlayedGame, error) {
	achievements, err := s.GetPlayerAchievements(ctx)
	if err != nil {
		s.logger.Error("Failed to get player achievements %v", err)
		return nil, err
	}

	titles := achievements.Titles
	timeStamp := time.Now().Unix() - timeRange

	// concurrency
	var (
		mu    sync.Mutex
		wg    sync.WaitGroup
		games []models.RecentlyPlayedGame
	)

	for _, title := range titles {
		lastPlayedTime := title.TitleHistory.LastTimePlayed.Unix()
		if lastPlayedTime < timeStamp {
			continue
		}

		wg.Add(1)
		// 匿名函数，创建后立即用title调用
		go func(title models.XboxGameTitles) {
			defer wg.Done()

			stats, err := s.GetGameStats(ctx, title.TitleId)
			if err != nil {
				s.logger.Error("Failed to get game stats %v", err)
				return
			}

			minutesPlayed := 0
			for _, stat := range stats.StatListsCollection[0].Stats {
				if stat.Name == "MinutesPlayed" {
					minutesPlayed, err = strconv.Atoi(stat.Value)
					if err != nil {
						s.logger.Error("Failed to convert minutes played to int %v", err)
						return
					}
				}
			}

			mu.Lock()
			games = append(games, models.RecentlyPlayedGame{
				Name:               title.Name,
				PlayTime:           minutesPlayed,
				LastPlayedTime:     lastPlayedTime,
				EarnedAchievements: title.Achievement.CurrentAchievements,
				ArtUrl:             title.DisplayImage,
				Platform:           "xbox",
			})
			mu.Unlock()
		}(title)
	}
	wg.Wait()

	return games, nil
}
