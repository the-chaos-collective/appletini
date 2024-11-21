package main

import (
	"net/http"

	"git_applet/gitter"
	"git_applet/queries"
	"git_applet/types"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/systray"
)

// global and default stuff.
const CONFIG_FILE = "config.json"

var (
	Config             types.Config
	currentHash        string = ""
	currentTrackedHash string = ""
	trackedPrs         map[string][]gitter.PullRequest
	gqlClient          gitter.GraphQLClient

	queryAggregator queries.Query
)

// visual stuff.
var (
	trackedPRs map[string]*fyne.MenuItem

	mprincipal *fyne.Menu
	status     *systray.MenuItem
	desk       desktop.App
)

// global singleton.
var client *http.Client
