package components

import (
	"fmt"

	"appletini/actions"
	"appletini/ui"
)

type PullRequestOptions struct {
	Permalink string
	// Approve approves the PR. When nil (e.g. for your own PRs, which GitHub
	// won't let you approve) the Approve option is not shown.
	Approve func() error
	// Merge and Close are shown only when non-nil.
	Merge func() error
	Close func() error
}

func (info PullRequestOptions) Build() ui.SystrayMenu {
	items := []ui.Itemable{
		ui.SystrayButton{
			Title: "Open in browser",
			Action: func() error {
				err := actions.OpenLink(info.Permalink)
				if err != nil {
					return fmt.Errorf("error opening PR in browser: %w", err)
				}
				return nil
			},
		},
	}

	optional := []struct {
		title  string
		action func() error
	}{
		{title: "Approve", action: info.Approve},
		{title: "Merge…", action: info.Merge},
		{title: "Close…", action: info.Close},
	}

	for _, opt := range optional {
		if opt.action != nil {
			items = append(items, ui.SystrayButton{Title: opt.title, Action: opt.action})
		}
	}

	return ui.SystrayMenu{Items: items}
}
