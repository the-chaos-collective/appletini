package main

import (
	_ "embed"
	"log"

	"git_applet/config"
	"git_applet/gitter"
	"git_applet/polling"
	"git_applet/ui/pages"
)

func main() {
	// * Logger
	logger := log.Default()

	// * Config
	config, err := config.Load(CONFIG_FILE, DUMP_MIGRATIONS)
	ehp(err, logger)

	// * GraphQL Client
	gqlClient := &gitter.GraphQLClient{
		Url:   config.Github.GraphQL,
		Token: config.Computed.GithubToken,
	}

	// * Polling
	poller := polling.Polling{
		Logger:    logger,
		GqlClient: gqlClient,
		Config:    config,
	}

	prs := make(chan map[string][]gitter.PullRequest)

	mockQueries := false
	err = poller.Setup(mockQueries)
	ehp(err, logger)

	go poller.PollPRs(prs)

	// * UI
	indexPage := pages.IndexPage{
		PullRequests: prs,
		Darkmode:     config.Darkmode,
		Trackers:     config.Tracking,
		Logger:       logger,
	}

	indexPage.Run()
}

func ehp(err error, logger *log.Logger) {
	if err != nil {
		logger.Fatal(err)
	}
}
