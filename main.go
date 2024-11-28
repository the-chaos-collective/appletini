package main

import (
	_ "embed"
	"git_applet/gitter"
	"git_applet/ui/pages"
	"log"
)

func main() {
	err := loadConfig()
	ehp(err)

	gqlClient = gitter.GraphQLClient{
		Url:   Config.Github.GraphQL,
		Token: getCurrentAccessToken(),
	}

	prs := make(chan map[string][]gitter.PullRequest)

	indexPage := pages.IndexPage{
		PullRequests: prs,
		Darkmode:     Config.Darkmode,
		Trackers:     Config.Tracking,
	}

	mockQueries := false
	err = setupPolling(mockQueries)
	ehp(err)

	go pollPRs(prs)

	indexPage.Run()
}

func ehp(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
