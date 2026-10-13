package mock

import (
	"context"
	"fmt"

	"appletini/gitter"
)

// GetAll returns fixed example data. The client is ignored, but cancellation
// is still honoured so the mock behaves like a real query.
func (MockQuery) GetAll(ctx context.Context, _ gitter.Querier) (map[string][]gitter.PullRequest, error) {
	err := ctx.Err()
	if err != nil {
		return map[string][]gitter.PullRequest{}, fmt.Errorf("mock query: %w", err)
	}

	prMap := map[string][]gitter.PullRequest{}

	//exhaustruct:ignore
	prMap["personal"] = []gitter.PullRequest{
		{
			Number:         1,
			Title:          "Example PR 1",
			HeadRefName:    "from-branch",
			BaseRefName:    "to-branch",
			Mergeable:      "MERGEABLE",
			ReviewDecision: "APPROVED",
		},
		{
			Number:         2,
			Title:          "Example PR 2",
			HeadRefName:    "from-branch",
			BaseRefName:    "to-branch",
			Mergeable:      "UNKNOWN",
			ReviewDecision: "APPROVED",
		},
		{
			Number:         3,
			Title:          "Example PR 3",
			HeadRefName:    "from-branch",
			BaseRefName:    "to-branch",
			Mergeable:      "CONFLICTING",
			ReviewDecision: "APPROVED",
		},
		{
			Number:         4,
			Title:          "Example PR 4",
			HeadRefName:    "from-branch",
			BaseRefName:    "to-branch",
			Mergeable:      "MERGEABLE",
			ReviewDecision: "CHANGES_REQUESTED",
		},
		{
			Number:         5,
			Title:          "Example PR 5",
			HeadRefName:    "from-branch",
			BaseRefName:    "to-branch",
			Mergeable:      "UNKNOWN",
			ReviewDecision: "CHANGES_REQUESTED",
		},
		{
			Number:         6,
			Title:          "Example PR 6",
			HeadRefName:    "from-branch",
			BaseRefName:    "to-branch",
			Mergeable:      "CONFLICTING",
			ReviewDecision: "CHANGES_REQUESTED",
		},
	}
	//exhaustruct:ignore
	prMap["label_0"] = []gitter.PullRequest{
		{
			Number:      3,
			Title:       "Example PR 3",
			HeadRefName: "from-branch",
			BaseRefName: "to-branch",
		},
	}
	//exhaustruct:ignore
	prMap["repo_0"] = []gitter.PullRequest{
		{
			Number:      4,
			Title:       "Example PR 4",
			HeadRefName: "from-branch",
			BaseRefName: "to-branch",
		},
		{
			Number:      5,
			Title:       "Example PR 5",
			HeadRefName: "from-branch",
			BaseRefName: "to-branch",
		},
	}
	//exhaustruct:ignore
	prMap["author_0"] = []gitter.PullRequest{
		{
			Number:      6,
			Title:       "Example PR 6",
			HeadRefName: "from-branch",
			BaseRefName: "to-branch",
		},
	}

	return prMap, nil
}
