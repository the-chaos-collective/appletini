package main

import (
	_ "embed"
	"log"

	"go.uber.org/dig"

	"appletini/config"
	"appletini/gitter"
	"appletini/logging"
	"appletini/polling"
	"appletini/ui/pages"
)

type PRChan chan map[string][]gitter.PullRequest

func main() {
	deps := dig.New()

	err := setupProviders(deps)
	ehp(deps, err)

	err = setupPolling(deps)
	ehp(deps, err)

	err = render(deps)
	ehp(deps, err)
}

func render(deps *dig.Container) error {
	return deps.Invoke(func(indexPage pages.IndexPage, logger logging.Logger) {
		logger.Info("Running")
		indexPage.Run()
	})
}

func setupPolling(deps *dig.Container) error {
	return deps.Invoke(func(
		flags FeatureFlags,
		logger logging.Logger,
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

func ehp(deps *dig.Container, err error) {
	if err != nil {
		err := deps.Invoke(func(logger logging.Logger) {
			logger.Fatalf("Runtime error: %v\n", err)
		})
		if err != nil {
			log.Fatalf("Runtime error: %v\n", err)
		}
	}
}
