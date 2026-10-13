package queries

import (
	"context"

	"appletini/gitter"
)

// Query fetches the pull requests for one kind of tracker, keyed by tracker id.
// Implementations must honour ctx cancellation.
type Query interface {
	GetAll(ctx context.Context, client gitter.Querier) (map[string][]gitter.PullRequest, error)
}
