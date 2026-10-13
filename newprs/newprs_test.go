package newprs_test

import (
	"errors"
	"slices"
	"testing"

	"appletini/gitter"
	"appletini/newprs"
	"appletini/trackerkey"
)

type prMap = map[string][]gitter.PullRequest

//exhaustruct:ignore
func pr(id string) gitter.PullRequest { return gitter.PullRequest{Id: id, Number: len(id)} }

func ids(found []newprs.Found) []string {
	out := make([]string, 0, len(found))
	for _, f := range found {
		out = append(out, f.Key+":"+f.PR.Id)
	}

	return out
}

func TestDetector(t *testing.T) {
	mine := pr("MINE")
	mine.ViewerAuthor = true

	approved := pr("APPROVED")
	approved.ViewerReview = "APPROVED"

	testCases := []struct {
		name     string
		polls    []prMap
		expected [][]string // expected new PRs after each poll
	}{
		{
			name:     "first poll is only a baseline",
			polls:    []prMap{{"repo_0": {pr("A"), pr("B")}}},
			expected: [][]string{{}},
		},
		{
			name: "a PR appearing later is new, once",
			polls: []prMap{
				{"repo_0": {pr("A")}},
				{"repo_0": {pr("A"), pr("B")}},
				{"repo_0": {pr("A"), pr("B")}},
			},
			expected: [][]string{{}, {"repo_0:B"}, {}},
		},
		{
			name: "label and author trackers count too",
			polls: []prMap{
				{},
				{"label_0": {pr("L")}, "author_1": {pr("U")}},
			},
			expected: [][]string{{}, {"author_1:U", "label_0:L"}},
		},
		{
			name: "your own PRs never notify",
			polls: []prMap{
				{},
				{"personal": {pr("P")}, "repo_0": {mine}},
			},
			expected: [][]string{{}, {}},
		},
		{
			name: "PRs you already approved don't notify",
			polls: []prMap{
				{},
				{"repo_0": {approved}},
			},
			expected: [][]string{{}, {}},
		},
		{
			name: "a PR in two trackers notifies once",
			polls: []prMap{
				{},
				{"repo_0": {pr("X")}, "label_0": {pr("X")}},
			},
			expected: [][]string{{}, {"label_0:X"}},
		},
		{
			name: "a partial poll after an error doesn't re-announce old PRs",
			polls: []prMap{
				{"repo_0": {pr("A")}, "label_0": {pr("B")}},
				{"repo_0": {pr("A")}}, // label query failed
				{"repo_0": {pr("A")}, "label_0": {pr("B")}},
			},
			expected: [][]string{{}, {}, {}},
		},
		{
			name: "a PR seen only in a personal group still isn't new later",
			polls: []prMap{
				{"personal": {pr("P")}},
				{"repo_0": {pr("P")}},
			},
			expected: [][]string{{}, {}},
		},
		{
			name: "PRs without an id are ignored",
			polls: []prMap{
				{},
				{"repo_0": {pr("")}},
			},
			expected: [][]string{{}, {}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var d newprs.Detector

			for i, poll := range tc.polls {
				found, err := d.Update(poll)
				if err != nil {
					t.Fatalf("poll %d: unexpected error: %v", i, err)
				}

				if got := ids(found); !slices.Equal(got, tc.expected[i]) {
					t.Fatalf("poll %d: got %v, expected %v", i, got, tc.expected[i])
				}
			}
		})
	}
}

func TestDetectorMalformedKeys(t *testing.T) {
	var d newprs.Detector

	_, err := d.Update(prMap{})
	if err != nil {
		t.Fatal(err)
	}

	found, err := d.Update(prMap{"garbage": {pr("G")}, "repo_0": {pr("R")}})
	if !errors.Is(err, trackerkey.ErrMalformed) {
		t.Fatalf("got error %v, expected %v", err, trackerkey.ErrMalformed)
	}

	if got := ids(found); !slices.Equal(got, []string{"repo_0:R"}) {
		t.Fatalf("valid groups should still be checked, got %v", got)
	}
}
