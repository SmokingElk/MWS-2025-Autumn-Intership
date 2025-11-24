package interfaces

import (
	moduleEntity "github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/domain/module/entity"
	repoEntity "github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/domain/repo/entity"
)

type RepoService interface {
	AddHubClient(name string, client RepoHubClient)
	GetRepoInfo(url string, showInderect bool) (repoEntity.Repo, []moduleEntity.Module, error)
}
