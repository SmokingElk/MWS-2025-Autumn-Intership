package di

import (
	reposervice "github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/application/repo"
	"github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/config"
	proxygolang "github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/infrastructure/module-clients/proxy-golang"
	"github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/infrastructure/repo-hub-clients/github"
)

func MustConfigureApp(cfg *config.Config) {
	moduleClient := proxygolang.NewModuleProxyGolang()

	githubClient := github.NewRepoHubClientGithub(cfg.GithubConfig.AuthToken)

	repoService := reposervice.NewRepoService(moduleClient)
	repoService.AddHubClient("github.com", githubClient)
}
