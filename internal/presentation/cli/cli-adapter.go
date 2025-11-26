package cli

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/url"
	"sort"
	"strings"
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

func (a *CLIAdapter) Serve(in io.Reader, out io.Writer) int {
	if *a.flags.Help {
		flag.CommandLine.SetOutput(out)
		flag.Usage()
		return exitcodes.OK
	}

	urlStr := *a.flags.Url

	if urlStr == "" {
		a.print(out, "Input repository url: ")

		scanner := bufio.NewScanner(in)
		scanner.Scan()
		urlStr = strings.TrimSpace(scanner.Text())

		if err := scanner.Err(); err != nil {
			a.print(out, fmt.Sprintf("Reader error: %s\n", err.Error()))
			return exitcodes.UnknownError
		}
	}

	if !a.isValidUrl(urlStr) {
		a.print(out, fmt.Sprintf("Invalid url: %s\n", urlStr))
		return exitcodes.BadURL
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(a.cfg.TimeoutSeconds)*time.Second)
	defer cancel()

	repo, upd, err := a.service.GetRepoInfo(ctx, urlStr, *a.flags.ShowIndirect)

	if err != nil {
		return a.handleError(ctx, err, urlStr, out)
	}

	sort.Slice(upd, func(i, j int) bool {
		return upd[i].Name < upd[j].Name
	})

	if *a.flags.Verbose {
		a.printRepoInfoVerbose(repo, upd, out)
	} else {
		a.printRepoInfo(repo, upd, out)
	}

	return exitcodes.OK
}

func (a *CLIAdapter) printRepoInfo(repo repoEntity.Repo, upd []moduleEntity.Module, out io.Writer) {
	a.print(out, fmt.Sprintln(repo.Name))
	a.print(out, fmt.Sprintln(repo.Version))

	for _, dependency := range upd {
		a.print(out, fmt.Sprintln(dependency.Name))
	}
}

func (a *CLIAdapter) printRepoInfoVerbose(repo repoEntity.Repo, upd []moduleEntity.Module, out io.Writer) {
	a.print(
		out,
		fmt.Sprintf("\nMODULE: %s\nGO VERSION: %v\n", repo.Name, repo.Version),
	)

	a.print(out, fmt.Sprintln(sep))

	if len(upd) == 0 {
		a.print(out, fmt.Sprintln(nothingToUpdateText))
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

		a.print(
			out,
			fmt.Sprintf(
				"%-*s | CURRENT %9v | LAST %9v | %s\n",
				depNameWidth,
				dependency.Name,
				currentVersion,
				dependency.Version,
				requireType,
			),
		)
	}
}

func (a *CLIAdapter) handleError(ctx context.Context, err error, url string, out io.Writer) int {
	switch {
	case errors.Is(err, repoErrors.ErrRepoNotFound):
		a.print(out, fmt.Sprintln("Repository is private or not found"))
		return exitcodes.RepoNotFound
	case errors.Is(err, repoErrors.ErrBadUrl):
		a.print(out, fmt.Sprintf("Invalid url: %s\n", url))
		return exitcodes.BadURL
	case errors.Is(err, repoErrors.ErrUnknownHub):
		a.print(out, fmt.Sprintln("Unknown repository hub"))
		return exitcodes.UnknownHub
	case errors.Is(err, repoErrors.ErrNotGoRepo):
		a.print(out, fmt.Sprintln("Not a go repo: go.mod not found"))
		return exitcodes.NotGoRepo
	case errors.Is(err, repoErrors.ErrBadGomod):
		a.print(out, fmt.Sprintln("Failed to parse go.mod"))
		return exitcodes.BadGoMod
	case errors.Is(ctx.Err(), context.DeadlineExceeded):
		a.print(out, fmt.Sprintln("Timeout exceeded"))
		return exitcodes.TimeoutExceeded
	default:
		a.print(out, fmt.Sprintf("An error occurred while getting repository info: %s\n", err.Error()))
		return exitcodes.UnknownError
	}
}

func (a *CLIAdapter) isValidUrl(urlStr string) bool {
	u, err := url.Parse(urlStr)
	return err == nil && u.Scheme != "" && u.Hostname() != ""
}

func (a *CLIAdapter) print(out io.Writer, msg string) {
	if _, err := out.Write([]byte(msg)); err != nil {
		log.Fatal("failed to write to out")
	}
}
