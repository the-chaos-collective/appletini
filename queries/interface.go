package queries

import "git_applet/gitter"

type Query interface {
	GetAll(client gitter.GraphQLClient) (map[string][]gitter.PullRequest, error)
}
