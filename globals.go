package main

import (
	"log"

	"git_applet/gitter"
	"git_applet/queries"
	"git_applet/types"
)

// global and default stuff.
const CONFIG_FILE = "config.json"

var (
	Config          types.Config
	currentHash     string = ""
	gqlClient       gitter.GraphQLClient
	logger          log.Logger
	queryAggregator queries.Query
)
