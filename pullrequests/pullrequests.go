// Package pullrequests merges and closes pull requests.
package pullrequests

import (
	"context"
	"errors"
	"fmt"

	"appletini/gitter"
)

var (
	ErrNilClient     = errors.New("graphql client must not be nil")
	ErrEmptyID       = errors.New("pull request id must not be empty")
	ErrEmptyResponse = errors.New("GitHub returned no pull request")
	ErrNotMerged     = errors.New("GitHub did not merge the pull request")
	ErrNotClosed     = errors.New("GitHub did not close the pull request")
)

// mergeMutation merges a pull request with the repository's default merge
// method. The id is sent as a variable.
const mergeMutation = `mutation Merge($id: ID!) {
  mergePullRequest(input: { pullRequestId: $id }) {
    pullRequest {
      number
      merged
      url
    }
  }
}`

// closeMutation closes a pull request without merging it.
const closeMutation = `mutation Close($id: ID!) {
  closePullRequest(input: { pullRequestId: $id }) {
    pullRequest {
      number
      closed
      url
    }
  }
}`

// PullRequest is the state GitHub reports back after a merge or close.
type PullRequest struct {
	Number int    `json:"number"`
	Merged bool   `json:"merged"`
	Closed bool   `json:"closed"`
	URL    string `json:"url"`
}

type mergeResponse struct {
	MergePullRequest struct {
		PullRequest *PullRequest `json:"pullRequest"`
	} `json:"mergePullRequest"`
}

type closeResponse struct {
	ClosePullRequest struct {
		PullRequest *PullRequest `json:"pullRequest"`
	} `json:"closePullRequest"`
}

// Actions merges and closes pull requests on GitHub.
type Actions struct {
	Client gitter.Mutator
}

func (a Actions) check(prID string) error {
	if a.Client == nil {
		return ErrNilClient
	}

	if prID == "" {
		return ErrEmptyID
	}

	return nil
}

// Merge merges the pull request with GraphQL node id prID.
func (a Actions) Merge(ctx context.Context, prID string) (PullRequest, error) {
	err := a.check(prID)
	if err != nil {
		return PullRequest{}, err
	}

	var res mergeResponse

	err = a.Client.Mutate(ctx, mergeMutation, map[string]any{"id": prID}, &res)
	if err != nil {
		return PullRequest{}, fmt.Errorf("merging pull request: %w", err)
	}

	pr := res.MergePullRequest.PullRequest
	if pr == nil {
		return PullRequest{}, ErrEmptyResponse
	}

	if !pr.Merged {
		return *pr, ErrNotMerged
	}

	return *pr, nil
}

// Close closes the pull request with GraphQL node id prID without merging it.
func (a Actions) Close(ctx context.Context, prID string) (PullRequest, error) {
	err := a.check(prID)
	if err != nil {
		return PullRequest{}, err
	}

	var res closeResponse

	err = a.Client.Mutate(ctx, closeMutation, map[string]any{"id": prID}, &res)
	if err != nil {
		return PullRequest{}, fmt.Errorf("closing pull request: %w", err)
	}

	pr := res.ClosePullRequest.PullRequest
	if pr == nil {
		return PullRequest{}, ErrEmptyResponse
	}

	if !pr.Closed {
		return *pr, ErrNotClosed
	}

	return *pr, nil
}
