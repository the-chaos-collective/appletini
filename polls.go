package main

import (
	"fmt"
	"log"
	"os"
	"time"

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

func setupLabelQuery() (queries.Query, error) {
	trackers := []by_label.Tracker{}
	for idx, tracker := range Config.Tracking.ByLabel {
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
		PrAmount:       Config.ItemCount,
		ReviewAmount:   10,
		CommentsAmount: 10,
	})
}

func setupRepoQuery() (queries.Query, error) {
	trackers := []by_repo.Tracker{}
	for idx, tracker := range Config.Tracking.ByRepo {
		trackers = append(trackers, by_repo.Tracker{
			Id:    fmt.Sprintf("repo_%d", idx),
			Title: tracker.Title,
			Owner: tracker.Owner,
			Repo:  tracker.RepoName,
		})
	}

	return by_repo.MakeQuery(by_repo.Config{
		Trackers:       trackers,
		PrAmount:       Config.ItemCount,
		ReviewAmount:   10,
		CommentsAmount: 10,
	})
}

func setupAuthorQuery() (queries.Query, error) {
	trackers := []by_author.Tracker{}
	for idx, tracker := range Config.Tracking.ByAuthor {
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
		PrAmount:       Config.ItemCount,
		ReviewAmount:   10,
		CommentsAmount: 10,
	})
}

func setupPolling(mock bool) error {
	labeled, err := setupLabelQuery()
	if err != nil {
		return fmt.Errorf("setting up polling: %w", err)
	}

	repo, err := setupRepoQuery()
	if err != nil {
		return fmt.Errorf("setting up polling: %w", err)
	}

	author, err := setupAuthorQuery()
	if err != nil {
		return fmt.Errorf("setting up  author polling: %w", err)
	}

	queries := []queries.Query{
		labeled,
		repo,
		author,
	}

	if Config.Tracking.Personal {
		personal, err := setupPersonalQuery()
		if err != nil {
			return fmt.Errorf("setting up polling: %w", err)
		}
		queries = append(queries, personal)
	}

	queryAggregator = aggregator.QueryAggregator{
		Queries: queries,
		Mock:    mock,
	}

	return nil
}

func getCurrentAccessToken() string {
	token, present := os.LookupEnv(Config.Github.Token)
	if !present {
		log.Fatal("token not present")
	}
	return token
}

func pollPRs(prs chan<- map[string][]gitter.PullRequest) {
	for {
		trackingPrs, err := queryAggregator.GetAll(gqlClient)
		if err != nil {
			log.Printf("when polling for PRs: %v", err)
		}

		hashCheck(trackingPrs, prs)

		time.Sleep(getPollDuration())
	}
}

func getPollDuration() time.Duration {
	return time.Duration(Config.Poll.Frequency * int(time.Second))
}
