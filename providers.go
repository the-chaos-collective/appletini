package main

import (
	"log"

	"git_applet/config"
	"git_applet/gitter"
	"git_applet/ui/pages"

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
	err = deps.Provide(log.Default)
	if err != nil {
		return err
	}

	// * Config
	err = deps.Provide(func(
		feat FeatureFlags,
		globals Globals,
		logger *log.Logger,
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
		logger *log.Logger,
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
