package github

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"

	repoEntity "github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/domain/repo/entity"
	repoErrors "github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/domain/repo/errors"
	"github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/domain/repo/interfaces"
	"github.com/google/go-github/v62/github"
)

const urlPrefix = "https://github.com/"
const moduleFile = "go.mod"

type RepoHubClientGithub struct {
	client *github.Client
}

func NewRepoHubClientGithub(token string) interfaces.RepoHubClient {
	res := &RepoHubClientGithub{}

	if token != "" {
		res.client = github.NewClient(nil).WithAuthToken(token)
	} else {
		res.client = github.NewClient(http.DefaultClient)
	}

	return res
}

func (r *RepoHubClientGithub) GetRepo(
	ctx context.Context,
	url string,
	builder func(gomod []byte) (repoEntity.Repo, error),
) (repoEntity.Repo, error) {
	if !strings.HasPrefix(url, urlPrefix) {
		return repoEntity.Repo{}, repoErrors.ErrBadUrl
	}

	parts := strings.Split(strings.TrimPrefix(url, urlPrefix), "/")

	if len(parts) < 2 {
		return repoEntity.Repo{}, repoErrors.ErrBadUrl
	}

	owner, repoName := parts[0], parts[1]

	_, rootContents, _, err := r.client.Repositories.GetContents(ctx, owner, repoName, "", nil)

	if err != nil {
		if e, ok := err.(*github.ErrorResponse); ok && e.Response.StatusCode == http.StatusNotFound {
			return repoEntity.Repo{}, repoErrors.ErrRepoNotFound
		}

		return repoEntity.Repo{}, fmt.Errorf("failed to get repo contents: %w", err)
	}

	containsModuleFile := false

	for _, file := range rootContents {
		if *file.Name == moduleFile {
			containsModuleFile = true
			break
		}
	}

	if !containsModuleFile {
		return repoEntity.Repo{}, repoErrors.ErrNotGoRepo
	}

	moduleFile, _, _, err := r.client.Repositories.GetContents(ctx, owner, repoName, moduleFile, nil)
	if err != nil {
		return repoEntity.Repo{}, fmt.Errorf("failed to get module file content: %w", err)
	}

	moduleFileContent, err := base64.StdEncoding.DecodeString(*moduleFile.Content)
	if err != nil {
		return repoEntity.Repo{}, repoErrors.ErrBadGomod
	}

	repo, err := builder(moduleFileContent)

	if err != nil {
		if errors.Is(err, repoErrors.ErrBadGomod) {
			return repoEntity.Repo{}, err
		}

		return repoEntity.Repo{}, fmt.Errorf("failed to build repo: %w", err)
	}

	return repo, nil
}
