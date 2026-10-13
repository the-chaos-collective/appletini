package aggregator

import (
	"context"
	"fmt"
	"maps"

	"appletini/gitter"
	"appletini/queries/mock"
)

// GetAll runs every query in order and merges the results. It stops at the
// first failing query and returns what was collected so far together with the
// error (partial results are intentional).
func (qa QueryAggregator) GetAll(ctx context.Context, client gitter.Querier) (map[string][]gitter.PullRequest, error) {
	prMap := make(map[string][]gitter.PullRequest)

	// mock override
	if qa.Mock {
		return mock.MockQuery{}.GetAll(ctx, client)
	}

	for _, query := range qa.Queries {
		err := ctx.Err()
		if err != nil {
			return prMap, fmt.Errorf("aggregating queries: %w", err)
		}

		prtmp, err := query.GetAll(ctx, client)
		if err != nil {
			return prMap, err
		}

		maps.Copy(prMap, prtmp)
	}

	return prMap, nil
}
