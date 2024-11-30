package main

import (
	"log"

	"git_applet/gitter"
	"git_applet/queries"
)

// global and default stuff.
const CONFIG_FILE = "config.json"

var (
	gqlClient       gitter.GraphQLClient
	logger          log.Logger
	queryAggregator queries.Query
)
