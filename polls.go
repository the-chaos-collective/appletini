package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"git_applet/config"
	"git_applet/gitter"
	"git_applet/queries"
	"git_applet/queries/aggregator"
	"git_applet/queries/by_author"
	"git_applet/queries/by_label"
	"git_applet/queries/by_repo"
	"git_applet/queries/personal"
)

func setupPersonalQuery() (queries.Query, error) {
	personalQuery := personal.Query{}

	return personalQuery, nil
}

func setupLabelQuery(config config.Config) (queries.Query, error) {
	trackers := []by_label.Tracker{}
	for idx, tracker := range config.Tracking.ByLabel {
		trackers = append(trackers, by_label.Tracker{
			Id:    fmt.Sprintf("label_%d", idx),
			Title: tracker.Title,
			Owner: tracker.Owner,
			Repo:  tracker.RepoName,
			Label: tracker.Label,
		})
	}

	return by_label.MakeQuery(by_label.Config{
		Trackers:       trackers,
		PrAmount:       config.ItemCount,
		ReviewAmount:   10,
		CommentsAmount: 10,
	})
}

func setupRepoQuery(config config.Config) (queries.Query, error) {
	trackers := []by_repo.Tracker{}
	for idx, tracker := range config.Tracking.ByRepo {
		trackers = append(trackers, by_repo.Tracker{
			Id:    fmt.Sprintf("repo_%d", idx),
			Title: tracker.Title,
			Owner: tracker.Owner,
			Repo:  tracker.RepoName,
		})
	}

	return by_repo.MakeQuery(by_repo.Config{
		Trackers:       trackers,
		PrAmount:       config.ItemCount,
		ReviewAmount:   10,
		CommentsAmount: 10,
	})
}

func setupAuthorQuery(config config.Config) (queries.Query, error) {
	trackers := []by_author.Tracker{}
	for idx, tracker := range config.Tracking.ByAuthor {
		trackers = append(trackers, by_author.Tracker{
			Id:      fmt.Sprintf("author_%d", idx),
			Title:   tracker.Title,
			Owner:   tracker.Owner,
			Repo:    tracker.RepoName,
			Authors: tracker.Authors,
		})
	}

	return by_author.MakeQuery(by_author.Config{
		Trackers:       trackers,
		PrAmount:       config.ItemCount,
		ReviewAmount:   10,
		CommentsAmount: 10,
	})
}

func setupPolling(config config.Config, mock bool) error {
	labeled, err := setupLabelQuery(config)
	if err != nil {
		return fmt.Errorf("setting up label polling: %w", err)
	}

	repo, err := setupRepoQuery(config)
	if err != nil {
		return fmt.Errorf("setting up repo polling: %w", err)
	}

	author, err := setupAuthorQuery(config)
	if err != nil {
		return fmt.Errorf("setting up author polling: %w", err)
	}

	queries := []queries.Query{
		labeled,
		repo,
		author,
	}

	if config.Tracking.Personal {
		personal, err := setupPersonalQuery()
		if err != nil {
			return fmt.Errorf("setting up personal polling: %w", err)
		}
		queries = append(queries, personal)
	}

	queryAggregator = aggregator.QueryAggregator{
		Queries: queries,
		Mock:    mock,
	}

	return nil
}

func getCurrentAccessToken(config config.Config) string {
	token, present := os.LookupEnv(config.Github.Token)
	if !present {
		log.Fatal("token not present")
	}
	return token
}

func pollPRs(config config.Config, prs chan<- map[string][]gitter.PullRequest) {
	currentHash := ""
	for {
		trackingPrs, err := queryAggregator.GetAll(gqlClient)
		if err != nil {
			log.Printf("when polling for PRs: %v", err)
		}

		hashCheck(&currentHash, trackingPrs, prs)

		time.Sleep(getPollDuration(config))
	}
}

func getPollDuration(config config.Config) time.Duration {
	return time.Duration(config.Poll.Frequency * int(time.Second))
}
