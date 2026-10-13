package reviews_test

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"appletini/gitter"
	"appletini/gitter/gittertest"
	"appletini/reviews"
)

const okResponse = `{"addPullRequestReview": {"pullRequestReview": {"id": "PRR_1", "url": "https://github.com/o/r/pull/7#pullrequestreview-1"}}}`

func TestApprove(t *testing.T) {
	errBoom := errors.New("Can not approve your own pull request")

	cancelled, cancel := context.WithCancel(context.Background())
	cancel()

	testCases := []struct {
		name         string
		ctx          context.Context
		client       gitter.Mutator
		prID         string
		comment      string
		expected     reviews.Review
		expectedVars map[string]any
		expectedErr  error
	}{
		{
			name:         "approves without a comment",
			ctx:          context.Background(),
			client:       &gittertest.Fake{Data: okResponse},
			prID:         "PR_7",
			expected:     reviews.Review{ID: "PRR_1", URL: "https://github.com/o/r/pull/7#pullrequestreview-1"},
			expectedVars: map[string]any{"id": "PR_7"},
		},
		{
			name:         "comment is sent as a variable, untouched",
			ctx:          context.Background(),
			client:       &gittertest.Fake{Data: okResponse},
			prID:         "PR_7",
			comment:      `LGTM "ship it" \o/`,
			expected:     reviews.Review{ID: "PRR_1", URL: "https://github.com/o/r/pull/7#pullrequestreview-1"},
			expectedVars: map[string]any{"id": "PR_7", "body": `LGTM "ship it" \o/`},
		},
		{
			name:        "GitHub error is returned",
			ctx:         context.Background(),
			client:      &gittertest.Fake{Err: errBoom},
			prID:        "PR_7",
			expectedErr: errBoom,
		},
		{
			name:        "no review in the response",
			ctx:         context.Background(),
			client:      &gittertest.Fake{Data: `{"addPullRequestReview": {"pullRequestReview": null}}`},
			prID:        "PR_7",
			expectedErr: reviews.ErrEmptyResponse,
		},
		{
			name:        "cancelled",
			ctx:         cancelled,
			client:      &gittertest.Fake{Data: okResponse},
			prID:        "PR_7",
			expectedErr: context.Canceled,
		},
		{name: "empty pull request id", ctx: context.Background(), client: &gittertest.Fake{}, prID: "", expectedErr: reviews.ErrEmptyID},
		{name: "nil client", ctx: context.Background(), client: nil, prID: "PR_7", expectedErr: reviews.ErrNilClient},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := reviews.Approver{Client: tc.client}.Approve(tc.ctx, tc.prID, tc.comment)
			if !errors.Is(err, tc.expectedErr) {
				t.Fatalf("got error %v, expected %v", err, tc.expectedErr)
			}

			if got != tc.expected {
				t.Fatalf("got %+v, expected %+v", got, tc.expected)
			}

			fake, ok := tc.client.(*gittertest.Fake)
			if !ok || tc.expectedVars == nil {
				return
			}

			if vars := fake.Vars(); len(vars) != 1 || !reflect.DeepEqual(vars[0], tc.expectedVars) {
				t.Fatalf("got vars %v, expected %v", vars, tc.expectedVars)
			}

			sent := fake.Queries()[0]
			if !strings.Contains(sent, "event: APPROVE") || strings.Contains(sent, tc.prID) {
				t.Fatalf("mutation should approve and take the id as a variable:\n%s", sent)
			}
		})
	}
}
