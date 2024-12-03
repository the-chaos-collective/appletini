package aggregator

import (
	"appletini/logging"
	"appletini/queries"
)

type QueryAggregator struct {
	Queries []queries.Query
	Mock    bool
	Logger  logging.Logger
}
