package main

import (
	_ "embed"
	"log"

	"git_applet/config"
	"git_applet/gitter"
	"git_applet/ui/pages"
)

func main() {
	config, err := config.Load(CONFIG_FILE)
	ehp(err)

	logger = *log.Default()
	gqlClient = gitter.GraphQLClient{
		Url:   config.Github.GraphQL,
		Token: getCurrentAccessToken(config),
	}

	prs := make(chan map[string][]gitter.PullRequest)

	indexPage := pages.IndexPage{
		PullRequests: prs,
		Darkmode:     config.Darkmode,
		Trackers:     config.Tracking,
	}

	mockQueries := false
	err = setupPolling(config, mockQueries)
	ehp(err)

	go pollPRs(config, prs)

	indexPage.Run()
}

func ehp(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
