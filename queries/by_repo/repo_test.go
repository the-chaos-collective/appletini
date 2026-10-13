package by_repo_test

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"appletini/gitter"
	"appletini/gitter/gittertest"
	"appletini/queries/by_repo"
)

func validTracker() by_repo.Tracker {
	return by_repo.Tracker{Id: "repo_0", Repo: "appletini", Owner: "the-chaos-collective", Title: "Appletini"}
}

func validConfig() by_repo.Config {
	return by_repo.Config{
		Trackers:       []by_repo.Tracker{validTracker()},
		ReviewAmount:   10,
		PrAmount:       10,
		CommentsAmount: 10,
	}
}

func TestMakeQueryValidation(t *testing.T) {
	testCases := []struct {
		name        string
		mutate      func(*by_repo.Config)
		errContains string
	}{
		{name: "valid", mutate: func(*by_repo.Config) {}},
		{name: "no trackers skips validation", mutate: func(c *by_repo.Config) {
			c.Trackers = nil
			c.ReviewAmount = 0
		}},
		{name: "empty title", mutate: func(c *by_repo.Config) { c.Trackers[0].Title = "" }, errContains: "Trackers[0].Title must not be empty"},
		{name: "empty id", mutate: func(c *by_repo.Config) { c.Trackers[0].Id = "" }, errContains: "Trackers[0].Id must not be empty"},
		{name: "empty repo", mutate: func(c *by_repo.Config) { c.Trackers[0].Repo = "" }, errContains: "Trackers[0].Name must not be empty"},
		{name: "empty owner", mutate: func(c *by_repo.Config) { c.Trackers[0].Owner = "" }, errContains: "Trackers[0].Owner must not be empty"},
		{name: "zero review amount", mutate: func(c *by_repo.Config) { c.ReviewAmount = 0 }, errContains: "invalid config: ReviewAmount must not be zero or less"},
		{name: "zero pr amount", mutate: func(c *by_repo.Config) { c.PrAmount = 0 }, errContains: "invalid config: PrAmount must not be zero or less"},
		{name: "zero comments amount", mutate: func(c *by_repo.Config) { c.CommentsAmount = 0 }, errContains: "invalid config: CommentsAmount must not be zero or less"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			conf := validConfig()
			tc.mutate(&conf)

			_, err := by_repo.MakeQuery(conf)

			if tc.errContains == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				return
			}

			if err == nil || !strings.Contains(err.Error(), tc.errContains) {
				t.Fatalf("got error %v, expected it to contain %q", err, tc.errContains)
			}
		})
	}
}

const repoResponse = `{
  "repo_0": {"pullRequests": {"edges": [
    {"node": {"id": "PR_1", "title": "One", "number": 1, "reviewDecision": "", "mergeable": "UNKNOWN",
              "reviewRequests": {"totalCount": 0}, "reviews": {"totalCount": 0}}},
    {"node": {"id": "PR_2", "title": "Two", "number": 2, "reviewDecision": "CHANGES_REQUESTED", "mergeable": "CONFLICTING",
              "baseRefName": "main", "headRefName": "feat", "permalink": "https://x/2",
              "reviewRequests": {"totalCount": 1}, "reviews": {"totalCount": 4}}}
  ]}},
  "repo_1": {"pullRequests": {"edges": []}}
}`

func TestGetAll(t *testing.T) {
	errBoom := errors.New("boom")

	twoTrackers := validConfig()
	second := validTracker()
	second.Id = "repo_1"
	twoTrackers.Trackers = append(twoTrackers.Trackers, second)

	//exhaustruct:ignore
	expected := map[string][]gitter.PullRequest{
		"repo_0": {
			{Id: "PR_1", Title: "One", Number: 1, Mergeable: "UNKNOWN"},
			{
				Id: "PR_2", Title: "Two", Number: 2, ReviewDecision: "CHANGES_REQUESTED", Mergeable: "CONFLICTING",
				BaseRefName: "main", HeadRefName: "feat", Permalink: "https://x/2", ReviewRequests: 1, ReviewCount: 4,
			},
		},
		"repo_1": {},
	}

	cancelled, cancel := context.WithCancel(context.Background())
	cancel()

	testCases := []struct {
		name        string
		ctx         context.Context
		client      gitter.Querier
		expected    map[string][]gitter.PullRequest
		expectedErr error
	}{
		{name: "maps response", ctx: context.Background(), client: &gittertest.Fake{Data: repoResponse}, expected: expected},
		{name: "client error", ctx: context.Background(), client: &gittertest.Fake{Err: errBoom}, expected: map[string][]gitter.PullRequest{}, expectedErr: errBoom},
		{name: "cancelled", ctx: cancelled, client: &gittertest.Fake{Data: repoResponse}, expected: map[string][]gitter.PullRequest{}, expectedErr: context.Canceled},
		{name: "nil client", ctx: context.Background(), client: nil, expected: map[string][]gitter.PullRequest{}, expectedErr: by_repo.ErrNilClient},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query, err := by_repo.MakeQuery(twoTrackers)
			if err != nil {
				t.Fatalf("MakeQuery: %v", err)
			}

			got, err := query.GetAll(tc.ctx, tc.client)
			if !errors.Is(err, tc.expectedErr) {
				t.Fatalf("got error %v, expected %v", err, tc.expectedErr)
			}

			if !reflect.DeepEqual(got, tc.expected) {
				t.Fatalf("got %+v, expected %+v", got, tc.expected)
			}
		})
	}
}

func TestGeneratedQuery(t *testing.T) {
	query, err := by_repo.MakeQuery(validConfig())
	if err != nil {
		t.Fatalf("MakeQuery: %v", err)
	}

	fake := &gittertest.Fake{Data: `{}`}

	_, err = query.GetAll(context.Background(), fake)
	if err != nil {
		t.Fatalf("GetAll: %v", err)
	}

	for _, want := range []string{
		`repo_0: repository(name: "appletini", owner: "the-chaos-collective")`,
		"first: 10",
		"totalCount",
	} {
		if !strings.Contains(fake.Queries()[0], want) {
			t.Fatalf("query does not contain %q:\n%s", want, fake.Queries()[0])
		}
	}
}
