package aggregator

import (
	"log"

	"git_applet/queries"
)

type QueryAggregator struct {
	Queries []queries.Query
	Mock    bool
	Logger  *log.Logger
}
