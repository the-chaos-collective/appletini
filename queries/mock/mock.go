package mock

import (
	"git_applet/gitter"
)

func (MockQuery) GetAll(client gitter.GraphQLClient) (map[string][]gitter.PullRequest, error) {
	prMap := map[string][]gitter.PullRequest{}
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
			Mergeable:      "UNKWNOWN",
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
			Mergeable:      "UNKWNOWN",
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
	prMap["foo"] = []gitter.PullRequest{
		{
			Number:      3,
			Title:       "Example PR 3",
			HeadRefName: "from-branch",
			BaseRefName: "to-branch",
		},
	}
	prMap["bar"] = []gitter.PullRequest{
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

	return prMap, nil
}
