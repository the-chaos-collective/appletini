package pages

import (
	"os"

	"git_applet/gitter"
	"git_applet/ui"
	"git_applet/ui/components"
	"git_applet/ui/icons"
)

type IndexPage struct {
	systray      *ui.Systray
	PullRequests <-chan map[string][]gitter.PullRequest
}

func (page IndexPage) makeTree(prs map[string][]gitter.PullRequest) []ui.Itemable {
	prPersonalItems := make([]ui.Itemable, 0)
	prLabeledItems := make([]ui.Itemable, 0)
	prRepoItems := make([]ui.Itemable, 0)

	for _, pr := range prs["personal"] {
		prPersonalItems = append(prPersonalItems, components.PullRequest{
			Title:          pr.Title,
			Number:         pr.Number,
			Mergeable:      pr.Mergeable,
			ReviewDecision: pr.ReviewDecision,
			HeadRefName:    pr.HeadRefName,
			BaseRefName:    pr.BaseRefName,
			Permalink:      pr.Permalink,
		}.Build())
	}

	for _, pr := range prs["m20"] {
		prLabeledItems = append(prLabeledItems, components.PullRequest{
			Title:          pr.Title,
			Number:         pr.Number,
			Mergeable:      pr.Mergeable,
			ReviewDecision: pr.ReviewDecision,
			HeadRefName:    pr.HeadRefName,
			BaseRefName:    pr.BaseRefName,
			Permalink:      pr.Permalink,
		}.Build())
	}
	for _, pr := range prs["m21"] {
		prRepoItems = append(prRepoItems, components.PullRequest{
			Title:          pr.Title,
			Number:         pr.Number,
			Mergeable:      pr.Mergeable,
			ReviewDecision: pr.ReviewDecision,
			HeadRefName:    pr.HeadRefName,
			BaseRefName:    pr.BaseRefName,
			Permalink:      pr.Permalink,
		}.Build())
	}

	return []ui.Itemable{
		ui.SystraySubmenu{
			Title: "My Pull Requests",
			Submenu: &ui.SystrayMenu{
				Items: prPersonalItems,
			},
		},
		ui.SystraySubmenu{
			Title: "Labeled Pull Requests",
			Submenu: &ui.SystrayMenu{
				Items: prLabeledItems,
			},
		},
		ui.SystraySubmenu{
			Title: "Repo Pull Requests",
			Submenu: &ui.SystrayMenu{
				Items: prRepoItems,
			},
		},
		ui.SystraySeparator{},
		ui.SystrayButton{
			Title: "Quit",
			Action: func() {
				os.Exit(0)
			},
		},
	}
}

func (page IndexPage) run() {
	for {
		page.systray.MainMenu.Items = page.makeTree(<-page.PullRequests)
		page.systray.Sync()
	}
}

func (page IndexPage) Run() {
	systray := ui.MakeSystray("Git Appletini", icons.ResIconDefault)

	systray.Setup()

	page.systray = &systray

	go page.run()

	page.systray.Run()
}
