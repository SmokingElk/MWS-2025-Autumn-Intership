package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"sort"
	"time"

	"github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/config"
	"github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/di/app"
	"github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/di/flags"
	moduleEntity "github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/domain/module/entity"
	repoEntity "github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/domain/repo/entity"
	repoErrors "github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/domain/repo/errors"
	"github.com/SmokingElk/MWS-2025-Autumn-Intership/internal/domain/repo/interfaces"
	exitcodes "github.com/SmokingElk/MWS-2025-Autumn-Intership/pkg/exit-codes"
)

const sep = "=========="
const nothingToUpdateText = "Nothing to update, dependency list is clean"

type CLIAdapter struct {
	service interfaces.RepoService
	flags   *flags.Flags
	cfg     *config.CLIConfig
}

func NewCLIAdapter(service interfaces.RepoService, flags *flags.Flags, cfg *config.CLIConfig) app.App {
	return &CLIAdapter{
		flags:   flags,
		cfg:     cfg,
		service: service,
	}
}

func (a *CLIAdapter) Serve() error {
	if *a.flags.Help {
		flag.Usage()
		return nil
	}

	urlStr := *a.flags.Url

	if urlStr == "" {
		fmt.Print("Input repository url: ")
		fmt.Scan(&urlStr)
	}

	if !a.isValidUrl(urlStr) {
		fmt.Printf("Invalid url: %s", urlStr)
		os.Exit(exitcodes.BadURL)
		return repoErrors.ErrBadUrl
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(a.cfg.TimeoutSeconds)*time.Second)
	defer cancel()

	repo, upd, err := a.service.GetRepoInfo(ctx, urlStr, *a.flags.ShowIndirect)

	if err != nil {
		a.handleError(ctx, err, urlStr)
		return err
	}

	sort.Slice(upd, func(i, j int) bool {
		return upd[i].Name < upd[j].Name
	})

	if *a.flags.Verbose {
		a.printRepoInfoVerbose(repo, upd)
	} else {
		a.printRepoInfo(repo, upd)
	}

	return nil
}

func (a *CLIAdapter) printRepoInfo(repo repoEntity.Repo, upd []moduleEntity.Module) {
	fmt.Println(repo.Name)
	fmt.Println(repo.Version)

	for _, dependency := range upd {
		fmt.Println(dependency.Name)
	}
}

func (a *CLIAdapter) printRepoInfoVerbose(repo repoEntity.Repo, upd []moduleEntity.Module) {
	fmt.Printf(
		"\nMODULE: %s\nGO VERSION: %v\n",
		repo.Name,
		repo.Version,
	)

	fmt.Println(sep)

	if len(upd) == 0 {
		fmt.Println(nothingToUpdateText)
		return
	}

	depNameWidth := 0
	for _, dependency := range upd {
		depNameWidth = max(depNameWidth, len(dependency.Name))
	}

	for _, dependency := range upd {
		requireType := "DIRECT"

		if !dependency.Direct {
			requireType = "INDIRECT"
		}

		currentVersion := repo.Dependencies[dependency.Name].Version

		fmt.Printf(
			"%-*s | CURRENT %9v | LAST %9v | %s\n",
			depNameWidth,
			dependency.Name,
			currentVersion,
			dependency.Version,
			requireType,
		)
	}
}

func (a *CLIAdapter) handleError(ctx context.Context, err error, url string) {
	switch {
	case errors.Is(err, repoErrors.ErrRepoNotFound):
		fmt.Println("Repository is private or not found")
		os.Exit(exitcodes.RepoNotFound)
	case errors.Is(err, repoErrors.ErrBadUrl):
		fmt.Printf("Invalid url: %s", url)
		os.Exit(exitcodes.BadURL)
	case errors.Is(err, repoErrors.ErrUnknownHub):
		fmt.Println("Unknown repository hub")
		os.Exit(exitcodes.UnknownHub)
	case errors.Is(err, repoErrors.ErrNotGoRepo):
		fmt.Println("Not a go repo: go.mod not found")
		os.Exit(exitcodes.NotGoRepo)
	case errors.Is(err, repoErrors.ErrBadGomod):
		fmt.Println("Failed to parse go.mod")
		os.Exit(exitcodes.BadGoMod)
	case errors.Is(ctx.Err(), context.DeadlineExceeded):
		fmt.Println("Timeout exceeded")
		os.Exit(exitcodes.TimeoutExceeded)
	default:
		fmt.Printf("An error occured while getting repository info: %s\n", err.Error())
		os.Exit(exitcodes.UnknownError)
	}
}

func (a *CLIAdapter) isValidUrl(urlStr string) bool {
	u, err := url.Parse(urlStr)
	return err == nil && u.Scheme != "" && u.Hostname() != ""
}
