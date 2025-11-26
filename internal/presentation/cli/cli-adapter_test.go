package cli_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	reposervice "github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/application/repo"
	"github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/config"
	"github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/di/flags"
	moduleEntity "github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/domain/module/entity"
	moduleMocks "github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/domain/module/mocks"
	repoEntity "github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/domain/repo/entity"
	repoErrors "github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/domain/repo/errors"
	repoMocks "github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/domain/repo/mocks"
	versionEntity "github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/domain/version/entity"
	"github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/presentation/cli"
	exitcodes "github.com/SmokingElk/MWS-2025-Autumn-Intership/pkg/exit-codes"
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

func TestServe(t *testing.T) {
	type testCase struct {
		what string

		showIndirect bool
		verbose      bool
		url          string

		expectedRepo      repoEntity.Repo
		githubClientError error

		lastVersions       map[string]versionEntity.Version
		moduleClientErrors map[string]error

		input            string
		expectedOutput   string
		expectedExitCode int
	}

	testCases := []testCase{
		{
			what:             "bad url",
			url:              "///",
			expectedExitCode: exitcodes.BadURL,
			expectedOutput:   "Invalid url: ///\n",
		},

		{
			what:              "bad url in service",
			url:               "https://github.com/Bob",
			githubClientError: repoErrors.ErrBadUrl,
			expectedExitCode:  exitcodes.BadURL,
			expectedOutput:    "Invalid url: https://github.com/Bob\n",
		},

		{
			what:             "unknown hub",
			url:              "https://gitlab.com/Bob/module",
			expectedExitCode: exitcodes.UnknownHub,
			expectedOutput:   "Unknown repository hub\n",
		},

		{
			what:              "not go repo",
			url:               "https://github.com/Bob/python-project",
			githubClientError: repoErrors.ErrNotGoRepo,
			expectedExitCode:  exitcodes.NotGoRepo,
			expectedOutput:    "Not a go repo: go.mod not found\n",
		},

		{
			what:              "repo not found",
			url:               "https://github.com/Bob/not-existing-repo",
			githubClientError: repoErrors.ErrRepoNotFound,
			expectedExitCode:  exitcodes.RepoNotFound,
			expectedOutput:    "Repository is private or not found\n",
		},

		{
			what:              "bad gomod",
			url:               "https://github.com/Bob/not-existing-repo",
			githubClientError: repoErrors.ErrBadGomod,
			expectedExitCode:  exitcodes.BadGoMod,
			expectedOutput:    "Failed to parse go.mod\n",
		},

		{
			what:              "unknown error",
			url:               "https://github.com/Bob/repo",
			githubClientError: errors.New("failed to get repo"),
			expectedExitCode:  exitcodes.UnknownError,
			expectedOutput:    "An error occured while getting repository info: failed to get repo from hub: failed to get repo\n",
		},

		{
			what: "default output",
			url:  "https://github.com/Bob/go-repo",
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
				"github.com/mod1": version("v1.1.1"),
				"github.com/mod2": version("v1.1.2"),
			},
			moduleClientErrors: map[string]error{
				"github.com/mod1": nil,
				"github.com/mod2": nil,
			},
			expectedExitCode: 0,
			expectedOutput:   "go-repo\nv1.24.0\ngithub.com/mod1\ngithub.com/mod2\n",
		},

		{
			what:    "verbose output",
			url:     "https://github.com/Bob/go-repo",
			verbose: true,
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
				"github.com/mod1": version("v1.1.1"),
				"github.com/mod2": version("v1.1.2"),
			},
			moduleClientErrors: map[string]error{
				"github.com/mod1": nil,
				"github.com/mod2": nil,
			},
			expectedExitCode: 0,
			expectedOutput: `
MODULE: go-repo
GO VERSION: v1.24.0
==========
github.com/mod1 | CURRENT    v1.0.1 | LAST    v1.1.1 | DIRECT
github.com/mod2 | CURRENT    v1.0.2 | LAST    v1.1.2 | DIRECT
`,
		},

		{
			what:         "verbose with inderect",
			url:          "https://github.com/Bob/go-repo",
			verbose:      true,
			showIndirect: true,
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
			expectedExitCode: 0,
			expectedOutput: `
MODULE: go-repo
GO VERSION: v1.24.0
==========
github.com/mod2 | CURRENT    v1.0.2 | LAST    v1.1.2 | DIRECT
github.com/mod4 | CURRENT    v1.0.4 | LAST    v1.1.4 | INDIRECT
`,
		},

		{
			what:    "verbose output",
			url:     "https://github.com/Bob/go-repo",
			verbose: true,
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
			expectedExitCode: 0,
			expectedOutput: `
MODULE: go-repo
GO VERSION: v1.24.0
==========
Nothing to update, dependency list is clean
`,
		},

		{
			what:  "user input",
			url:   "",
			input: "https://github.com/Bob/go-repo\n",
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
				"github.com/mod1": version("v1.1.1"),
				"github.com/mod2": version("v1.1.2"),
			},
			moduleClientErrors: map[string]error{
				"github.com/mod1": nil,
				"github.com/mod2": nil,
			},
			expectedExitCode: 0,
			expectedOutput:   "Input repository url: go-repo\nv1.24.0\ngithub.com/mod1\ngithub.com/mod2\n",
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("Test %d: %s", i, tc.what), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockGithubClient := repoMocks.NewMockRepoHubClient(ctrl)

			url := tc.url
			if url == "" {
				url = strings.TrimSpace(tc.input)
			}

			mockGithubClient.EXPECT().GetRepo(
				gomock.Any(),
				url,
				gomock.Any(),
			).Return(tc.expectedRepo, tc.githubClientError).MaxTimes(1)

			mockModuleClient := moduleMocks.NewMockModuleClient(ctrl)

			for _, module := range tc.expectedRepo.Dependencies {
				mockModuleClient.EXPECT().GetLastVersion(
					gomock.Any(),
					module,
				).Return(tc.lastVersions[module.Name], tc.moduleClientErrors[module.Name]).MaxTimes(1)
			}

			repoService := reposervice.NewRepoService(mockModuleClient)
			repoService.AddHubClient("github.com", mockGithubClient)

			help := false
			mockFlags := flags.Flags{
				ShowIndirect: &tc.showIndirect,
				Verbose:      &tc.verbose,
				Help:         &help,
				Url:          &tc.url,
			}

			cfg := config.CLIConfig{
				TimeoutSeconds: 10,
			}

			cliAdapter := cli.NewCLIAdapter(repoService, &mockFlags, &cfg)

			in := strings.NewReader(tc.input)
			out := bytes.Buffer{}

			exitCode := cliAdapter.Serve(in, &out)

			assert.Equal(t, tc.expectedExitCode, exitCode)
			assert.Equal(t, tc.expectedOutput, out.String())
		})
	}
}

func TestServe_Timeout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGithubClient := repoMocks.NewMockRepoHubClient(ctrl)

	help := false
	url := "http://github.com/Bob/go-repo"
	verbose := false
	showIndirect := false

	expectedExitCode := exitcodes.TimeoutExceeded
	expectedOutput := "Timeout exceeded\n"

	mockGithubClient.EXPECT().GetRepo(
		gomock.Any(),
		url,
		gomock.Any(),
	).DoAndReturn(func(
		ctx context.Context,
		url string,
		builder func(content []byte) (repoEntity.Repo, error),
	) (repoEntity.Repo, error) {
		var err error
		select {
		case <-time.After(time.Second * 10):
		case <-ctx.Done():
			err = context.DeadlineExceeded
		}

		return repoEntity.Repo{}, err
	})

	mockModuleClient := moduleMocks.NewMockModuleClient(ctrl)

	repoService := reposervice.NewRepoService(mockModuleClient)
	repoService.AddHubClient("github.com", mockGithubClient)

	mockFlags := flags.Flags{
		ShowIndirect: &showIndirect,
		Verbose:      &verbose,
		Help:         &help,
		Url:          &url,
	}

	cfg := config.CLIConfig{
		TimeoutSeconds: 1,
	}

	cliAdapter := cli.NewCLIAdapter(repoService, &mockFlags, &cfg)

	in := strings.NewReader("")
	out := bytes.Buffer{}

	exitCode := cliAdapter.Serve(in, &out)

	assert.Equal(t, expectedExitCode, exitCode)
	assert.Equal(t, expectedOutput, out.String())
}
