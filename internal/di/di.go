package di

import (
	reposervice "github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/application/repo"
	"github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/config"
	"github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/di/app"
	"github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/di/flags"
	proxygolang "github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/infrastructure/module-clients/proxy-golang"
	"github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/infrastructure/repo-hub-clients/github"
	"github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/presentation/cli"
)

func MustConfigureApp(flags *flags.Flags, cfg *config.Config) app.App {
	moduleClient := proxygolang.NewModuleProxyGolang()

	githubClient := github.NewRepoHubClientGithub(cfg.GithubConfig.AuthToken)

	repoService := reposervice.NewRepoService(moduleClient)
	repoService.AddHubClient("github.com", githubClient)

	app := cli.NewCLIAdapter(repoService, flags, &cfg.CLIConfig)

	return app
}
