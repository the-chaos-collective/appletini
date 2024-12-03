package queries

import "appletini/gitter"

type Query interface {
	GetAll(client gitter.GraphQLClient) (map[string][]gitter.PullRequest, error)
}
