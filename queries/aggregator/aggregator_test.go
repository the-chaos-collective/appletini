package aggregator_test

import (
	"context"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/charmbracelet/log"

	"appletini/gitter"
	"appletini/logging"
	"appletini/queries"
	"appletini/queries/aggregator"
)

// stubQuery returns fixed results and counts how often it ran.
type stubQuery struct {
	result map[string][]gitter.PullRequest
	err    error
	calls  *int
}

func (s stubQuery) GetAll(_ context.Context, _ gitter.Querier) (map[string][]gitter.PullRequest, error) {
	*s.calls++

	return s.result, s.err
}

func TestGetAll(t *testing.T) {
	errBoom := errors.New("boom")
	logger := logging.NewCharm(log.New(io.Discard))

	//exhaustruct:ignore
	prA := gitter.PullRequest{Number: 1}
	//exhaustruct:ignore
	prB := gitter.PullRequest{Number: 2}

	cancelled, cancel := context.WithCancel(context.Background())
	cancel()

	testCases := []struct {
		name          string
		ctx           context.Context
		mock          bool
		results       []map[string][]gitter.PullRequest
		errs          []error
		expected      map[string][]gitter.PullRequest
		expectedErr   error
		expectedCalls []int
	}{
		{
			name:          "merges results",
			ctx:           context.Background(),
			results:       []map[string][]gitter.PullRequest{{"label_0": {prA}}, {"repo_0": {prB}}},
			errs:          []error{nil, nil},
			expected:      map[string][]gitter.PullRequest{"label_0": {prA}, "repo_0": {prB}},
			expectedCalls: []int{1, 1},
		},
		{
			name:          "stops at first error and keeps partial results",
			ctx:           context.Background(),
			results:       []map[string][]gitter.PullRequest{{"label_0": {prA}}, nil, {"repo_0": {prB}}},
			errs:          []error{nil, errBoom, nil},
			expected:      map[string][]gitter.PullRequest{"label_0": {prA}},
			expectedErr:   errBoom,
			expectedCalls: []int{1, 1, 0},
		},
		{
			name:          "cancelled before running",
			ctx:           cancelled,
			results:       []map[string][]gitter.PullRequest{{"label_0": {prA}}},
			errs:          []error{nil},
			expected:      map[string][]gitter.PullRequest{},
			expectedErr:   context.Canceled,
			expectedCalls: []int{0},
		},
		{
			name:          "no queries",
			ctx:           context.Background(),
			expected:      map[string][]gitter.PullRequest{},
			expectedCalls: []int{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			calls := make([]int, len(tc.results))
			qs := make([]queries.Query, len(tc.results))

			for i := range tc.results {
				qs[i] = stubQuery{result: tc.results[i], err: tc.errs[i], calls: &calls[i]}
			}

			qa := aggregator.QueryAggregator{Queries: qs, Mock: tc.mock, Logger: logger}

			got, err := qa.GetAll(tc.ctx, nil)
			if !errors.Is(err, tc.expectedErr) {
				t.Fatalf("got error %v, expected %v", err, tc.expectedErr)
			}

			if !reflect.DeepEqual(got, tc.expected) {
				t.Fatalf("got %+v, expected %+v", got, tc.expected)
			}

			if !reflect.DeepEqual(calls, tc.expectedCalls) {
				t.Fatalf("got calls %v, expected %v", calls, tc.expectedCalls)
			}
		})
	}
}

func TestGetAllMock(t *testing.T) {
	calls := 0
	qa := aggregator.QueryAggregator{
		Queries: []queries.Query{stubQuery{err: errors.New("must not run"), calls: &calls}},
		Mock:    true,
		Logger:  logging.NewCharm(log.New(io.Discard)),
	}

	got, err := qa.GetAll(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if calls != 0 {
		t.Fatalf("real queries ran %d times in mock mode", calls)
	}

	if len(got["personal"]) == 0 {
		t.Fatal("expected mock data")
	}
}
