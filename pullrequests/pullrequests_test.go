package pullrequests_test

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"appletini/gitter"
	"appletini/gitter/gittertest"
	"appletini/pullrequests"
)

type action func(pullrequests.Actions, context.Context, string) (pullrequests.PullRequest, error)

func merge(a pullrequests.Actions, ctx context.Context, id string) (pullrequests.PullRequest, error) {
	return a.Merge(ctx, id)
}

func closePR(a pullrequests.Actions, ctx context.Context, id string) (pullrequests.PullRequest, error) {
	return a.Close(ctx, id)
}

func TestActions(t *testing.T) {
	errGitHub := errors.New("Pull request is not mergeable")

	cancelled, cancel := context.WithCancel(context.Background())
	cancel()

	testCases := []struct {
		name         string
		run          action
		ctx          context.Context
		client       gitter.Mutator
		prID         string
		expected     pullrequests.PullRequest
		expectedErr  error
		expectedDoc  string
		expectedVars map[string]any
	}{
		{
			name:         "merge",
			run:          merge,
			ctx:          context.Background(),
			client:       &gittertest.Fake{Data: `{"mergePullRequest": {"pullRequest": {"number": 7, "merged": true, "url": "https://x/7"}}}`},
			prID:         "PR_7",
			expected:     pullrequests.PullRequest{Number: 7, Merged: true, URL: "https://x/7"},
			expectedDoc:  "mergePullRequest(input: { pullRequestId: $id })",
			expectedVars: map[string]any{"id": "PR_7"},
		},
		{
			name:        "merge reported as not merged",
			run:         merge,
			ctx:         context.Background(),
			client:      &gittertest.Fake{Data: `{"mergePullRequest": {"pullRequest": {"number": 7, "merged": false}}}`},
			prID:        "PR_7",
			expected:    pullrequests.PullRequest{Number: 7},
			expectedErr: pullrequests.ErrNotMerged,
		},
		{
			name:        "merge rejected by GitHub",
			run:         merge,
			ctx:         context.Background(),
			client:      &gittertest.Fake{Err: errGitHub},
			prID:        "PR_7",
			expectedErr: errGitHub,
		},
		{
			name:        "merge with no pull request in the response",
			run:         merge,
			ctx:         context.Background(),
			client:      &gittertest.Fake{Data: `{"mergePullRequest": {"pullRequest": null}}`},
			prID:        "PR_7",
			expectedErr: pullrequests.ErrEmptyResponse,
		},
		{
			name:         "close",
			run:          closePR,
			ctx:          context.Background(),
			client:       &gittertest.Fake{Data: `{"closePullRequest": {"pullRequest": {"number": 7, "closed": true, "url": "https://x/7"}}}`},
			prID:         "PR_7",
			expected:     pullrequests.PullRequest{Number: 7, Closed: true, URL: "https://x/7"},
			expectedDoc:  "closePullRequest(input: { pullRequestId: $id })",
			expectedVars: map[string]any{"id": "PR_7"},
		},
		{
			name:        "close reported as still open",
			run:         closePR,
			ctx:         context.Background(),
			client:      &gittertest.Fake{Data: `{"closePullRequest": {"pullRequest": {"number": 7, "closed": false}}}`},
			prID:        "PR_7",
			expected:    pullrequests.PullRequest{Number: 7},
			expectedErr: pullrequests.ErrNotClosed,
		},
		{
			name:        "close with no pull request in the response",
			run:         closePR,
			ctx:         context.Background(),
			client:      &gittertest.Fake{Data: `{"closePullRequest": {"pullRequest": null}}`},
			prID:        "PR_7",
			expectedErr: pullrequests.ErrEmptyResponse,
		},
		{name: "cancelled", run: merge, ctx: cancelled, client: &gittertest.Fake{Data: `{}`}, prID: "PR_7", expectedErr: context.Canceled},
		{name: "empty id (merge)", run: merge, ctx: context.Background(), client: &gittertest.Fake{}, prID: "", expectedErr: pullrequests.ErrEmptyID},
		{name: "empty id (close)", run: closePR, ctx: context.Background(), client: &gittertest.Fake{}, prID: "", expectedErr: pullrequests.ErrEmptyID},
		{name: "nil client", run: closePR, ctx: context.Background(), client: nil, prID: "PR_7", expectedErr: pullrequests.ErrNilClient},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.run(pullrequests.Actions{Client: tc.client}, tc.ctx, tc.prID)
			if !errors.Is(err, tc.expectedErr) {
				t.Fatalf("got error %v, expected %v", err, tc.expectedErr)
			}

			if got != tc.expected {
				t.Fatalf("got %+v, expected %+v", got, tc.expected)
			}

			fake, ok := tc.client.(*gittertest.Fake)
			if !ok || tc.expectedDoc == "" {
				return
			}

			sent := fake.Queries()[0]
			if !strings.Contains(sent, tc.expectedDoc) || strings.Contains(sent, tc.prID) {
				t.Fatalf("mutation should contain %q and take the id as a variable:\n%s", tc.expectedDoc, sent)
			}

			if vars := fake.Vars()[0]; !reflect.DeepEqual(vars, tc.expectedVars) {
				t.Fatalf("got vars %v, expected %v", vars, tc.expectedVars)
			}
		})
	}
}
