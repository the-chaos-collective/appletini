package gitter

import (
	"context"
	"fmt"

	"github.com/machinebox/graphql"
)

func ApprovePullRequest(url string, token string, ctx context.Context, id string, body string) error { // TODO: fix the graphql injection xD

	req := graphql.NewRequest(fmt.Sprintf(`mutation {
		addPullRequestReview(input: {
		  pullRequestId: "%s",
		  event: APPROVE,
		  body: "%s"
		}) {
		  pullRequestReview {
			id
			url
		  }
		}
	  }`, id, body))

	client := graphql.NewClient(url)
	// TODO: do the same for organization

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", token))

	err := client.Run(context.Background(), req, nil)
	if err != nil {
		return fmt.Errorf("error approving PR: %w", err)
	}

	return nil
}

type PullRequest struct {
	Title       string `yaml:"title"`
	BaseRefName string `yaml:"baseRefName"`
	HeadRefName string `yaml:"headRefName"`
	Number      int    `yaml:"number"`
	Permalink   string `yaml:"permalink"`
	ReviewCount int    `yaml:"reviewCount"`

	ReviewRequests int    `yaml:"reviewRequests"`
	ReviewDecision string `yaml:"reviewDecision"`
	Id             string `yaml:"id"`

	Mergeable  string         `yaml:"mergeable"`
	Remainder_ map[string]any `yaml:",inline"`
}
