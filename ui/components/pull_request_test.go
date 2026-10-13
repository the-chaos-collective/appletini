package components_test

import (
	"strings"
	"testing"

	"appletini/status"
	"appletini/ui"
	"appletini/ui/components"
)

func TestPullRequestBuild(t *testing.T) {
	testCases := []struct {
		name          string
		emoji         []string
		expectedLines []string
	}{
		{
			name:          "all three emoji",
			emoji:         []string{"🔴", "🚧", "📝"},
			expectedLines: []string{"🔴\t(#7) Fix tray", "🚧\tfix/tray ➜ main", "📝\tmsg"},
		},
		{
			name:          "extra emoji share the status line",
			emoji:         []string{"🔴", "🚧", "📝", "❌"},
			expectedLines: []string{"🔴\t(#7) Fix tray", "🚧\tfix/tray ➜ main", "📝❌\tmsg"},
		},
		{
			name:          "one emoji",
			emoji:         []string{"🟢"},
			expectedLines: []string{"🟢\t(#7) Fix tray", "\tfix/tray ➜ main", "\tmsg"},
		},
		{
			name:          "no emoji",
			emoji:         nil,
			expectedLines: []string{"\t(#7) Fix tray", "\tfix/tray ➜ main", "\tmsg"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			item := components.PullRequest{
				Number:      7,
				Title:       "Fix tray",
				HeadRefName: "fix/tray",
				BaseRefName: "main",
				Permalink:   "https://github.com/o/r/pull/7",
				//exhaustruct:ignore
				Status: status.Status{Emoji: tc.emoji, Message: "msg"},
			}.Build()

			sub, ok := item.(ui.SystraySubmenu)
			if !ok {
				t.Fatalf("got %T, expected a submenu", item)
			}

			if got := strings.Split(sub.Title, "\n"); strings.Join(got, "|") != strings.Join(tc.expectedLines, "|") {
				t.Fatalf("got %q, expected %q", got, tc.expectedLines)
			}

			if len(sub.Submenu.Items) != 1 {
				t.Fatalf("expected one option, got %d", len(sub.Submenu.Items))
			}

			button, ok := sub.Submenu.Items[0].(ui.SystrayButton)
			if !ok || button.Title != "Open in browser" {
				t.Fatalf("unexpected option %+v", sub.Submenu.Items[0])
			}
		})
	}
}

func TestOpenInBrowserRejectsBadLinks(t *testing.T) {
	testCases := []struct {
		name      string
		permalink string
	}{
		{name: "empty", permalink: ""},
		{name: "not http", permalink: "file:///etc/passwd"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			menu := components.PullRequestOptions{Permalink: tc.permalink}.Build()

			button, ok := menu.Items[0].(ui.SystrayButton)
			if !ok {
				t.Fatalf("got %T", menu.Items[0])
			}

			err := button.Action()
			if err == nil || !strings.Contains(err.Error(), "error opening PR in browser") {
				t.Fatalf("got error %v", err)
			}
		})
	}
}

func TestApproveOption(t *testing.T) {
	approved := 0
	approve := func() error { approved++; return nil }

	testCases := []struct {
		name     string
		approve  func() error
		expected []string
	}{
		{name: "approve shown when available", approve: approve, expected: []string{"Open in browser", "Approve"}},
		{name: "approve hidden when nil", approve: nil, expected: []string{"Open in browser"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			menu := components.PullRequestOptions{Permalink: "https://x", Approve: tc.approve}.Build()

			got := make([]string, 0, len(menu.Items))
			for _, item := range menu.Items {
				button, ok := item.(ui.SystrayButton)
				if !ok {
					t.Fatalf("option is %T", item)
				}

				got = append(got, button.Title)

				if button.Title == "Approve" {
					err := button.Action()
					if err != nil {
						t.Fatal(err)
					}
				}
			}

			if strings.Join(got, ",") != strings.Join(tc.expected, ",") {
				t.Fatalf("got %v, expected %v", got, tc.expected)
			}
		})
	}

	if approved != 1 {
		t.Fatalf("approve ran %d times, expected 1", approved)
	}
}
