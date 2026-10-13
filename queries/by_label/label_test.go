package by_label_test

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"appletini/gitter"
	"appletini/gitter/gittertest"
	"appletini/queries/by_label"
)

func validTracker() by_label.Tracker {
	return by_label.Tracker{Id: "label_0", Label: "bug", Repo: "appletini", Owner: "the-chaos-collective", Title: "Bugs"}
}

func validConfig() by_label.Config {
	return by_label.Config{
		Trackers:       []by_label.Tracker{validTracker()},
		ReviewAmount:   10,
		PrAmount:       10,
		CommentsAmount: 10,
	}
}

func TestMakeQueryValidation(t *testing.T) {
	testCases := []struct {
		name        string
		mutate      func(*by_label.Config)
		errContains string
	}{
		{name: "valid", mutate: func(*by_label.Config) {}},
		{name: "no trackers skips validation", mutate: func(c *by_label.Config) {
			c.Trackers = nil
			c.PrAmount = 0
		}},
		{name: "empty id", mutate: func(c *by_label.Config) { c.Trackers[0].Id = "" }, errContains: "Trackers[0].Identifier must not be empty"},
		{name: "empty label", mutate: func(c *by_label.Config) { c.Trackers[0].Label = "" }, errContains: "Trackers[0].Label must not be empty"},
		{name: "empty repo", mutate: func(c *by_label.Config) { c.Trackers[0].Repo = "" }, errContains: "Trackers[0].Repo must not be empty"},
		{name: "empty owner", mutate: func(c *by_label.Config) { c.Trackers[0].Owner = "" }, errContains: "Trackers[0].Owner must not be empty"},
		{name: "second tracker invalid", mutate: func(c *by_label.Config) {
			bad := validTracker()
			bad.Label = ""
			c.Trackers = append(c.Trackers, bad)
		}, errContains: "Trackers[1].Label must not be empty"},
		{name: "zero review amount", mutate: func(c *by_label.Config) { c.ReviewAmount = 0 }, errContains: "invalid config: ReviewAmount must not be zero or less"},
		{name: "zero pr amount", mutate: func(c *by_label.Config) { c.PrAmount = 0 }, errContains: "invalid config: PrAmount must not be zero or less"},
		{name: "negative comments amount", mutate: func(c *by_label.Config) { c.CommentsAmount = -1 }, errContains: "invalid config: CommentsAmount must not be zero or less"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			conf := validConfig()
			tc.mutate(&conf)

			_, err := by_label.MakeQuery(conf)

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

func TestGeneratedQueryQuotesValues(t *testing.T) {
	testCases := []struct {
		name     string
		label    string
		expected string
	}{
		{name: "plain", label: "bug", expected: `label(name: "bug")`},
		{name: "apostrophe is not html-escaped", label: "won't fix", expected: `label(name: "won't fix")`},
		{name: "ampersand is not html-escaped", label: "R&D", expected: `label(name: "R&D")`},
		{name: "double quote is escaped", label: `say "hi"`, expected: `label(name: "say \"hi\"")`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			conf := validConfig()
			conf.Trackers[0].Label = tc.label

			query, err := by_label.MakeQuery(conf)
			if err != nil {
				t.Fatalf("MakeQuery: %v", err)
			}

			fake := &gittertest.Fake{Data: `{}`}

			_, err = query.GetAll(context.Background(), fake)
			if err != nil {
				t.Fatalf("GetAll: %v", err)
			}

			sent := fake.Queries()[0]
			if !strings.Contains(sent, tc.expected) {
				t.Fatalf("query does not contain %s:\n%s", tc.expected, sent)
			}

			if !strings.Contains(sent, `repository(name: "appletini", owner: "the-chaos-collective")`) {
				t.Fatalf("query does not quote repository:\n%s", sent)
			}
		})
	}
}

const labelResponse = `{
  "label_0": {"label": {"pullRequests": {"edges": [
    {"node": {
      "id": "PR_1", "title": "Fix tray", "number": 7,
      "baseRefName": "main", "headRefName": "fix/tray",
      "permalink": "https://github.com/o/r/pull/7",
      "reviewDecision": "APPROVED", "mergeable": "MERGEABLE",
      "reviewRequests": {"totalCount": 2},
      "reviews": {"totalCount": 3}
    }}
  ]}}}
}`

func TestGetAll(t *testing.T) {
	errBoom := errors.New("boom")

	expectedPRs := map[string][]gitter.PullRequest{
		"label_0": {{
			Title:          "Fix tray",
			BaseRefName:    "main",
			HeadRefName:    "fix/tray",
			Number:         7,
			Permalink:      "https://github.com/o/r/pull/7",
			ReviewCount:    3,
			ReviewRequests: 2,
			ReviewDecision: "APPROVED",
			Id:             "PR_1",
			Mergeable:      "MERGEABLE",
		}},
	}

	testCases := []struct {
		name        string
		conf        by_label.Config
		client      gitter.Querier
		expected    map[string][]gitter.PullRequest
		expectedErr error
		noCalls     bool
	}{
		{name: "maps response", conf: validConfig(), client: &gittertest.Fake{Data: labelResponse}, expected: expectedPRs},
		{name: "client error", conf: validConfig(), client: &gittertest.Fake{Err: errBoom}, expected: map[string][]gitter.PullRequest{}, expectedErr: errBoom},
		{name: "nil client", conf: validConfig(), client: nil, expected: map[string][]gitter.PullRequest{}, expectedErr: by_label.ErrNilClient},
		{
			name:     "no trackers does not query",
			conf:     by_label.Config{ReviewAmount: 1, PrAmount: 1, CommentsAmount: 1},
			client:   &gittertest.Fake{Err: errBoom},
			expected: map[string][]gitter.PullRequest{},
			noCalls:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query, err := by_label.MakeQuery(tc.conf)
			if err != nil {
				t.Fatalf("MakeQuery: %v", err)
			}

			got, err := query.GetAll(context.Background(), tc.client)
			if !errors.Is(err, tc.expectedErr) {
				t.Fatalf("got error %v, expected %v", err, tc.expectedErr)
			}

			if !reflect.DeepEqual(got, tc.expected) {
				t.Fatalf("got %+v, expected %+v", got, tc.expected)
			}

			if fake, ok := tc.client.(*gittertest.Fake); ok && tc.noCalls && len(fake.Queries()) != 0 {
				t.Fatalf("expected no queries, got %d", len(fake.Queries()))
			}
		})
	}
}
