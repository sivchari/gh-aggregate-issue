package gh

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/cli/go-gh/v2"
)

func PullRequests(ctx context.Context) ([]*PullRequest, error) {
	stdout, _, err := gh.ExecContext(ctx, "pr", "ls", "-s", "all", "--json", "number,title,url", "--limit", "50")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pull requests: %w", err)
	}
	var prs []*PullRequest
	if err := json.Unmarshal(stdout.Bytes(), &prs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal pull requests: %w", err)
	}
	return prs, nil
}

type PullRequest struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	URL    string `json:"url"`
}

func Issues(ctx context.Context) ([]*Issue, error) {
	stdout, _, err := gh.ExecContext(ctx, "issue", "ls", "-s", "all", "--json", "number,title,url", "--limit", "50")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch issues: %w", err)
	}
	var issues []*Issue
	if err := json.Unmarshal(stdout.Bytes(), &issues); err != nil {
		return nil, fmt.Errorf("failed to unmarshal issues: %w", err)
	}
	return issues, nil
}

type Issue struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	URL    string `json:"url"`
}

func CreateIssue(ctx context.Context, title string, prs []*PullRequest, issues []*Issue) error {
	str := "This issue aggregates the following pull requests and issues by gh aggregate-issue:\n\n"
	if len(prs) > 0 {
		sort.SliceStable(prs, func(i, j int) bool {
			return prs[i].Number < prs[j].Number
		})
		str += "## Pull Requests:\n\n"
		for _, p := range prs {
			str += fmt.Sprintf("  - #%d\n", p.Number)
		}
		str += "\n"
	}
	if len(issues) > 0 {
		sort.SliceStable(issues, func(i, j int) bool {
			return issues[i].Number < issues[j].Number
		})
		str += "## Issues:\n\n"
		for _, i := range issues {
			str += fmt.Sprintf("  - #%d\n", i.Number)
		}
		str += "\n"
	}
	if _, _, err := gh.ExecContext(ctx, "issue", "create", "--title", title, "--body", str, "--assignee", "@me"); err != nil {
		return fmt.Errorf("failed to create issue: %w", err)
	}
	return nil
}
