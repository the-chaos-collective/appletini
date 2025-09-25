package components

import (
	"fmt"
	"strings"

	"appletini/status"
	"appletini/ui"
)

type PullRequest struct {
	Number      int
	Title       string
	HeadRefName string
	BaseRefName string
	Permalink   string
	Status      status.Status
}

func (pr PullRequest) makeTitle() string {
	emoji := ""
	if len(pr.Status.Emoji) > 0 {
		emoji = pr.Status.Emoji[0]
	}
	return fmt.Sprintf("%s\t(#%d) %s", emoji, pr.Number, pr.Title)
}

func (pr PullRequest) makeBranchLine() string {
	emoji := ""
	if len(pr.Status.Emoji) > 1 {
		emoji = pr.Status.Emoji[1]
	}
	return fmt.Sprintf("%s\t%s ➜ %s", emoji, pr.HeadRefName, pr.BaseRefName)
}

func (pr PullRequest) makeStatus() string {
	emoji := ""
	if len(pr.Status.Emoji) > 2 {
		emoji = pr.Status.Emoji[2]
	}
	return fmt.Sprintf("%s\t%s", emoji, pr.Status.Message)
}

func (pr PullRequest) makeText() string {
	return strings.Join([]string{
		pr.makeTitle(),
		pr.makeBranchLine(),
		pr.makeStatus(),
	}, "\n")
}

func (pr PullRequest) Build() ui.Itemable {
	submenu := PullRequestOptions{
		Permalink: pr.Permalink,
	}.Build()

	return ui.SystraySubmenu{
		Title:   pr.makeText(),
		Submenu: &submenu,
	}
}
