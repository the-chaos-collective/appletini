package main

import (
	"git_applet/gitter"
	"git_applet/queries"
	"git_applet/types"
)

// global and default stuff.
const CONFIG_FILE = "config.json"

var (
	Config      types.Config
	currentHash string = ""
	trackedPrs  map[string][]gitter.PullRequest
	gqlClient   gitter.GraphQLClient

	queryAggregator queries.Query
)
