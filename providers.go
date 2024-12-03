package main

import (
	"os"
	"time"

	"appletini/config"
	"appletini/gitter"
	"appletini/logging"
	"appletini/ui/pages"

	"github.com/charmbracelet/log"
	"go.uber.org/dig"
)

func setupProviders(deps *dig.Container) error {
	// * Feature Flags
	err := deps.Provide(featureFlags)
	if err != nil {
		return err
	}

	// * Globals
	err = deps.Provide(globals)
	if err != nil {
		return err
	}

	// * Logger
	err = deps.Provide(func(globals Globals) logging.Logger {
		//exhaustruct:ignore
		return logging.NewCharm(log.NewWithOptions(os.Stderr, log.Options{
			Level:           globals.LogLevel,
			ReportCaller:    false,
			ReportTimestamp: true,
			TimeFormat:      time.DateTime,
			Prefix:          "Appletini",
		}))
	})
	if err != nil {
		return err
	}

	// * Config
	err = deps.Provide(func(
		feat FeatureFlags,
		globals Globals,
		logger logging.Logger,
	) (config.Config, error) {
		loader := config.Loader{
			DumpMigrations: feat.DumpMigrations,
			Logger:         logger,
		}

		return loader.Load(globals.ConfigPath)
	})
	if err != nil {
		return err
	}

	// * GraphQL Client
	err = deps.Provide(func(conf config.Config) *gitter.GraphQLClient {
		return &gitter.GraphQLClient{
			Url:   conf.Github.GraphQL,
			Token: conf.Computed.GithubToken,
		}
	})
	if err != nil {
		return err
	}

	// * PR Channel
	err = deps.Provide(func() PRChan {
		return make(PRChan)
	})
	if err != nil {
		return err
	}

	// * UI
	err = deps.Provide(func(
		conf config.Config,
		logger logging.Logger,
		prs PRChan,
	) pages.IndexPage {
		return pages.IndexPage{
			PullRequests: prs,
			Darkmode:     conf.Darkmode,
			Trackers:     conf.Tracking,
			Logger:       logger,
		}
	})
	if err != nil {
		return err
	}

	return nil
}
