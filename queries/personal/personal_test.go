package personal_test

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"appletini/gitter"
	"appletini/gitter/gittertest"
	"appletini/queries/personal"
)

const viewerResponse = `{
  "viewer": {"pullRequests": {"edges": [
    {"node": {
      "id": "PR_9", "title": "Mine", "number": 9,
      "baseRefName": "main", "headRefName": "me/feature",
      "permalink": "https://github.com/o/r/pull/9",
      "reviewDecision": "REVIEW_REQUIRED", "mergeable": "MERGEABLE",
      "reviewRequests": {"totalCount": 1},
      "reviews": {"totalCount": 5, "nodes": []},
      "commits": {"nodes": [{"commit": {"statusCheckRollup": {"state": "PENDING"}}}]}
    }}
  ]}}
}`

func TestGetAll(t *testing.T) {
	errBoom := errors.New("boom")

	cancelled, cancel := context.WithCancel(context.Background())
	cancel()

	testCases := []struct {
		name        string
		ctx         context.Context
		client      gitter.Querier
		expected    map[string][]gitter.PullRequest
		expectedErr error
	}{
		{
			name:   "maps response under the personal key",
			ctx:    context.Background(),
			client: &gittertest.Fake{Data: viewerResponse},
			expected: map[string][]gitter.PullRequest{
				personal.Key: {{
					Id:             "PR_9",
					Title:          "Mine",
					Number:         9,
					BaseRefName:    "main",
					HeadRefName:    "me/feature",
					Permalink:      "https://github.com/o/r/pull/9",
					ReviewDecision: "REVIEW_REQUIRED",
					Mergeable:      "MERGEABLE",
					ReviewRequests: 1,
					ReviewCount:    5,
					ChecksState:    "PENDING",
				}},
			},
		},
		{
			name:     "no open PRs",
			ctx:      context.Background(),
			client:   &gittertest.Fake{Data: `{"viewer": {"pullRequests": {"edges": []}}}`},
			expected: map[string][]gitter.PullRequest{personal.Key: {}},
		},
		{name: "client error", ctx: context.Background(), client: &gittertest.Fake{Err: errBoom}, expected: map[string][]gitter.PullRequest{}, expectedErr: errBoom},
		{name: "cancelled", ctx: cancelled, client: &gittertest.Fake{Data: viewerResponse}, expected: map[string][]gitter.PullRequest{}, expectedErr: context.Canceled},
		{name: "nil client", ctx: context.Background(), client: nil, expected: map[string][]gitter.PullRequest{}, expectedErr: personal.ErrNilClient},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := personal.Query{}.GetAll(tc.ctx, tc.client)
			if !errors.Is(err, tc.expectedErr) {
				t.Fatalf("got error %v, expected %v", err, tc.expectedErr)
			}

			if !reflect.DeepEqual(got, tc.expected) {
				t.Fatalf("got %+v, expected %+v", got, tc.expected)
			}
		})
	}
}

func TestViewerQueryRequestsChecks(t *testing.T) {
	if !strings.Contains(personal.ViewerQuery, "statusCheckRollup") {
		t.Fatal("the personal query does not request the CI state")
	}
}
