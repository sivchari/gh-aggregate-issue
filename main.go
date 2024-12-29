package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/sivchari/gh-aggregate-issue/gh"
	"github.com/sivchari/gh-aggregate-issue/tui"
)

func init() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, nil)))
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	prs, err := gh.PullRequests(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to fetch pull requests", slog.String("error", err.Error()))
		return
	}

	issues, err := gh.Issues(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to fetch issues", slog.String("error", err.Error()))
		return
	}

	tm := tui.NewTitle()
	if _, err := tea.NewProgram(tm).Run(); err != nil {
		slog.ErrorContext(ctx, "failed to run TUI", slog.String("error", err.Error()))
		return
	}

	prm := tui.NewPRModel(prs)
	if _, err := tea.NewProgram(prm).Run(); err != nil {
		slog.ErrorContext(ctx, "failed to run TUI", slog.String("error", err.Error()))
		return
	}

	im := tui.NewIssueModel(issues)
	if _, err := tea.NewProgram(im).Run(); err != nil {
		slog.ErrorContext(ctx, "failed to run TUI", slog.String("error", err.Error()))
		return
	}

	if err := gh.CreateIssue(ctx, tm.Value(), prm.Value(), im.Value()); err != nil {
		slog.ErrorContext(ctx, "failed to create issue", slog.String("error", err.Error()))
		return
	}
}
