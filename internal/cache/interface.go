package cache

import (
	"context"

	"github.com/camtrik/ebbilogue-backend/internal/models"
)

type PSNCache interface {
	GetUserTitles(ctx context.Context, accountId string) (*models.UserTitlesResponse, error)
	SetUserTitles(ctx context.Context, accountId string, userTitles *models.UserTitlesResponse) error
}
