package personal

import (
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
					}
				}
				}
			}
			}`

func (Query) GetAll(client gitter.GraphQLClient) (map[string][]gitter.PullRequest, error) {
	var res Response

	err := gitter.AuthorizedGraphQLQuery[Response](client, ViewerQuery, &res)
	if err != nil {
		return map[string][]gitter.PullRequest{}, fmt.Errorf("requesting personal PRs: %w", err)
	}

	prMap := map[string][]gitter.PullRequest{}
	prMap["personal"] = res.Extract()

	return prMap, nil
}
