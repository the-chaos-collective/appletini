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
	conf, err := config.Load(CONFIG_FILE, DUMP_MIGRATIONS, logger)
	ehp(err, logger)

	// * GraphQL Client
	gqlClient := &gitter.GraphQLClient{
		Url:   conf.Github.GraphQL,
		Token: conf.Computed.GithubToken,
	}

	// * Polling
	poller := polling.Polling{
		Logger:    logger,
		GqlClient: gqlClient,
		Config:    conf,
	}

	prs := make(chan map[string][]gitter.PullRequest)

	err = poller.Setup(MOCK_QUERIES)
	ehp(err, logger)

	go poller.PollPRs(prs)

	// * UI
	indexPage := pages.IndexPage{
		PullRequests: prs,
		Darkmode:     conf.Darkmode,
		Trackers:     conf.Tracking,
		Logger:       logger,
	}

	indexPage.Run()
}

func ehp(err error, logger *log.Logger) {
	if err != nil {
		logger.Fatalf("Runtime error: %v\n", err)
	}
}
