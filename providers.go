package main

import (
	"os"
	"time"

	"appletini/config"
	"appletini/config/migration"
	"appletini/gitter"
	"appletini/logging"
	"appletini/reviews"
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
			Prefix:          globals.LogPrefix,
		}))
	})
	if err != nil {
		return err
	}

	// * Migrator
	err = deps.Provide(func(feat FeatureFlags, logger logging.Logger) migration.Migrator {
		return migration.Migrator{
			DumpMigrations: feat.DumpMigrations,
			Logger:         logger,
		}
	})
	if err != nil {
		return err
	}

	// * Config
	err = deps.Provide(func(
		globals Globals,
		migrator migration.Migrator,
		logger logging.Logger,
	) (config.Config, error) {
		loader := config.Loader{
			Migrator: migrator,
			Logger:   logger,
		}

		return loader.Load(globals.ConfigPath)
	})
	if err != nil {
		return err
	}

	// * GraphQL Client
	err = deps.Provide(func(conf config.Config) (*gitter.GraphQLClient, error) {
		return gitter.NewGraphQLClient(conf.Github.GraphQL, conf.Computed.GithubToken, gitter.DefaultTimeout)
	})
	if err != nil {
		return err
	}

	// The same client reads (queries) and writes (mutations).
	err = deps.Provide(func(client *gitter.GraphQLClient) gitter.Querier { return client })
	if err != nil {
		return err
	}

	err = deps.Provide(func(client *gitter.GraphQLClient) gitter.Mutator { return client })
	if err != nil {
		return err
	}

	// * Reviews (Approve option)
	err = deps.Provide(func(client gitter.Mutator) reviews.Approver {
		return reviews.Approver{Client: client}
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

	// * Error Channel
	err = deps.Provide(func() HasErrChan {
		return make(HasErrChan)
	})
	if err != nil {
		return err
	}

	// * UI
	err = deps.Provide(func(
		conf config.Config,
		logger logging.Logger,
		prs PRChan,
		hasErr HasErrChan,
		approver reviews.Approver,
	) *pages.IndexPage {
		//exhaustruct:ignore
		return &pages.IndexPage{
			PullRequests: prs,
			HasErr:       hasErr,
			Approver:     approver,
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
