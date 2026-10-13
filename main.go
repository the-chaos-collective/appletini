package main

import (
	"context"
	_ "embed"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/dig"

	"appletini/config"
	"appletini/gitter"
	"appletini/logging"
	"appletini/polling"
	"appletini/ui/pages"
)

type (
	PRChan     chan map[string][]gitter.PullRequest
	HasErrChan chan bool
)

func main() {
	// Cancelled on Ctrl+C / SIGTERM, or when the tray app quits.
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	deps := dig.New()

	err := setupProviders(deps)
	ehp(deps, err)

	err = setupPolling(ctx, deps)
	ehp(deps, err)

	err = render(ctx, deps)
	ehp(deps, err)
}

func render(ctx context.Context, deps *dig.Container) error {
	return deps.Invoke(func(indexPage *pages.IndexPage, logger logging.Logger) error {
		logger.Info("Running")
		return indexPage.Run(ctx)
	})
}

func setupPolling(ctx context.Context, deps *dig.Container) error {
	return deps.Invoke(func(
		flags FeatureFlags,
		logger logging.Logger,
		gqlClient gitter.Querier,
		conf config.Config,
		prs PRChan,
		hasErr HasErrChan,
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

		go func() {
			err := poller.PollPRs(ctx, prs, hasErr)
			if err != nil {
				logger.Error("polling stopped", "err", err)
			}
		}()

		return nil
	})
}

func ehp(deps *dig.Container, err error) {
	if err != nil {
		log_err := deps.Invoke(func(logger logging.Logger) {
			logger.Fatalf("Runtime error: %v\n", err)
		})
		if log_err != nil {
			log.Fatalf("Runtime error: %v\nAdditionally, there was an error retrieving logger: %v\n", err, log_err)
		}
	}
}
