package main

import (
	_ "embed"
	"git_applet/gitter"
	"git_applet/ui/pages"
	"log"
)

func main() {
	config, err := loadConfig()
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
