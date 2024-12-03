package main

import (
	_ "embed"
	"log"

	"go.uber.org/dig"

	"git_applet/config"
	"git_applet/gitter"
	"git_applet/polling"
	"git_applet/ui/pages"
)

type PRChan chan map[string][]gitter.PullRequest

func main() {
	deps := dig.New()

	err := setupProviders(deps)
	ehp(err)

	err = setupPolling(deps)
	ehp(err)

	err = render(deps)
	ehp(err)
}

func render(deps *dig.Container) error {
	return deps.Invoke(func(indexPage pages.IndexPage) {
		indexPage.Run()
	})
}

func setupPolling(deps *dig.Container) error {
	return deps.Invoke(func(
		flags FeatureFlags,
		logger *log.Logger,
		gqlClient *gitter.GraphQLClient,
		conf config.Config,
		prs PRChan,
	) error {
		poller := polling.Polling{
			Logger:    logger,
			GqlClient: gqlClient,
			Config:    conf,
		}

		err := poller.Setup(flags.MockQueries)
		if err != nil {
			return err
		}

		go poller.PollPRs(prs)

		return nil
	})
}

func ehp(err error) {
	if err != nil {
		log.Fatalf("Runtime error: %v\n", err)
	}
}
