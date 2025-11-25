package reposervice

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"sync"

	moduleEntity "github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/domain/module/entity"
	moduleInterfaces "github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/domain/module/interfaces"
	repoEntity "github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/domain/repo/entity"
	repoErrors "github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/domain/repo/errors"
	repoInterfaces "github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/domain/repo/interfaces"
	versionEntity "github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/domain/version/entity"
	"golang.org/x/mod/modfile"
)

type RepoService struct {
	moduleClient   moduleInterfaces.ModuleClient
	repoHubClients map[string]repoInterfaces.RepoHubClient
}

func NewRepoService(moduleClient moduleInterfaces.ModuleClient) repoInterfaces.RepoService {
	return &RepoService{
		moduleClient:   moduleClient,
		repoHubClients: make(map[string]repoInterfaces.RepoHubClient),
	}
}

func (s *RepoService) AddHubClient(host string, client repoInterfaces.RepoHubClient) {
	s.repoHubClients[host] = client
}

func (s *RepoService) GetRepoInfo(
	ctx context.Context,
	repoUrl string,
	showInderect bool,
) (repoEntity.Repo, []moduleEntity.Module, error) {
	u, err := url.Parse(repoUrl)

	if err != nil {
		return repoEntity.Repo{}, nil, repoErrors.ErrBadUrl
	}

	hubHost := strings.Split(strings.TrimPrefix(u.Hostname(), "www."), ":")[0]

	hubClient, ok := s.repoHubClients[hubHost]

	if !ok {
		return repoEntity.Repo{}, nil, repoErrors.ErrUnknownHub
	}

	repo, err := hubClient.GetRepo(ctx, repoUrl, s.builder(showInderect))

	if err != nil {
		if errors.Is(err, repoErrors.ErrBadUrl) ||
			errors.Is(err, repoErrors.ErrNotGoRepo) ||
			errors.Is(err, repoErrors.ErrBadGomod) ||
			errors.Is(err, repoErrors.ErrRepoNotFound) {

			return repoEntity.Repo{}, nil, err
		}

		return repoEntity.Repo{}, nil, fmt.Errorf("failed to get repo from hub: %w", err)
	}

	dependenciesToUpdate, err := s.getDependenciesToUpdate(ctx, repo)

	if err != nil {
		return repoEntity.Repo{}, nil, fmt.Errorf("failed to get dependencies to update: %w", err)
	}

	return repo, dependenciesToUpdate, nil
}

func (s *RepoService) builder(showInderect bool) func([]byte) (repoEntity.Repo, error) {
	return func(content []byte) (repoEntity.Repo, error) {
		file, err := modfile.Parse("go.mod", content, nil)

		if err != nil {
			return repoEntity.Repo{}, repoErrors.ErrBadGomod
		}

		repo, err := repoEntity.NewRepo(file.Module.Mod.Path, file.Go.Version)

		if err != nil {
			return repoEntity.Repo{}, repoErrors.ErrBadGomod
		}

		for _, dependency := range file.Require {
			if !showInderect && dependency.Indirect {
				continue
			}

			version, err := versionEntity.NewVersion(dependency.Mod.Version)

			if err != nil {
				return repoEntity.Repo{}, repoErrors.ErrBadGomod
			}

			repo.AddDependency(moduleEntity.NewModule(dependency.Mod.Path, version, !dependency.Indirect))
		}

		return repo, nil
	}
}

func (s *RepoService) getDependenciesToUpdate(ctx context.Context, repo repoEntity.Repo) ([]moduleEntity.Module, error) {
	mtx := sync.Mutex{}
	res := make([]moduleEntity.Module, 0)

	wg := sync.WaitGroup{}

	for _, dependency := range repo.Dependencies {
		wg.Add(1)

		go func(dependency moduleEntity.Module) {
			defer wg.Done()

			lastVersion, err := s.moduleClient.GetLastVersion(ctx, dependency)

			if err != nil {
				// TODO: add warning
				return
			}

			if dependency.Version.Less(lastVersion) {
				actualDependency := moduleEntity.NewModule(
					dependency.Name,
					lastVersion,
					dependency.Direct,
				)

				mtx.Lock()
				defer mtx.Unlock()
				res = append(res, actualDependency)
			}
		}(dependency)
	}

	wg.Wait()

	return res, nil
}
