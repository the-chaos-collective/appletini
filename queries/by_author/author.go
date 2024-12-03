package by_author

import (
	"fmt"
	"strings"

	"appletini/gitter"
)

type Result map[string]NodeMap

type NodeMap struct {
	Nodes []PullRequest `json:"nodes"`
}

type PullRequest struct {
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

func (result NodeMap) Extract() []gitter.PullRequest {
	all := []gitter.PullRequest{}

	for _, node := range result.Nodes {

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

func (query Query) GetAll(client gitter.GraphQLClient) (map[string][]gitter.PullRequest, error) {
	prs := make(map[string][]gitter.PullRequest)

	if !query.shouldBeExecuted {
		return prs, nil
	}

	res := Result{}

	err := gitter.AuthorizedGraphQLQuery[Result](client, query.generatedQuery, &res)
	if err != nil {
		return map[string][]gitter.PullRequest{}, fmt.Errorf("author response failed: %w", err)
	}

	for key, prMap := range res {
		parts := strings.Split(key, "_")
		baseKey := strings.Join(parts[0:2], "_")
		_, exists := prs[baseKey]
		if !exists {
			prs[baseKey] = prMap.Extract()
		} else {
			prs[baseKey] = append(prs[baseKey], prMap.Extract()...)
		}
	}

	return prs, nil
}
