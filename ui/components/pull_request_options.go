package components

import (
	"fmt"

	"git_applet/actions"
	"git_applet/ui"
)

type PullRequestOptions struct {
	Mergeable       string
	ReviewDecision  string
	RepoBranchRules struct{}
	Permalink       string
}

func (info PullRequestOptions) Build() ui.SystrayMenu {
	return ui.SystrayMenu{
		Items: []ui.Itemable{
			ui.SystrayButton{
				Title: "Open in browser",
				Action: func() error {
					err := actions.OpenLink(info.Permalink)
					if err != nil {
						return fmt.Errorf("error opening tracked PR in browser: %w", err)
					}
					return nil
				},
			},
		},
	}
}
