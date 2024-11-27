package components

import (
	"log"

	"git_applet/actions"
	"git_applet/ui"
)

type PullRequestOptions struct {
	Mergeable       string
	ReviewDecision  string
	RepoBranchRules struct{}
	Permalink       string
}

func (box PullRequestOptions) Build() ui.SystrayMenu {
	return ui.SystrayMenu{
		Items: []ui.Itemable{
			ui.SystrayButton{
				Title: "Open in browser",
				Action: func() {
					err := actions.OpenLink(box.Permalink)
					if err != nil {
						log.Fatalf("error opening tracked PR in browser: %w", err)
					}
				},
			},
		},
	}
}
