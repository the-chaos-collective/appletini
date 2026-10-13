// Package reviews submits pull request reviews.
package reviews

import (
	"context"
	"errors"
	"fmt"

	"appletini/gitter"
)

var (
	ErrNilClient     = errors.New("graphql client must not be nil")
	ErrEmptyID       = errors.New("pull request id must not be empty")
	ErrEmptyResponse = errors.New("GitHub returned no review")
)

// approveMutation approves a pull request. Values are sent as variables, so a
// comment containing quotes can't break the mutation.
const approveMutation = `mutation Approve($id: ID!, $body: String) {
  addPullRequestReview(input: { pullRequestId: $id, event: APPROVE, body: $body }) {
    pullRequestReview {
      id
      url
    }
  }
}`

// Review is the submitted review.
type Review struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type approveResponse struct {
	AddPullRequestReview struct {
		PullRequestReview *Review `json:"pullRequestReview"`
	} `json:"addPullRequestReview"`
}

// Approver approves pull requests on GitHub.
type Approver struct {
	Client gitter.Mutator
}

// Approve submits an approving review on the pull request with GraphQL node
// id prID, with an optional comment. It returns the created review.
func (a Approver) Approve(ctx context.Context, prID, comment string) (Review, error) {
	if a.Client == nil {
		return Review{}, ErrNilClient
	}

	if prID == "" {
		return Review{}, ErrEmptyID
	}

	vars := map[string]any{"id": prID}
	if comment != "" {
		vars["body"] = comment
	}

	var res approveResponse

	err := a.Client.Mutate(ctx, approveMutation, vars, &res)
	if err != nil {
		return Review{}, fmt.Errorf("approving pull request: %w", err)
	}

	review := res.AddPullRequestReview.PullRequestReview
	if review == nil {
		return Review{}, ErrEmptyResponse
	}

	return *review, nil
}
