package psn

import (
	"context"
	"time"

	"github.com/camtrik/psn-steam-api/internal/models"
)

func (s *PSNService) GetRecentlyPlayedGames(ctx context.Context, accountId string, timeRange int64) ([]models.RecentlyPlayedGame, error) {
	userTitles, err := s.GetUserTitles(ctx, accountId, nil, models.TrophyFilter{})
	if err != nil {
		s.logger.Errorf("failed to get user titles: %v", err)
		return nil, err
	}

	timeStamp := time.Now().Unix() - timeRange

	titles := userTitles.TrophyTitles

	games := []models.RecentlyPlayedGame{}

	for _, title := range titles {
		lastPlayedTime := title.LastUpdatedDateTime.Unix()
		if lastPlayedTime < timeStamp {
			continue
		}

		earnedAchievements := title.DefinedTrophies.Bronze + title.DefinedTrophies.Silver + title.DefinedTrophies.Gold + title.DefinedTrophies.Platinum
		games = append(games, models.RecentlyPlayedGame{
			Name:               title.TrophyTitleName,
			LastPlayedTime:     lastPlayedTime,
			EarnedAchievements: earnedAchievements,
			ArtUrl:             title.TrophyTitleIconUrl,
			Platform:           "psn",
		})
	}

	return games, nil
}
