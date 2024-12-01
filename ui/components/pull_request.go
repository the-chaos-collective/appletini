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
	Mergeable   status.MergeableState
	Review      status.ReviewState
	HeadRefName string
	BaseRefName string
	Permalink   string
}

func (pr PullRequest) status() status.PRInfo {
	return status.PRInfo{
		Review:    pr.Review,
		Mergeable: pr.Mergeable,
	}
}

func (pr PullRequest) makeTitle(s status.Status) string {
	emoji := ""
	if len(s.Emoji) > 0 {
		emoji = s.Emoji[0]
	}
	return fmt.Sprintf("%s\t(#%d) %s", emoji, pr.Number, pr.Title)
}

func (pr PullRequest) makeBranchLine(s status.Status) string {
	emoji := ""
	if len(s.Emoji) > 1 {
		emoji = s.Emoji[1]
	}
	return fmt.Sprintf("%s\t%s ➜ %s", emoji, pr.HeadRefName, pr.BaseRefName)
}

func (pr PullRequest) makeStatus(s status.Status) string {
	emoji := ""
	if len(s.Emoji) > 2 {
		emoji = s.Emoji[2]
	}
	return fmt.Sprintf("%s\t%s", emoji, s.Message)
}

func (pr PullRequest) makeText() string {
	s := pr.status().Classify()

	return strings.Join([]string{
		pr.makeTitle(s),
		pr.makeBranchLine(s),
		pr.makeStatus(s),
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
