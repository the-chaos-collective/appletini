package by_author

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"appletini/gitter"
)

var (
	ErrNilClient    = errors.New("graphql client must not be nil")
	ErrMalformedKey = errors.New("malformed author result key (expected [trackerId]_[authorIndex])")
)

type Result map[string]NodeMap

type NodeMap struct {
	Nodes []PullRequest `json:"nodes"`
}

type PullRequest struct {
	Id             string     `json:"id"`
	Title          string     `json:"title"`
	Url            string     `json:"url"`
	BaseRefName    string     `json:"baseRefName"`
	HeadRefName    string     `json:"headRefName"`
	ReviewRequests TotalCount `json:"reviewRequests"`
	ReviewDecision string     `json:"reviewDecision"`
	Permalink      string     `json:"permalink"`
	Mergeable      string     `json:"mergeable"`
	Number         int        `json:"number"`
	State          string     `json:"state"`
	Reviews        TotalCount `json:"reviews"`
}

type TotalCount struct {
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
			ReviewCount:    node.Reviews.TotalCount,
			ReviewRequests: node.ReviewRequests.TotalCount,
			ReviewDecision: node.ReviewDecision,
			Id:             node.Id,
			Mergeable:      node.Mergeable,
		}
		all = append(all, pr)
	}

	return all
}

// resultKey is one "[trackerId]_[authorIndex]" alias from the query response.
type resultKey struct {
	raw         string
	trackerID   string
	authorIndex int
}

func parseResultKey(raw string) (resultKey, error) {
	sep := strings.LastIndex(raw, "_")
	if sep <= 0 || sep == len(raw)-1 {
		return resultKey{}, fmt.Errorf("%w: %q", ErrMalformedKey, raw)
	}

	idx, err := strconv.Atoi(raw[sep+1:])
	if err != nil {
		return resultKey{}, fmt.Errorf("%w: %q: %w", ErrMalformedKey, raw, err)
	}

	return resultKey{raw: raw, trackerID: raw[:sep], authorIndex: idx}, nil
}

// sortedKeys orders keys by tracker then author index so the merged PR list
// is stable between polls (map iteration order is random, which otherwise
// changes the hash and rebuilds the menu on every poll).
func sortedKeys(res Result) ([]resultKey, error) {
	keys := make([]resultKey, 0, len(res))

	for raw := range res {
		key, err := parseResultKey(raw)
		if err != nil {
			return nil, err
		}

		keys = append(keys, key)
	}

	slices.SortFunc(keys, func(a, b resultKey) int {
		if c := strings.Compare(a.trackerID, b.trackerID); c != 0 {
			return c
		}

		return a.authorIndex - b.authorIndex
	})

	return keys, nil
}

func (query Query) GetAll(ctx context.Context, client gitter.Querier) (map[string][]gitter.PullRequest, error) {
	prs := make(map[string][]gitter.PullRequest)

	if !query.shouldBeExecuted {
		return prs, nil
	}

	if client == nil {
		return map[string][]gitter.PullRequest{}, ErrNilClient
	}

	res := Result{}

	err := client.Query(ctx, query.generatedQuery, &res)
	if err != nil {
		return map[string][]gitter.PullRequest{}, fmt.Errorf("requesting author PRs: %w", err)
	}

	keys, err := sortedKeys(res)
	if err != nil {
		return map[string][]gitter.PullRequest{}, fmt.Errorf("reading author PRs: %w", err)
	}

	for _, key := range keys {
		prs[key.trackerID] = append(prs[key.trackerID], res[key.raw].Extract()...)
	}

	return prs, nil
}
