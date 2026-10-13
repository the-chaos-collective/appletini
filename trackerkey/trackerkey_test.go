package trackerkey_test

import (
	"errors"
	"slices"
	"testing"

	"appletini/trackerkey"
)

func TestNew(t *testing.T) {
	testCases := []struct {
		name        string
		kind        trackerkey.Kind
		idx         int
		expected    string
		expectedErr error
	}{
		{name: "personal ignores index", kind: trackerkey.Personal, idx: 3, expected: "personal"},
		{name: "label", kind: trackerkey.Label, idx: 0, expected: "label_0"},
		{name: "repo", kind: trackerkey.Repo, idx: 12, expected: "repo_12"},
		{name: "author", kind: trackerkey.Author, idx: 1, expected: "author_1"},
		{name: "negative index", kind: trackerkey.Label, idx: -1, expectedErr: trackerkey.ErrNegativeIndex},
		{name: "unknown kind", kind: "milestone", idx: 0, expectedErr: trackerkey.ErrUnknownKind},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := trackerkey.New(tc.kind, tc.idx)
			if !errors.Is(err, tc.expectedErr) {
				t.Fatalf("got error %v, expected %v", err, tc.expectedErr)
			}

			if got != tc.expected {
				t.Fatalf("got %q, expected %q", got, tc.expected)
			}
		})
	}
}

func TestParse(t *testing.T) {
	testCases := []struct {
		name        string
		raw         string
		expected    trackerkey.Key
		expectedErr error
	}{
		{name: "personal", raw: "personal", expected: trackerkey.Key{Kind: trackerkey.Personal}},
		{name: "label", raw: "label_0", expected: trackerkey.Key{Kind: trackerkey.Label, Index: 0}},
		{name: "repo multi-digit", raw: "repo_10", expected: trackerkey.Key{Kind: trackerkey.Repo, Index: 10}},
		{name: "author", raw: "author_2", expected: trackerkey.Key{Kind: trackerkey.Author, Index: 2}},
		{name: "empty", raw: "", expectedErr: trackerkey.ErrMalformed},
		{name: "no separator", raw: "label", expectedErr: trackerkey.ErrMalformed},
		{name: "missing index", raw: "label_", expectedErr: trackerkey.ErrMalformed},
		{name: "missing kind", raw: "_1", expectedErr: trackerkey.ErrMalformed},
		{name: "non-numeric index", raw: "label_x", expectedErr: trackerkey.ErrMalformed},
		{name: "extra suffix", raw: "author_0_1", expectedErr: trackerkey.ErrMalformed},
		{name: "negative index", raw: "repo_-1", expectedErr: trackerkey.ErrNegativeIndex},
		{name: "unknown kind", raw: "milestone_0", expectedErr: trackerkey.ErrUnknownKind},
		{name: "personal with index", raw: "personal_0", expectedErr: trackerkey.ErrMalformed},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := trackerkey.Parse(tc.raw)
			if !errors.Is(err, tc.expectedErr) {
				t.Fatalf("got error %v, expected %v", err, tc.expectedErr)
			}

			if got != tc.expected {
				t.Fatalf("got %+v, expected %+v", got, tc.expected)
			}
		})
	}
}

func TestNewParseRoundTrip(t *testing.T) {
	for _, kind := range []trackerkey.Kind{trackerkey.Label, trackerkey.Repo, trackerkey.Author} {
		for _, idx := range []int{0, 1, 9, 10, 123} {
			raw, err := trackerkey.New(kind, idx)
			if err != nil {
				t.Fatalf("New(%s, %d): %v", kind, idx, err)
			}

			key, err := trackerkey.Parse(raw)
			if err != nil {
				t.Fatalf("Parse(%q): %v", raw, err)
			}

			if key.Kind != kind || key.Index != idx {
				t.Fatalf("round trip of %s/%d gave %+v", kind, idx, key)
			}
		}
	}
}

func TestCompare(t *testing.T) {
	keys := []trackerkey.Key{
		{Kind: trackerkey.Author, Index: 0},
		{Kind: trackerkey.Repo, Index: 10},
		{Kind: trackerkey.Label, Index: 1},
		{Kind: trackerkey.Repo, Index: 2},
		{Kind: trackerkey.Personal},
		{Kind: trackerkey.Label, Index: 0},
	}

	expected := []trackerkey.Key{
		{Kind: trackerkey.Personal},
		{Kind: trackerkey.Label, Index: 0},
		{Kind: trackerkey.Label, Index: 1},
		{Kind: trackerkey.Repo, Index: 2},
		{Kind: trackerkey.Repo, Index: 10},
		{Kind: trackerkey.Author, Index: 0},
	}

	slices.SortFunc(keys, trackerkey.Compare)

	if !slices.Equal(keys, expected) {
		t.Fatalf("got %+v, expected %+v", keys, expected)
	}
}
