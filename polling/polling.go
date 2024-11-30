package polling

import (
	"fmt"
	"log"
	"time"

	"git_applet/config"
	"git_applet/gitter"
	"git_applet/hasher"
	"git_applet/queries"
	"git_applet/queries/aggregator"
	"git_applet/queries/by_author"
	"git_applet/queries/by_label"
	"git_applet/queries/by_repo"
	"git_applet/queries/personal"
)

type Polling struct {
	Logger          *log.Logger
	GqlClient       *gitter.GraphQLClient
	Config          config.Config
	queryAggregator queries.Query
}

func setupPersonalQuery() (queries.Query, error) {
	personalQuery := personal.Query{}

	return personalQuery, nil
}

func (p Polling) setupLabelQuery() (queries.Query, error) {
	trackers := []by_label.Tracker{}
	for idx, tracker := range p.Config.Tracking.ByLabel {
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
		PrAmount:       p.Config.ItemCount,
		ReviewAmount:   10,
		CommentsAmount: 10,
	})
}

func (p Polling) setupRepoQuery() (queries.Query, error) {
	trackers := []by_repo.Tracker{}
	for idx, tracker := range p.Config.Tracking.ByRepo {
		trackers = append(trackers, by_repo.Tracker{
			Id:    fmt.Sprintf("repo_%d", idx),
			Title: tracker.Title,
			Owner: tracker.Owner,
			Repo:  tracker.RepoName,
		})
	}

	return by_repo.MakeQuery(by_repo.Config{
		Trackers:       trackers,
		PrAmount:       p.Config.ItemCount,
		ReviewAmount:   10,
		CommentsAmount: 10,
	})
}

func (p Polling) setupAuthorQuery() (queries.Query, error) {
	trackers := []by_author.Tracker{}
	for idx, tracker := range p.Config.Tracking.ByAuthor {
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
		PrAmount:       p.Config.ItemCount,
		ReviewAmount:   10,
		CommentsAmount: 10,
	})
}

func (p *Polling) Setup(mock bool) error {
	labeled, err := p.setupLabelQuery()
	if err != nil {
		return fmt.Errorf("setting up label polling: %w", err)
	}

	repo, err := p.setupRepoQuery()
	if err != nil {
		return fmt.Errorf("setting up repo polling: %w", err)
	}

	author, err := p.setupAuthorQuery()
	if err != nil {
		return fmt.Errorf("setting up author polling: %w", err)
	}

	queries := []queries.Query{
		labeled,
		repo,
		author,
	}

	if p.Config.Tracking.Personal {
		personal, err := setupPersonalQuery()
		if err != nil {
			return fmt.Errorf("setting up personal polling: %w", err)
		}
		queries = append(queries, personal)
	}

	p.queryAggregator = aggregator.QueryAggregator{
		Queries: queries,
		Mock:    mock,
		Logger:  p.Logger,
	}

	return nil
}

func (p Polling) PollPRs(prs chan<- map[string][]gitter.PullRequest) {
	hasher := hasher.Hasher{
		Logger: p.Logger,
	}
	for {
		trackingPrs, err := p.queryAggregator.GetAll(*p.GqlClient)
		if err != nil {
			p.Logger.Printf("when polling for PRs: %v", err)
		}

		hasher.Check(trackingPrs, prs)

		time.Sleep(p.getPollDuration())
	}
}

func (p Polling) getPollDuration() time.Duration {
	return time.Duration(p.Config.Poll.Frequency * int(time.Second))
}
