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

	if info.Approve != nil {
		items = append(items, ui.SystrayButton{
			Title:  "Approve",
			Action: info.Approve,
		})
	}

	return ui.SystrayMenu{Items: items}
}
