package main

import (
	"fmt"
	"git_applet/gitter"
	"git_applet/queries"
	"git_applet/queries/aggregator"
	"git_applet/queries/labeled"
	"git_applet/queries/personal"
	"git_applet/queries/repo"
	"log"
	"os"
	"time"
)

func setupPersonalQuery() (queries.Query, error) {
	personalQuery := personal.PersonalQuery{}

	return personalQuery, nil
}

func setupLabeledQuery() (queries.Query, error) {
	trackers := []labeled.Tracker{}
	for idx, tracker := range Config.Tracking.ByLabel {
		trackers = append(trackers, labeled.Tracker{
			Id:    fmt.Sprintf("labeled_%d", idx),
			Title: tracker.Title,
			Owner: tracker.Owner,
			Repo:  tracker.RepoName,
			Label: tracker.Label,
		})
	}

	return labeled.MakeLabeledQuery(labeled.Config{
		Trackers:       trackers,
		PrAmount:       Config.ItemCount,
		ReviewAmount:   10,
		CommentsAmount: 10,
	})
}

func setupRepoQuery() (queries.Query, error) {
	trackers := []repo.Tracker{}
	for idx, tracker := range Config.Tracking.ByRepo {
		trackers = append(trackers, repo.Tracker{
			Id:    fmt.Sprintf("repo_%d", idx),
			Title: tracker.Title,
			Owner: tracker.Owner,
			Repo:  tracker.RepoName,
		})
	}

	return repo.MakeRepoQuery(repo.Config{
		Trackers:       trackers,
		PrAmount:       Config.ItemCount,
		ReviewAmount:   10,
		CommentsAmount: 10,
	})
}

func setupPolling(mock bool) error {
	personal, err := setupPersonalQuery()
	if err != nil {
		return fmt.Errorf("setting up polling: %w", err)
	}

	labeled, err := setupLabeledQuery()
	if err != nil {
		return fmt.Errorf("setting up polling: %w", err)
	}

	repo, err := setupRepoQuery()
	if err != nil {
		return fmt.Errorf("setting up polling: %w", err)
	}

	queryAggregator = aggregator.QueryAggregator{
		Queries: []queries.Query{
			personal,
			labeled,
			repo,
		},

		Mock: mock,
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
		trackedPrs, err := queryAggregator.GetAll(gqlClient)
		if err != nil {
			log.Printf("when polling for PRs: %v", err)
		}

		hashCheck(trackedPrs, prs)

		time.Sleep(getPollDuration())
	}
}

func getPollDuration() time.Duration {
	return time.Duration(Config.Poll.Frequency * int(time.Second))
}
