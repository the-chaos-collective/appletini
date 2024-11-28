package repo

import (
	"fmt"
	"git_applet/gitter"
	"strings"
)

type Result map[string]struct {
	PullRequests PullRequest `json:"pullRequests"`
}

type PullRequest struct {
	Edges []Edge `json:"edges"`
}

type Edge struct {
	Node PullNode `json:"node"`
}

type PullNode struct {
	Id             string                 `json:"id"`
	Title          string                 `json:"title"`
	Url            string                 `json:"url"`
	BaseRefName    string                 `json:"baseRefName"`
	HeadRefName    string                 `json:"headRefName"`
	ReviewRequests ReviewRequest          `json:"reviewRequests"`
	ReviewDecision string                 `json:"reviewDecision"`
	Permalink      string                 `json:"permalink"`
	Mergeable      string                 `json:"mergeable"`
	Number         int                    `json:"number"`
	State          string                 `json:"state"`
	Review         map[string]interface{} `json:"review"`
}

type ReviewRequest struct {
	TotalCount int `json:"totalCount"`
}

func (result PullRequest) Extract() []gitter.PullRequest {
	all := []gitter.PullRequest{}

	for _, edge := range result.Edges {
		node := edge.Node

		pr := gitter.PullRequest{
			Title:          node.Title,
			BaseRefName:    node.BaseRefName,
			HeadRefName:    node.HeadRefName,
			Number:         node.Number,
			Permalink:      node.Permalink,
			ReviewCount:    0, // TODO
			ReviewRequests: node.ReviewRequests.TotalCount,
			ReviewDecision: node.ReviewDecision,
			Id:             node.Id,
			Mergeable:      node.Mergeable,
		}
		all = append(all, pr)
	}

	return all
}

func (query RepoQuery) GetAll(client gitter.GraphQLClient) (map[string][]gitter.PullRequest, error) {
	res := Result{} // make(map[string]interface{})

	if strings.Trim(query.generatedQuery, "\n") != "" {
		err := gitter.AuthorizedGraphQLQuery[Result](client, query.generatedQuery, &res)
		if err != nil {
			return map[string][]gitter.PullRequest{}, fmt.Errorf("requesting repo PRs: %w", err)
		}
	}

	prs := make(map[string][]gitter.PullRequest)

	for key, entry := range res {
		prs[key] = entry.PullRequests.Extract()
	}

	return prs, nil
}
