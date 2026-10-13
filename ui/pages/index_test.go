package pages

import (
	"errors"
	"testing"

	"appletini/config"
	v2 "appletini/config/migration/v2"
	"appletini/gitter"
	"appletini/trackerkey"
	"appletini/ui"
	"appletini/ui/icons"

	"fyne.io/fyne/v2"
)

func trackers() config.Tracking {
	return config.Tracking{
		Personal: true,
		ByLabel:  v2.LabeledSet{{Title: "Bugs"}, {Title: "Docs"}},
		ByRepo:   v2.RepoSet{{Title: "Appletini"}},
		ByAuthor: v2.AuthorSet{{Title: "Team"}},
	}
}

func TestGroupTitle(t *testing.T) {
	testCases := []struct {
		name        string
		key         trackerkey.Key
		expected    string
		expectedErr error
	}{
		{name: "personal", key: trackerkey.Key{Kind: trackerkey.Personal}, expected: "My Pull Requests"},
		{name: "label 0", key: trackerkey.Key{Kind: trackerkey.Label, Index: 0}, expected: "Bugs"},
		{name: "label 1", key: trackerkey.Key{Kind: trackerkey.Label, Index: 1}, expected: "Docs"},
		{name: "repo", key: trackerkey.Key{Kind: trackerkey.Repo, Index: 0}, expected: "Appletini"},
		{name: "author", key: trackerkey.Key{Kind: trackerkey.Author, Index: 0}, expected: "Team"},
		{name: "label out of range", key: trackerkey.Key{Kind: trackerkey.Label, Index: 2}, expectedErr: ErrTrackerOutOfRange},
		{name: "repo out of range", key: trackerkey.Key{Kind: trackerkey.Repo, Index: 1}, expectedErr: ErrTrackerOutOfRange},
		{name: "author out of range", key: trackerkey.Key{Kind: trackerkey.Author, Index: 5}, expectedErr: ErrTrackerOutOfRange},
		{name: "unknown kind", key: trackerkey.Key{Kind: "milestone"}, expectedErr: ErrUnsupportedKind},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := groupTitle(tc.key, trackers())
			if !errors.Is(err, tc.expectedErr) {
				t.Fatalf("got error %v, expected %v", err, tc.expectedErr)
			}

			if got != tc.expected {
				t.Fatalf("got %q, expected %q", got, tc.expected)
			}
		})
	}
}

func TestIconFor(t *testing.T) {
	testCases := []struct {
		name     string
		state    iconState
		dark     bool
		expected fyne.Resource
	}{
		{name: "light default", state: iconState{}, expected: icons.ResIconDefault},
		{name: "light reviewable", state: iconState{red: true}, expected: icons.ResIconReviewable},
		{name: "light mergeable", state: iconState{green: true}, expected: icons.ResIconMergeable},
		{name: "light both", state: iconState{green: true, red: true}, expected: icons.ResIconBoth},
		{name: "dark default", state: iconState{}, dark: true, expected: icons.ResIconDefaultDark},
		{name: "dark both", state: iconState{green: true, red: true}, dark: true, expected: icons.ResIconBothDark},
		{name: "error wins (light)", state: iconState{green: true, err: true}, expected: icons.ResIconWarning},
		{name: "error wins (dark)", state: iconState{red: true, err: true}, dark: true, expected: icons.ResIconWarning},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := iconFor(tc.state, tc.dark)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tc.expected {
				t.Fatalf("got %s, expected %s", got.Name(), tc.expected.Name())
			}
		})
	}
}

func TestIconForCoversEveryState(t *testing.T) {
	for _, dark := range []bool{false, true} {
		for i := range 8 {
			state := iconState{green: i&1 != 0, red: i&2 != 0, err: i&4 != 0}

			_, err := iconFor(state, dark)
			if err != nil {
				t.Fatalf("dark=%v state=%+v: %v", dark, state, err)
			}
		}
	}
}

func titles(items []ui.Itemable) []string {
	out := make([]string, 0, len(items))

	for _, item := range items {
		switch v := item.(type) {
		case ui.SystraySubmenu:
			out = append(out, v.Title)
		case ui.SystrayButton:
			out = append(out, v.Title)
		case ui.SystraySeparator:
			out = append(out, "-")
		}
	}

	return out
}

func TestBuildTree(t *testing.T) {
	//exhaustruct:ignore
	mergeable := gitter.PullRequest{Number: 1, ReviewDecision: "APPROVED", Mergeable: "MERGEABLE"}
	//exhaustruct:ignore
	conflicting := gitter.PullRequest{Number: 2, ReviewDecision: "APPROVED", Mergeable: "CONFLICTING"}
	//exhaustruct:ignore
	waiting := gitter.PullRequest{Number: 3, ReviewDecision: "REVIEW_REQUIRED", Mergeable: "MERGEABLE"}

	//exhaustruct:ignore
	failingCI := gitter.PullRequest{Number: 4, ReviewDecision: "APPROVED", Mergeable: "MERGEABLE", ChecksState: "FAILURE"}
	//exhaustruct:ignore
	runningCI := gitter.PullRequest{Number: 5, ReviewDecision: "APPROVED", Mergeable: "MERGEABLE", ChecksState: "PENDING"}

	quit := func() error { return nil }

	testCases := []struct {
		name           string
		prs            map[string][]gitter.PullRequest
		expectedTitles []string
		expectedState  iconState
		expectErr      bool
	}{
		{
			name:           "no results still has quit",
			prs:            map[string][]gitter.PullRequest{},
			expectedTitles: []string{"-", "Quit"},
		},
		{
			name: "groups are ordered personal, label, repo, author by index",
			prs: map[string][]gitter.PullRequest{
				"author_0": nil, "repo_0": nil, "label_1": nil, "label_0": nil, "personal": nil,
			},
			expectedTitles: []string{"My Pull Requests", "Bugs", "Docs", "Appletini", "Team", "-", "Quit"},
		},
		{
			name:           "own mergeable PR turns the icon green",
			prs:            map[string][]gitter.PullRequest{"personal": {mergeable, waiting}},
			expectedTitles: []string{"My Pull Requests", "-", "Quit"},
			expectedState:  iconState{green: true},
		},
		{
			name:           "own PR with failing CI turns the icon red",
			prs:            map[string][]gitter.PullRequest{"personal": {failingCI}},
			expectedTitles: []string{"My Pull Requests", "-", "Quit"},
			expectedState:  iconState{red: true},
		},
		{
			name:           "own PR with running CI is not ready yet",
			prs:            map[string][]gitter.PullRequest{"personal": {runningCI}},
			expectedTitles: []string{"My Pull Requests", "-", "Quit"},
		},
		{
			name:           "own conflicting PR turns the icon red",
			prs:            map[string][]gitter.PullRequest{"personal": {conflicting}},
			expectedTitles: []string{"My Pull Requests", "-", "Quit"},
			expectedState:  iconState{red: true},
		},
		{
			name:           "other trackers do not affect the icon",
			prs:            map[string][]gitter.PullRequest{"repo_0": {mergeable, conflicting}},
			expectedTitles: []string{"Appletini", "-", "Quit"},
		},
		{
			name: "bad groups are skipped and reported, the rest is built",
			prs: map[string][]gitter.PullRequest{
				"repo_0": nil, "repo_7": nil, "garbage": nil,
			},
			expectedTitles: []string{"Appletini", "-", "Quit"},
			expectErr:      true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Repeat: map iteration order must not affect the result.
			for range 20 {
				items, state, err := buildTree(tc.prs, trackers(), quit)

				if (err != nil) != tc.expectErr {
					t.Fatalf("got error %v, expectErr=%v", err, tc.expectErr)
				}

				got := titles(items)
				if len(got) != len(tc.expectedTitles) {
					t.Fatalf("got %v, expected %v", got, tc.expectedTitles)
				}

				for i := range got {
					if got[i] != tc.expectedTitles[i] {
						t.Fatalf("got %v, expected %v", got, tc.expectedTitles)
					}
				}

				if state != tc.expectedState {
					t.Fatalf("got state %+v, expected %+v", state, tc.expectedState)
				}

				if state.err {
					t.Fatal("buildTree must never set the error flag")
				}
			}
		})
	}
}

func TestBuildTreeGroupContents(t *testing.T) {
	//exhaustruct:ignore
	prs := map[string][]gitter.PullRequest{"label_0": {{Number: 1, Title: "a"}, {Number: 2, Title: "b"}}}

	items, _, err := buildTree(prs, trackers(), func() error { return nil })
	if err != nil {
		t.Fatal(err)
	}

	group, ok := items[0].(ui.SystraySubmenu)
	if !ok {
		t.Fatalf("first item is %T", items[0])
	}

	if len(group.Submenu.Items) != 2 {
		t.Fatalf("group has %d PRs, expected 2", len(group.Submenu.Items))
	}
}
