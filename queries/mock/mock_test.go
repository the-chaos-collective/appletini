package mock_test

import (
	"context"
	"errors"
	"testing"

	"appletini/queries/mock"
	"appletini/trackerkey"
)

func TestMockQuery(t *testing.T) {
	cancelled, cancel := context.WithCancel(context.Background())
	cancel()

	testCases := []struct {
		name         string
		ctx          context.Context
		expectedKeys []string
		expectedErr  error
	}{
		{name: "returns example data", ctx: context.Background(), expectedKeys: []string{"personal", "label_0", "repo_0", "author_0"}},
		{name: "cancelled", ctx: cancelled, expectedErr: context.Canceled},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := mock.MockQuery{}.GetAll(tc.ctx, nil)
			if !errors.Is(err, tc.expectedErr) {
				t.Fatalf("got error %v, expected %v", err, tc.expectedErr)
			}

			if len(got) != len(tc.expectedKeys) {
				t.Fatalf("got %d groups, expected %d", len(got), len(tc.expectedKeys))
			}

			for _, key := range tc.expectedKeys {
				if len(got[key]) == 0 {
					t.Fatalf("missing PRs for %s", key)
				}

				// Mock keys must be parseable by the UI.
				_, err := trackerkey.Parse(key)
				if err != nil {
					t.Fatalf("mock key %q is not a valid tracker key: %v", key, err)
				}
			}
		})
	}
}
