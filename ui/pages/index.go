package pages

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"appletini/config"
	"appletini/gitter"
	"appletini/logging"
	"appletini/ui"
	"appletini/ui/components"
	"appletini/ui/icons"

	"fyne.io/fyne/v2"
)

type IndexPage struct {
	systray       *ui.Systray
	Darkmode      bool
	PullRequests  <-chan map[string][]gitter.PullRequest
	HasErr        <-chan bool
	Trackers      config.Tracking
	Logger        logging.Logger
	showGreenIcon bool
	showRedIcon   bool
}

type iconState struct {
	green bool
	red   bool
	err   bool
}

// TODO: If this becomes unwieldy consider composing/generating these icon combinations during compile-time

var lightIcons = map[iconState]fyne.Resource{
	{false, false, false}: icons.ResIconDefault,
	{false, true, false}:  icons.ResIconReviewable,
	{true, false, false}:  icons.ResIconMergeable,
	{true, true, false}:   icons.ResIconBoth,
	// TODO: Icons with error
	{false, false, true}: icons.ResIconWarning,
	{false, true, true}:  icons.ResIconWarning,
	{true, false, true}:  icons.ResIconWarning,
	{true, true, true}:   icons.ResIconWarning,
}

var darkIcons = map[iconState]fyne.Resource{
	{false, false, false}: icons.ResIconDefaultDark,
	{false, true, false}:  icons.ResIconReviewableDark,
	{true, false, false}:  icons.ResIconMergeableDark,
	{true, true, false}:   icons.ResIconBothDark,
	// TODO: Icons with error
	{false, false, true}: icons.ResIconWarning,
	{false, true, true}:  icons.ResIconWarning,
	{true, false, true}:  icons.ResIconWarning,
	{true, true, true}:   icons.ResIconWarning,
}

func (page IndexPage) renderIcons() {
	for {
		iconState := iconState{
			green: page.showGreenIcon,
			red:   page.showRedIcon,
			err:   <-page.HasErr,
		}

		if page.Darkmode {
			page.systray.SetIcon(darkIcons[iconState])
		} else {
			page.systray.SetIcon(lightIcons[iconState])
		}
	}
}

func (page *IndexPage) makeTree(prs map[string][]gitter.PullRequest) []ui.Itemable {
	result := make([]ui.Itemable, 0, 6) // separator + quit button + 4 tracking types by default

	page.showGreenIcon = false
	page.showRedIcon = false

	for key, value := range prs {
		prList := make([]ui.Itemable, 0, 1) // at least one pr

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

			if key == "personal" {
				if status.ShowGreenIcon {
					page.showGreenIcon = true
				}

				if status.ShowRedIcon {
					page.showRedIcon = true
				}
			}
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

func (page IndexPage) renderMenu() {
	for {
		page.systray.MainMenu.Items = page.makeTree(<-page.PullRequests)
		page.systray.Sync()
	}
}

func (page IndexPage) render() {
	go page.renderMenu()
	go page.renderIcons()

	for {
		time.Sleep(1 * time.Microsecond)
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

	go page.render()

	page.systray.Run()
}
