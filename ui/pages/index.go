package pages

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"git_applet/gitter"
	"git_applet/types"
	"git_applet/ui"
	"git_applet/ui/components"
	"git_applet/ui/icons"

	"fyne.io/fyne/v2"
)

type IndexPage struct {
	systray      *ui.Systray
	Darkmode     bool
	PullRequests <-chan map[string][]gitter.PullRequest
	Trackers     types.Tracking
}

func (page IndexPage) makeTree(prs map[string][]gitter.PullRequest) []ui.Itemable {
	result := make([]ui.Itemable, 0, 5) // separator + quit button + 3 tracking types by default
	for key, value := range prs {
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

		groupTitle := ""

		if key == "personal" {
			groupTitle = "My Pull Requests"
		} else {
			trackerType := strings.Split(key, "_")[0]

			idx, err := strconv.Atoi(strings.Split(key, "_")[1])
			if err != nil {
				log.Fatalf("unrecognized query response - malformed key (expected format [type]_[number]): %v", key)
			}

			switch trackerType {
			case "repo":
				groupTitle = page.Trackers.ByRepo[idx].Title
			case "labeled":
				groupTitle = page.Trackers.ByLabel[idx].Title
			case "author":
				groupTitle = page.Trackers.ByAuthor[idx].Title
			default:
				groupTitle = fmt.Sprintf("Unsupported tracker type: %s", trackerType)
			}
		}

		tmp := ui.SystraySubmenu{
			Title: groupTitle,
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
	var icon fyne.Resource
	if page.Darkmode {
		icon = icons.ResIconDefault
	} else {
		icon = icons.ResIconDefaultDark
	}
	systray := ui.MakeSystray("Git Appletini", icon)

	systray.Setup()

	page.systray = &systray

	go page.run()

	page.systray.Run()
}
