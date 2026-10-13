package by_author_test

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"appletini/gitter"
	"appletini/gitter/gittertest"
	"appletini/queries/by_author"
)

func validTracker() by_author.Tracker {
	return by_author.Tracker{
		Id:      "author_0",
		Authors: []string{"darvoid", "tobyselway"},
		Repo:    "appletini",
		Owner:   "the-chaos-collective",
		Title:   "Team",
	}
}

func validConfig() by_author.Config {
	return by_author.Config{
		Trackers:       []by_author.Tracker{validTracker()},
		ReviewAmount:   10,
		PrAmount:       10,
		CommentsAmount: 10,
	}
}

func TestMakeQueryValidation(t *testing.T) {
	testCases := []struct {
		name        string
		mutate      func(*by_author.Config)
		errContains string
	}{
		{name: "valid", mutate: func(*by_author.Config) {}},
		{name: "no trackers skips validation", mutate: func(c *by_author.Config) {
			c.Trackers = nil
			c.CommentsAmount = 0
		}},
		{name: "empty id", mutate: func(c *by_author.Config) { c.Trackers[0].Id = "" }, errContains: "Trackers[0].Identifier must not be empty"},
		{name: "empty repo", mutate: func(c *by_author.Config) { c.Trackers[0].Repo = "" }, errContains: "Trackers[0].Repo must not be empty"},
		{name: "empty owner", mutate: func(c *by_author.Config) { c.Trackers[0].Owner = "" }, errContains: "Trackers[0].Owner must not be empty"},
		{name: "no authors", mutate: func(c *by_author.Config) { c.Trackers[0].Authors = nil }, errContains: "Trackers[0].Authors must have at least one author"},
		{name: "zero review amount", mutate: func(c *by_author.Config) { c.ReviewAmount = 0 }, errContains: "invalid config: ReviewAmount must not be zero or less"},
		{name: "zero pr amount", mutate: func(c *by_author.Config) { c.PrAmount = 0 }, errContains: "invalid config: PrAmount must not be zero or less"},
		{name: "zero comments amount", mutate: func(c *by_author.Config) { c.CommentsAmount = 0 }, errContains: "invalid config: CommentsAmount must not be zero or less"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			conf := validConfig()
			tc.mutate(&conf)

			_, err := by_author.MakeQuery(conf)

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

func TestGeneratedQuery(t *testing.T) {
	query, err := by_author.MakeQuery(validConfig())
	if err != nil {
		t.Fatalf("MakeQuery: %v", err)
	}

	fake := &gittertest.Fake{Data: `{}`}

	_, err = query.GetAll(context.Background(), fake)
	if err != nil {
		t.Fatalf("GetAll: %v", err)
	}

	sent := fake.Queries()[0]
	for _, want := range []string{
		`author_0_0: search(query: "repo:the-chaos-collective/appletini is:open type:pr author:darvoid", type: ISSUE, first: 10)`,
		`author_0_1: search(query: "repo:the-chaos-collective/appletini is:open type:pr author:tobyselway", type: ISSUE, first: 10)`,
	} {
		if !strings.Contains(sent, want) {
			t.Fatalf("query does not contain %q:\n%s", want, sent)
		}
	}
}

func node(id string, number int) string {
	return fmt.Sprintf(`{"id": %q, "title": "T", "number": %d, "reviews": {"totalCount": 1}}`, id, number)
}

func TestGetAll(t *testing.T) {
	errBoom := errors.New("boom")

	// Two authors for tracker 0, one for tracker 1. Results must be merged per
	// tracker in author order regardless of map iteration order.
	response := `{
	  "author_0_1": {"nodes": [` + node("B1", 3) + `]},
	  "author_0_0": {"nodes": [` + node("A1", 1) + `, ` + node("A2", 2) + `]},
	  "author_1_0": {"nodes": []}
	}`

	//exhaustruct:ignore
	expected := map[string][]gitter.PullRequest{
		"author_0": {
			{Id: "A1", Title: "T", Number: 1, ReviewCount: 1},
			{Id: "A2", Title: "T", Number: 2, ReviewCount: 1},
			{Id: "B1", Title: "T", Number: 3, ReviewCount: 1},
		},
		"author_1": nil,
	}

	testCases := []struct {
		name        string
		client      gitter.Querier
		expected    map[string][]gitter.PullRequest
		expectedErr error
	}{
		{name: "merges authors in order", client: &gittertest.Fake{Data: response}, expected: expected},
		{
			name:        "malformed key",
			client:      &gittertest.Fake{Data: `{"author": {"nodes": []}}`},
			expected:    map[string][]gitter.PullRequest{},
			expectedErr: by_author.ErrMalformedKey,
		},
		{
			name:        "non-numeric author index",
			client:      &gittertest.Fake{Data: `{"author_0_x": {"nodes": []}}`},
			expected:    map[string][]gitter.PullRequest{},
			expectedErr: by_author.ErrMalformedKey,
		},
		{name: "client error", client: &gittertest.Fake{Err: errBoom}, expected: map[string][]gitter.PullRequest{}, expectedErr: errBoom},
		{name: "nil client", client: nil, expected: map[string][]gitter.PullRequest{}, expectedErr: by_author.ErrNilClient},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query, err := by_author.MakeQuery(validConfig())
			if err != nil {
				t.Fatalf("MakeQuery: %v", err)
			}

			// Repeat to catch map-order dependent results.
			for range 20 {
				got, err := query.GetAll(context.Background(), tc.client)
				if !errors.Is(err, tc.expectedErr) {
					t.Fatalf("got error %v, expected %v", err, tc.expectedErr)
				}

				if !reflect.DeepEqual(got, tc.expected) {
					t.Fatalf("got %+v, expected %+v", got, tc.expected)
				}
			}
		})
	}
}
