package components

import (
	"fmt"
	"strings"

	"git_applet/ui"
)

type PullRequest struct {
	Number          int
	Title           string
	Mergeable       string
	ReviewDecision  string
	HeadRefName     string
	BaseRefName     string
	RepoBranchRules struct{}
	Permalink       string
}

type Status struct {
	Mergeable      string
	ReviewDecision string
}

var statusMap map[Status]func(PullRequest) string = map[Status]func(PullRequest) string{
	{Mergeable: "MERGEABLE", ReviewDecision: "APPROVED"}: func(pr PullRequest) string {
		return "Go ahead, merge it. I'm sure it's fine"
	},
	{Mergeable: "MERGEABLE", ReviewDecision: "CHANGES_REQUESTED"}: func(pr PullRequest) string {
		return "Go make your team happy"
	},
	{Mergeable: "MERGEABLE", ReviewDecision: "REVIEW_REQUIRED"}: func(pr PullRequest) string {
		return "Go ask for reviews"
	},
	{Mergeable: "CONFLICTING", ReviewDecision: "APPROVED"}: func(pr PullRequest) string {
		return "Git doesn't like you"
	},
	{Mergeable: "CONFLICTING", ReviewDecision: "CHANGES_REQUESTED"}: func(pr PullRequest) string {
		return "Nothing works, just give up"
	},
	{Mergeable: "UNKWNOWN", ReviewDecision: "APPROVED"}: func(pr PullRequest) string {
		return "Bro, your PR is cooked"
	},
	{Mergeable: "UNKWNOWN", ReviewDecision: "CHANGES_REQUESTED"}: func(pr PullRequest) string {
		return "The dog is dead but we can keep working on it"
	},
}

func (status Status) Infer(pr PullRequest) string {
	text := fmt.Sprintf("UNKNOWN (%s/%s)", pr.Mergeable, pr.ReviewDecision)

	res, ok := statusMap[status]
	if ok {
		text = res(pr)
	}

	return fmt.Sprintf("🔹 %s", text)
}

func (pr PullRequest) makeTitle() string {
	return fmt.Sprintf("🔷 (#%d) %s", pr.Number, pr.Title)
}

func (pr PullRequest) makeBranchLine() string {
	return fmt.Sprintf("%s ↦ %s", pr.HeadRefName, pr.BaseRefName)
}

func (pr PullRequest) makeStatus() string {
	return Status{
		Mergeable:      pr.Mergeable,
		ReviewDecision: pr.ReviewDecision,
	}.Infer(pr)
}

func (pr PullRequest) makeFullTitle() string {
	return strings.Join([]string{
		pr.makeTitle(),
		pr.makeBranchLine(),
		pr.makeStatus(),
	}, "\n")
}

func (pr PullRequest) Build() ui.Itemable {
	submenu := PullRequestOptions{
		Mergeable:       pr.Mergeable,
		ReviewDecision:  pr.ReviewDecision,
		RepoBranchRules: pr.RepoBranchRules,
		Permalink:       pr.Permalink,
	}.Build()

	return ui.SystraySubmenu{
		Title:   pr.makeFullTitle(),
		Submenu: &submenu,
	}
}
