package interfaces

import (
	"context"

	"github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/domain/repo/entity"
)

type RepoHubClient interface {
	GetRepo(ctx context.Context, url string, builder func(gomod []byte) (entity.Repo, error)) (entity.Repo, error)
}
