package pages

import (
	"fmt"
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

	result := make([]ui.Itemable, 0, 5) // separator + quit button + 3 tracking types by default
	for key, value := range prs {
		fmt.Println(key)
		prList := make([]ui.Itemable, 0, 1) // at least one pr
		for _, pr := range value {
			prList = append(prList, components.PullRequest{
				Title:          pr.Title,
				Number:         pr.Number,
				Mergeable:      pr.Mergeable,
				ReviewDecision: pr.ReviewDecision,
				HeadRefName:    pr.HeadRefName,
				BaseRefName:    pr.BaseRefName,
				Permalink:      pr.Permalink,
			}.Build())
		}
		tmp := ui.SystraySubmenu{
			Title: key,
			Submenu: &ui.SystrayMenu{
				Items: prList,
			},
		}
		result = append(result, tmp)
	}
	finalItems := []ui.Itemable{
		ui.SystraySeparator{},
		ui.SystrayButton{
			Title: "Quit",
			Action: func() {
				os.Exit(0)
			},
		},
	}
	result = append(result, finalItems...)

	return result
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
