package reposervice_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"testing"

	reposervice "github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/application/repo"
	moduleEntity "github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/domain/module/entity"
	moduleErrors "github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/domain/module/errors"
	moduleMocks "github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/domain/module/mocks"
	repoEntity "github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/domain/repo/entity"
	repoErrors "github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/domain/repo/errors"
	repoMocks "github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/domain/repo/mocks"
	versionEntity "github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/domain/version/entity"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func version(str string) versionEntity.Version {
	version, err := versionEntity.NewVersion(str)

	if err != nil {
		log.Fatal(err)
	}

	return version
}

func TestGetRepoInfo(t *testing.T) {

	type testCase struct {
		what string

		url          string
		showIndirect bool

		gomod             string
		githubClientError error

		lastVersions       map[string]versionEntity.Version
		moduleClientErrors map[string]error

		expectedError error
		expectedUpd   []moduleEntity.Module
		expectedRepo  repoEntity.Repo
	}

	testCases := []testCase{
		{
			what:          "bad url",
			url:           "https://git hub.com/Bob/module",
			expectedError: repoErrors.ErrBadUrl,
		},

		{
			what:          "unknown hub",
			url:           "https://gitlab.com/Bob/module",
			expectedError: repoErrors.ErrUnknownHub,
		},

		{
			what:              "client url parsing fail",
			url:               "https://github.com/Bob",
			githubClientError: repoErrors.ErrBadUrl,
			expectedError:     repoErrors.ErrBadUrl,
		},

		{
			what:              "not go repo",
			url:               "https://github.com/Bob/python-project",
			githubClientError: repoErrors.ErrNotGoRepo,
			expectedError:     repoErrors.ErrNotGoRepo,
		},

		{
			what:              "repo not found",
			url:               "https://github.com/Bob/not-existing-repo",
			githubClientError: repoErrors.ErrRepoNotFound,
			expectedError:     repoErrors.ErrRepoNotFound,
		},

		{
			what:              "corrupted go.mod",
			url:               "https://github.com/Bob/go-repo",
			githubClientError: repoErrors.ErrBadGomod,
			gomod:             "lolkek",
			expectedError:     repoErrors.ErrBadGomod,
		},

		{
			what:              "unexpected client error",
			url:               "https://github.com/Bob/go-repo",
			githubClientError: errors.New("failed to get repo"),
			expectedError:     errors.New("failed to get repo from hub: failed to get repo"),
		},

		{
			what:         "ignore indirect",
			url:          "https://github.com/Bob/go-repo",
			showIndirect: false,
			gomod: `
module go-repo
go 1.24.0

require (
	github.com/mod1 v1.0.1
	github.com/mod2 v1.0.2
)

require (
	github.com/mod3 v1.0.3 // indirect
	github.com/mod4 v1.0.4 // indirect
)
			`,
			githubClientError: nil,
			expectedRepo: repoEntity.Repo{
				Name:    "go-repo",
				Version: version("1.24.0"),
				Dependencies: map[string]moduleEntity.Module{
					"github.com/mod1": {
						Name:    "github.com/mod1",
						Version: version("v1.0.1"),
						Direct:  true,
					},
					"github.com/mod2": {
						Name:    "github.com/mod2",
						Version: version("v1.0.2"),
						Direct:  true,
					},
				},
			},
			lastVersions: map[string]versionEntity.Version{
				"github.com/mod1": version("v1.0.1"),
				"github.com/mod2": version("v1.0.2"),
			},
			moduleClientErrors: map[string]error{
				"github.com/mod1": nil,
				"github.com/mod2": nil,
			},
			expectedError: nil,
			expectedUpd:   []moduleEntity.Module{},
		},

		{
			what:         "show indirect",
			url:          "https://github.com/Bob/go-repo",
			showIndirect: true,
			gomod: `
module go-repo
go 1.24.0

require (
	github.com/mod1 v1.0.1
	github.com/mod2 v1.0.2
)

require (
	github.com/mod3 v1.0.3 // indirect
	github.com/mod4 v1.0.4 // indirect
)
			`,
			githubClientError: nil,
			expectedRepo: repoEntity.Repo{
				Name:    "go-repo",
				Version: version("1.24.0"),
				Dependencies: map[string]moduleEntity.Module{
					"github.com/mod1": {
						Name:    "github.com/mod1",
						Version: version("v1.0.1"),
						Direct:  true,
					},
					"github.com/mod2": {
						Name:    "github.com/mod2",
						Version: version("v1.0.2"),
						Direct:  true,
					},
					"github.com/mod3": {
						Name:    "github.com/mod3",
						Version: version("v1.0.3"),
						Direct:  false,
					},
					"github.com/mod4": {
						Name:    "github.com/mod4",
						Version: version("v1.0.4"),
						Direct:  false,
					},
				},
			},
			lastVersions: map[string]versionEntity.Version{
				"github.com/mod1": version("v1.0.1"),
				"github.com/mod2": version("v1.0.2"),
				"github.com/mod3": version("v1.0.3"),
				"github.com/mod4": version("v1.0.4"),
			},
			moduleClientErrors: map[string]error{
				"github.com/mod1": nil,
				"github.com/mod2": nil,
				"github.com/mod3": nil,
				"github.com/mod4": nil,
			},
			expectedError: nil,
			expectedUpd:   []moduleEntity.Module{},
		},

		{
			what:         "to update",
			url:          "https://github.com/Bob/go-repo",
			showIndirect: true,
			gomod: `
module go-repo
go 1.24.0

require (
	github.com/mod1 v1.0.1
	github.com/mod2 v1.0.2
)

require (
	github.com/mod3 v1.0.3 // indirect
	github.com/mod4 v1.0.4 // indirect
)
			`,
			githubClientError: nil,
			expectedRepo: repoEntity.Repo{
				Name:    "go-repo",
				Version: version("1.24.0"),
				Dependencies: map[string]moduleEntity.Module{
					"github.com/mod1": {
						Name:    "github.com/mod1",
						Version: version("v1.0.1"),
						Direct:  true,
					},
					"github.com/mod2": {
						Name:    "github.com/mod2",
						Version: version("v1.0.2"),
						Direct:  true,
					},
					"github.com/mod3": {
						Name:    "github.com/mod3",
						Version: version("v1.0.3"),
						Direct:  false,
					},
					"github.com/mod4": {
						Name:    "github.com/mod4",
						Version: version("v1.0.4"),
						Direct:  false,
					},
				},
			},
			lastVersions: map[string]versionEntity.Version{
				"github.com/mod1": version("v1.0.1"),
				"github.com/mod2": version("v1.1.2"),
				"github.com/mod3": version("v1.0.3"),
				"github.com/mod4": version("v1.1.4"),
			},
			moduleClientErrors: map[string]error{
				"github.com/mod1": nil,
				"github.com/mod2": nil,
				"github.com/mod3": nil,
				"github.com/mod4": nil,
			},
			expectedError: nil,
			expectedUpd: []moduleEntity.Module{
				{
					Name:    "github.com/mod2",
					Version: version("v1.1.2"),
					Direct:  true,
				},
				{
					Name:    "github.com/mod4",
					Version: version("v1.1.4"),
					Direct:  false,
				},
			},
		},

		{
			what:         "filter last version errors",
			url:          "https://github.com/Bob/go-repo",
			showIndirect: true,
			gomod: `
module go-repo
go 1.24.0

require (
	github.com/mod1 v1.0.1
	github.com/mod2 v1.0.2
)

require (
	github.com/mod3 v1.0.3 // indirect
	github.com/mod4 v1.0.4 // indirect
)
			`,
			githubClientError: nil,
			expectedRepo: repoEntity.Repo{
				Name:    "go-repo",
				Version: version("1.24.0"),
				Dependencies: map[string]moduleEntity.Module{
					"github.com/mod1": {
						Name:    "github.com/mod1",
						Version: version("v1.0.1"),
						Direct:  true,
					},
					"github.com/mod2": {
						Name:    "github.com/mod2",
						Version: version("v1.0.2"),
						Direct:  true,
					},
					"github.com/mod3": {
						Name:    "github.com/mod3",
						Version: version("v1.0.3"),
						Direct:  false,
					},
					"github.com/mod4": {
						Name:    "github.com/mod4",
						Version: version("v1.0.4"),
						Direct:  false,
					},
				},
			},
			lastVersions: map[string]versionEntity.Version{
				"github.com/mod1": version("v1.0.1"),
				"github.com/mod2": version("v1.1.2"),
				"github.com/mod3": version("v1.0.3"),
				"github.com/mod4": version("v1.1.4"),
			},
			moduleClientErrors: map[string]error{
				"github.com/mod1": nil,
				"github.com/mod2": nil,
				"github.com/mod3": nil,
				"github.com/mod4": moduleErrors.ErrModuleNotFound,
			},
			expectedError: nil,
			expectedUpd: []moduleEntity.Module{
				{
					Name:    "github.com/mod2",
					Version: version("v1.1.2"),
					Direct:  true,
				},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("Test %d: %s", i, tc.what), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockGithubClient := repoMocks.NewMockRepoHubClient(ctrl)

			mockGithubClient.EXPECT().GetRepo(
				gomock.Any(),
				tc.url,
				gomock.Any(),
			).DoAndReturn(func(
				ctx context.Context,
				url string,
				builder func(content []byte) (repoEntity.Repo, error),
			) (repoEntity.Repo, error) {
				if tc.githubClientError != nil && !errors.Is(tc.githubClientError, repoErrors.ErrBadGomod) {
					return tc.expectedRepo, tc.githubClientError
				}

				repo, err := builder([]byte(tc.gomod))

				if tc.githubClientError == nil {
					assert.NoError(t, err)
					assert.Equal(t, tc.expectedRepo, repo)
				} else {
					assert.EqualError(t, err, tc.githubClientError.Error())
				}

				return tc.expectedRepo, tc.githubClientError
			}).MaxTimes(1)

			mockModuleClient := moduleMocks.NewMockModuleClient(ctrl)

			for _, module := range tc.expectedRepo.Dependencies {
				mockModuleClient.EXPECT().GetLastVersion(
					gomock.Any(),
					module,
				).Return(tc.lastVersions[module.Name], tc.moduleClientErrors[module.Name]).MaxTimes(1)
			}

			repoService := reposervice.NewRepoService(mockModuleClient)
			repoService.AddHubClient("github.com", mockGithubClient)

			repo, upd, err := repoService.GetRepoInfo(context.Background(), tc.url, tc.showIndirect)

			if tc.expectedError == nil {
				assert.NoError(t, err)
				assert.ElementsMatch(t, tc.expectedUpd, upd)
				assert.Equal(t, tc.expectedRepo, repo)
			} else {
				assert.EqualError(t, err, tc.expectedError.Error())
			}
		})
	}
}
