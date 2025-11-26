package interfaces

import (
	"context"

	moduleEntity "github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/domain/module/entity"
	repoEntity "github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/domain/repo/entity"
)

type RepoService interface {
	AddHubClient(host string, client RepoHubClient)
	GetRepoInfo(ctx context.Context, repoUrl string, showIndirect bool) (repoEntity.Repo, []moduleEntity.Module, error)
}
