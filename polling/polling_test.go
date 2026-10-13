package polling

import (
	"context"
	"errors"
	"io"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/charmbracelet/log"

	"appletini/config"
	v2 "appletini/config/migration/v2"
	"appletini/gitter"
	"appletini/gitter/gittertest"
	"appletini/hasher"
	"appletini/logging"
	"appletini/queries/aggregator"
)

func discardLogger() logging.Logger {
	return logging.NewCharm(log.New(io.Discard))
}

func validConfig() config.Config {
	//exhaustruct:ignore
	conf := config.Config{ItemCount: 10}
	conf.Poll.Frequency = 0
	conf.Tracking.ByLabel = v2.LabeledSet{{Title: "Bugs", Owner: "o", RepoName: "r", Label: "bug"}}

	return conf
}

func TestSetup(t *testing.T) {
	testCases := []struct {
		name          string
		mutate        func(*config.Config)
		expectQueries int
		errContains   string
	}{
		{name: "valid without personal", mutate: func(*config.Config) {}, expectQueries: 3},
		{name: "valid with personal", mutate: func(c *config.Config) { c.Tracking.Personal = true }, expectQueries: 4},
		{name: "invalid label tracker", mutate: func(c *config.Config) { c.Tracking.ByLabel[0].Label = "" }, errContains: "setting up label polling"},
		{name: "invalid item count", mutate: func(c *config.Config) { c.ItemCount = 0 }, errContains: "PrAmount must not be zero or less"},
		{name: "invalid repo tracker", mutate: func(c *config.Config) {
			c.Tracking.ByRepo = v2.RepoSet{{Title: "", Owner: "o", RepoName: "r"}}
		}, errContains: "setting up repo polling"},
		{name: "invalid author tracker", mutate: func(c *config.Config) {
			//exhaustruct:ignore
			c.Tracking.ByAuthor = v2.AuthorSet{{Title: "t", Owner: "o", RepoName: "r"}}
		}, errContains: "setting up author polling"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			conf := validConfig()
			tc.mutate(&conf)

			//exhaustruct:ignore
			p := Polling{Logger: discardLogger(), GqlClient: &gittertest.Fake{}, Config: conf}

			err := p.Setup(false)

			if tc.errContains != "" {
				if err == nil || !strings.Contains(err.Error(), tc.errContains) {
					t.Fatalf("got error %v, expected it to contain %q", err, tc.errContains)
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			qa, ok := p.queryAggregator.(aggregator.QueryAggregator)
			if !ok {
				t.Fatalf("aggregator has type %T", p.queryAggregator)
			}

			if len(qa.Queries) != tc.expectQueries {
				t.Fatalf("got %d queries, expected %d", len(qa.Queries), tc.expectQueries)
			}
		})
	}
}

// scriptedQuery returns results[i] / errs[i] on the i-th call, then repeats the last.
type scriptedQuery struct {
	mu      sync.Mutex
	calls   int
	results []map[string][]gitter.PullRequest
	errs    []error
}

func (s *scriptedQuery) GetAll(ctx context.Context, _ gitter.Querier) (map[string][]gitter.PullRequest, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	i := min(s.calls, len(s.results)-1)
	s.calls++

	if ctx.Err() != nil {
		return map[string][]gitter.PullRequest{}, ctx.Err()
	}

	return s.results[i], s.errs[i]
}

type event struct {
	hasErr *bool
	prs    map[string][]gitter.PullRequest
}

func TestPollPRs(t *testing.T) {
	errBoom := errors.New("boom")

	//exhaustruct:ignore
	prsA := map[string][]gitter.PullRequest{"repo_0": {{Number: 1}}}
	//exhaustruct:ignore
	prsB := map[string][]gitter.PullRequest{"repo_0": {{Number: 2}}}

	b := func(v bool) *bool { return &v }

	testCases := []struct {
		name     string
		results  []map[string][]gitter.PullRequest
		errs     []error
		expected []event
	}{
		{
			name:    "reports success then sends PRs",
			results: []map[string][]gitter.PullRequest{prsA},
			errs:    []error{nil},
			expected: []event{
				{hasErr: b(false)},
				{prs: prsA},
			},
		},
		{
			name:    "reports error, recovers, only sends changed PRs",
			results: []map[string][]gitter.PullRequest{prsA, {}, prsA, prsB},
			errs:    []error{nil, errBoom, nil, nil},
			expected: []event{
				{hasErr: b(false)},
				{prs: prsA},
				{hasErr: b(true)},
				{prs: map[string][]gitter.PullRequest{}}, // partial (empty) result on error is sent: intended
				{hasErr: b(false)},
				{prs: prsA},
				{hasErr: b(false)},
				{prs: prsB},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			//exhaustruct:ignore
			p := Polling{
				Logger:          discardLogger(),
				GqlClient:       &gittertest.Fake{},
				Config:          validConfig(), // frequency 0: poll again immediately
				queryAggregator: &scriptedQuery{results: tc.results, errs: tc.errs},
			}

			prs := make(chan map[string][]gitter.PullRequest)
			hasErr := make(chan bool)
			done := make(chan error, 1)

			go func() { done <- p.PollPRs(ctx, prs, hasErr) }()

			for i, want := range tc.expected {
				select {
				case got := <-hasErr:
					if want.hasErr == nil || got != *want.hasErr {
						t.Fatalf("event %d: got hasErr=%v, expected %+v", i, got, want)
					}
				case got := <-prs:
					if want.prs == nil || hasher.Hash(got) != hasher.Hash(want.prs) {
						t.Fatalf("event %d: got prs %v, expected %+v", i, got, want)
					}
				case <-time.After(2 * time.Second):
					t.Fatalf("event %d: timed out", i)
				}
			}

			cancel()

			select {
			case err := <-done:
				if err != nil {
					t.Fatalf("PollPRs returned %v after cancel, expected nil", err)
				}
			case <-time.After(2 * time.Second):
				t.Fatal("PollPRs did not stop after cancel")
			}
		})
	}
}

func TestPollPRsStopsWhileWaiting(t *testing.T) {
	conf := validConfig()
	conf.Poll.Frequency = 3600

	//exhaustruct:ignore
	p := Polling{
		Logger:    discardLogger(),
		GqlClient: &gittertest.Fake{},
		Config:    conf,
		queryAggregator: &scriptedQuery{
			results: []map[string][]gitter.PullRequest{{}},
			errs:    []error{nil},
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	prs := make(chan map[string][]gitter.PullRequest, 1)
	hasErr := make(chan bool, 1)
	done := make(chan error, 1)

	go func() { done <- p.PollPRs(ctx, prs, hasErr) }()

	<-hasErr // first round finished; now sleeping for an hour
	cancel()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("got %v, expected nil", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("PollPRs did not stop while waiting for the next poll")
	}
}

func TestPollPRsPreconditions(t *testing.T) {
	testCases := []struct {
		name        string
		polling     Polling
		expectedErr error
	}{
		//exhaustruct:ignore
		{name: "not set up", polling: Polling{Logger: discardLogger(), GqlClient: &gittertest.Fake{}}, expectedErr: ErrNotSetup},
		//exhaustruct:ignore
		{name: "nil client", polling: Polling{Logger: discardLogger(), queryAggregator: &scriptedQuery{}}, expectedErr: ErrNilClient},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.polling.PollPRs(context.Background(), make(chan map[string][]gitter.PullRequest), make(chan bool))
			if !errors.Is(err, tc.expectedErr) {
				t.Fatalf("got error %v, expected %v", err, tc.expectedErr)
			}
		})
	}
}

func TestGetPollDuration(t *testing.T) {
	testCases := []struct {
		seconds  int
		expected time.Duration
	}{
		{seconds: 0, expected: 0},
		{seconds: 30, expected: 30 * time.Second},
		{seconds: 3600, expected: time.Hour},
	}

	for _, tc := range testCases {
		conf := validConfig()
		conf.Poll.Frequency = tc.seconds

		//exhaustruct:ignore
		if got := (Polling{Config: conf}).getPollDuration(); got != tc.expected {
			t.Fatalf("frequency %d: got %v, expected %v", tc.seconds, got, tc.expected)
		}
	}
}

func TestPollPRsRefreshPollsEarly(t *testing.T) {
	conf := validConfig()
	conf.Poll.Frequency = 3600 // without a refresh, the second poll would be an hour away

	query := &scriptedQuery{
		results: []map[string][]gitter.PullRequest{{}},
		errs:    []error{nil},
	}

	refresh := make(chan struct{}, 1)

	//exhaustruct:ignore
	p := Polling{
		Logger:          discardLogger(),
		GqlClient:       &gittertest.Fake{},
		Config:          conf,
		Refresh:         refresh,
		queryAggregator: query,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	prs := make(chan map[string][]gitter.PullRequest, 1)
	hasErr := make(chan bool)

	go func() { _ = p.PollPRs(ctx, prs, hasErr) }()

	<-hasErr // first poll done; now waiting for an hour
	refresh <- struct{}{}

	select {
	case <-hasErr: // second poll happened straight away
	case <-time.After(2 * time.Second):
		t.Fatal("refresh did not trigger an early poll")
	}

	query.mu.Lock()
	calls := query.calls
	query.mu.Unlock()

	if calls != 2 {
		t.Fatalf("polled %d times, expected 2", calls)
	}
}
