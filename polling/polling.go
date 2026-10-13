package polling

import (
	"context"
	"errors"
	"fmt"
	"time"

	"appletini/config"
	"appletini/gitter"
	"appletini/hasher"
	"appletini/logging"
	"appletini/queries"
	"appletini/queries/aggregator"
	"appletini/queries/by_author"
	"appletini/queries/by_label"
	"appletini/queries/by_repo"
	"appletini/queries/personal"
	"appletini/trackerkey"
)

const (
	reviewAmount   = 10
	commentsAmount = 10
)

var (
	ErrNotSetup  = errors.New("polling has not been set up")
	ErrNilClient = errors.New("graphql client must not be nil")
)

type Polling struct {
	Logger          logging.Logger
	GqlClient       gitter.Querier
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
		id, err := trackerkey.New(trackerkey.Label, idx)
		if err != nil {
			return nil, err
		}

		trackers = append(trackers, by_label.Tracker{
			Id:    id,
			Title: tracker.Title,
			Owner: tracker.Owner,
			Repo:  tracker.RepoName,
			Label: tracker.Label,
		})
	}

	return by_label.MakeQuery(by_label.Config{
		Trackers:       trackers,
		PrAmount:       p.Config.ItemCount,
		ReviewAmount:   reviewAmount,
		CommentsAmount: commentsAmount,
	})
}

func (p Polling) setupRepoQuery() (queries.Query, error) {
	trackers := []by_repo.Tracker{}
	for idx, tracker := range p.Config.Tracking.ByRepo {
		id, err := trackerkey.New(trackerkey.Repo, idx)
		if err != nil {
			return nil, err
		}

		trackers = append(trackers, by_repo.Tracker{
			Id:    id,
			Title: tracker.Title,
			Owner: tracker.Owner,
			Repo:  tracker.RepoName,
		})
	}

	return by_repo.MakeQuery(by_repo.Config{
		Trackers:       trackers,
		PrAmount:       p.Config.ItemCount,
		ReviewAmount:   reviewAmount,
		CommentsAmount: commentsAmount,
	})
}

func (p Polling) setupAuthorQuery() (queries.Query, error) {
	trackers := []by_author.Tracker{}
	for idx, tracker := range p.Config.Tracking.ByAuthor {
		id, err := trackerkey.New(trackerkey.Author, idx)
		if err != nil {
			return nil, err
		}

		trackers = append(trackers, by_author.Tracker{
			Id:      id,
			Title:   tracker.Title,
			Owner:   tracker.Owner,
			Repo:    tracker.RepoName,
			Authors: tracker.Authors,
		})
	}

	return by_author.MakeQuery(by_author.Config{
		Trackers:       trackers,
		PrAmount:       p.Config.ItemCount,
		ReviewAmount:   reviewAmount,
		CommentsAmount: commentsAmount,
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

// PollPRs polls until ctx is cancelled. Each round reports whether polling
// failed on hasErr and sends the PRs on prs when they changed. It returns nil
// on cancellation, which is the normal way to stop it.
func (p Polling) PollPRs(ctx context.Context, prs chan<- map[string][]gitter.PullRequest, hasErr chan<- bool) error {
	if p.queryAggregator == nil {
		return ErrNotSetup
	}

	if p.GqlClient == nil {
		return ErrNilClient
	}

	prHasher := hasher.Hasher{
		Logger: p.Logger,
	}

	for {
		err := p.poll(ctx, &prHasher, prs, hasErr)
		if ctx.Err() != nil {
			return nil
		}

		if err != nil {
			return err
		}

		timer := time.NewTimer(p.getPollDuration())

		select {
		case <-ctx.Done():
			timer.Stop()
			return nil
		case <-timer.C:
		}
	}
}

// poll runs a single polling round.
func (p Polling) poll(
	ctx context.Context,
	prHasher *hasher.Hasher,
	prs chan<- map[string][]gitter.PullRequest,
	hasErr chan<- bool,
) error {
	trackingPrs, queryErr := p.queryAggregator.GetAll(ctx, p.GqlClient)
	if ctx.Err() != nil {
		return nil
	}

	if queryErr != nil {
		p.Logger.Errorf("When polling for PRs: %v", queryErr)
	}

	select {
	case hasErr <- queryErr != nil:
	case <-ctx.Done():
		return nil
	}

	err := prHasher.Check(ctx, trackingPrs, prs)
	if err != nil && ctx.Err() == nil {
		return fmt.Errorf("publishing PRs: %w", err)
	}

	return nil
}

func (p Polling) getPollDuration() time.Duration {
	return time.Duration(p.Config.Poll.Frequency) * time.Second
}
