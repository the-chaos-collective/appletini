package pages

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"appletini/config"
	"appletini/gitter"
	"appletini/logging"
	"appletini/ui"
	"appletini/ui/components"
	"appletini/ui/icons"

	"fyne.io/fyne/v2"
)

type IndexPage struct {
	systray      *ui.Systray
	Darkmode     bool
	PullRequests <-chan map[string][]gitter.PullRequest
	Trackers     config.Tracking
	Logger       logging.Logger
}

type iconState struct {
	green bool
	red   bool
}

var lightIcons = map[iconState]fyne.Resource{
	{false, false}: icons.ResIconDefault,
	{false, true}:  icons.ResIconReviewable,
	{true, false}:  icons.ResIconMergeable,
	{true, true}:   icons.ResIconBoth,
}

var darkIcons = map[iconState]fyne.Resource{
	{false, false}: icons.ResIconDefaultDark,
	{false, true}:  icons.ResIconReviewableDark,
	{true, false}:  icons.ResIconMergeableDark,
	{true, true}:   icons.ResIconBothDark,
}

func (page IndexPage) makeTree(prs map[string][]gitter.PullRequest) []ui.Itemable {
	result := make([]ui.Itemable, 0, 5) // separator + quit button + 3 tracking types by default

	for key, value := range prs {
		prList := make([]ui.Itemable, 0, 1) // at least one pr

		showGreenIcon := false
		showRedIcon := false

		for _, pr := range value {
			status := pr.PRInfo().Classify()

			prList = append(prList, components.PullRequest{
				Title:       pr.Title,
				Number:      pr.Number,
				HeadRefName: pr.HeadRefName,
				BaseRefName: pr.BaseRefName,
				Permalink:   pr.Permalink,
				Status:      status,
			}.Build())

			// TODO: Perhaps don't look at all tracked PRs
			// Only look at personal maybe?
			// Or anything related to the user

			if status.ShowGreenIcon {
				showGreenIcon = true
			}

			if status.ShowRedIcon {
				showRedIcon = true
			}
		}

		iconState := iconState{
			green: showGreenIcon,
			red:   showRedIcon,
		}

		if page.Darkmode {
			page.systray.SetIcon(darkIcons[iconState])
		} else {
			page.systray.SetIcon(lightIcons[iconState])
		}

		groupTitle := ""

		if key == "personal" {
			groupTitle = "My Pull Requests"
		} else {
			trackerType := strings.Split(key, "_")[0]

			idx, err := strconv.Atoi(strings.Split(key, "_")[1])
			if err != nil {
				page.Logger.Fatalf("unrecognized query response - malformed key (expected format [type]_[number]): %v", key)
			}

			switch trackerType {
			case "repo":
				groupTitle = page.Trackers.ByRepo[idx].Title
			case "label":
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
			Action: func() error {
				os.Exit(0)
				return nil
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
		icon = icons.ResIconDefaultDark
	} else {
		icon = icons.ResIconDefault
	}

	systray := ui.MakeSystray("Appletini", icon, page.Logger)

	systray.Setup()

	page.systray = &systray

	go page.run()

	page.systray.Run()
}
