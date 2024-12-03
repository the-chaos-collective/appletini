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

func main() {
	deps := dig.New()

	// * Feature Flags
	err := deps.Provide(LoadFeatureFlags)
	ehp(err)

	// * Globals
	err = deps.Provide(LoadGlobals)
	ehp(err)

	// * Logger
	err = deps.Provide(log.Default)
	ehp(err)

	// * Config
	err = deps.Provide(func(feat FeatureFlags, globals Globals, logger *log.Logger) (config.Config, error) {
		loader := config.Loader{
			DumpMigrations: feat.DumpMigrations,
			Logger:         logger,
		}

		return loader.Load(globals.ConfigPath)
	})
	ehp(err)

	// * GraphQL Client
	err = deps.Provide(func(conf config.Config) *gitter.GraphQLClient {
		return &gitter.GraphQLClient{
			Url:   conf.Github.GraphQL,
			Token: conf.Computed.GithubToken,
		}
	})
	ehp(err)

	// * Polling
	var poller polling.Polling

	err = deps.Invoke(func(flags FeatureFlags, logger *log.Logger, gqlClient *gitter.GraphQLClient, conf config.Config) error {
		poller = polling.Polling{
			Logger:    logger,
			GqlClient: gqlClient,
			Config:    conf,
		}

		return poller.Setup(flags.MockQueries)
	})
	ehp(err)

	prs := make(chan map[string][]gitter.PullRequest)

	// * UI
	var indexPage pages.IndexPage

	err = deps.Invoke(func(conf config.Config, logger *log.Logger) {
		indexPage = pages.IndexPage{
			PullRequests: prs,
			Darkmode:     conf.Darkmode,
			Trackers:     conf.Tracking,
			Logger:       logger,
		}
	})
	ehp(err)

	go poller.PollPRs(prs)
	indexPage.Run()
}

func ehp(err error) {
	if err != nil {
		log.Fatalf("Runtime error: %v\n", err)
	}
}
