package personal

import (
	"context"
	"errors"
	"fmt"

	"appletini/gitter"
)

var ViewerQuery = `query fetchPRs {
			viewer {
				pullRequests(
				orderBy: {field: CREATED_AT, direction: ASC}
				first: 100
				states: [OPEN]
				) {
				edges {
					node {
					id
					repository {
						branchProtectionRules(first: 100) {
						edges {
							node {
							allowsDeletions
							allowsForcePushes
							requiresApprovingReviews
							}
						}
						}
						name
						url
						owner {
						login
						}
					}
					reviewDecision
					title
					baseRefName
					headRefName
					number
					permalink
					reviewRequests {
						totalCount
					}
					reviews(first: 12) {
						totalCount
						nodes {
						state
						body
						url
						}
					}
					mergeable
					viewerLatestReview {
						state
					}
					commits(last: 1) {
						nodes {
						commit {
							statusCheckRollup {
							state
							}
						}
						}
					}
					}
				}
				}
			}
			}`

// Key is the result key used for the viewer's own pull requests.
const Key = "personal"

var ErrNilClient = errors.New("graphql client must not be nil")

func (Query) GetAll(ctx context.Context, client gitter.Querier) (map[string][]gitter.PullRequest, error) {
	if client == nil {
		return map[string][]gitter.PullRequest{}, ErrNilClient
	}

	var res Response

	err := client.Query(ctx, ViewerQuery, &res)
	if err != nil {
		return map[string][]gitter.PullRequest{}, fmt.Errorf("requesting personal PRs: %w", err)
	}

	prMap := map[string][]gitter.PullRequest{}
	prMap[Key] = res.Extract()

	return prMap, nil
}
