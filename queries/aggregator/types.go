package aggregator

import "git_applet/queries"

type QueryAggregator struct {
	Queries []queries.Query
	Mock    bool
}
