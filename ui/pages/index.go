package pages

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"appletini/config"
	"appletini/gitter"
	"appletini/logging"
	"appletini/reviews"
	"appletini/trackerkey"
	"appletini/ui"
	"appletini/ui/components"
	"appletini/ui/icons"

	"fyne.io/fyne/v2"
)

var (
	ErrTrackerOutOfRange = errors.New("tracker index out of range for current config")
	ErrUnsupportedKind   = errors.New("unsupported tracker kind")
	ErrMissingIcon       = errors.New("no icon for state")
)

// Approver approves pull requests (reviews.Approver in production).
type Approver interface {
	Approve(ctx context.Context, prID, comment string) (reviews.Review, error)
}

type IndexPage struct {
	systray      *ui.Systray
	ctx          context.Context
	Darkmode     bool
	PullRequests <-chan map[string][]gitter.PullRequest
	HasErr       <-chan bool
	Trackers     config.Tracking
	Logger       logging.Logger
	// Approver enables the Approve option on other people's PRs; nil hides it.
	Approver Approver
}

// treeActions are the callbacks the menu items run.
type treeActions struct {
	quit func() error
	// approve returns the Approve action for a PR, or nil to hide the option.
	approve func(pr gitter.PullRequest) func() error
}

// iconState is what the tray icon reflects. green/red come from the user's
// own PRs in the latest results; err comes only from the poller, so a menu
// rebuild can no longer clear the warning icon.
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

func iconFor(state iconState, darkmode bool) (fyne.Resource, error) {
	set := lightIcons
	if darkmode {
		set = darkIcons
	}

	icon, ok := set[state]
	if !ok {
		return nil, fmt.Errorf("%w: %+v", ErrMissingIcon, state)
	}

	return icon, nil
}

// groupTitle is the submenu title for a tracker's results.
func groupTitle(key trackerkey.Key, trackers config.Tracking) (string, error) {
	titleAt := func(n int, title func(int) string) (string, error) {
		if key.Index >= n {
			return "", fmt.Errorf("%w: %s_%d (have %d)", ErrTrackerOutOfRange, key.Kind, key.Index, n)
		}

		return title(key.Index), nil
	}

	switch key.Kind {
	case trackerkey.Personal:
		return "My Pull Requests", nil
	case trackerkey.Repo:
		return titleAt(len(trackers.ByRepo), func(i int) string { return trackers.ByRepo[i].Title })
	case trackerkey.Label:
		return titleAt(len(trackers.ByLabel), func(i int) string { return trackers.ByLabel[i].Title })
	case trackerkey.Author:
		return titleAt(len(trackers.ByAuthor), func(i int) string { return trackers.ByAuthor[i].Title })
	}

	return "", fmt.Errorf("%w: %q", ErrUnsupportedKind, key.Kind)
}

// sortedKeys parses the result keys and orders them for display, so the
// menu doesn't reshuffle between rebuilds. Unparseable keys are returned as
// a joined error and left out.
func sortedKeys(prs map[string][]gitter.PullRequest) ([]trackerkey.Key, map[trackerkey.Key]string, error) {
	keys := make([]trackerkey.Key, 0, len(prs))
	raw := make(map[trackerkey.Key]string, len(prs))

	var errs []error

	for rawKey := range prs {
		key, err := trackerkey.Parse(rawKey)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		keys = append(keys, key)
		raw[key] = rawKey
	}

	slices.SortFunc(keys, trackerkey.Compare)

	return keys, raw, errors.Join(errs...)
}

func quitItems(quit func() error) []ui.Itemable {
	return []ui.Itemable{
		ui.SystraySeparator{},
		ui.SystrayButton{
			Title:  "Quit",
			Action: quit,
		},
	}
}

// buildTree turns the latest results into menu items and works out which
// icons the user's own PRs call for. Groups that can't be titled are skipped
// and reported in the returned error; the rest of the menu is still built.
func buildTree(
	prs map[string][]gitter.PullRequest,
	trackers config.Tracking,
	acts treeActions,
) ([]ui.Itemable, iconState, error) {
	//exhaustruct:ignore
	state := iconState{}
	result := make([]ui.Itemable, 0, len(prs)+2) // groups + separator + quit button

	keys, rawKeys, keyErr := sortedKeys(prs)
	errs := []error{keyErr}

	for _, key := range keys {
		title, err := groupTitle(key, trackers)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		prList := make([]ui.Itemable, 0, len(prs[rawKeys[key]]))

		for _, pr := range prs[rawKeys[key]] {
			status := pr.PRInfo().Classify()

			// No Approve on your own PRs (GitHub doesn't allow it); otherwise
			// see ApprovableByViewer (already approved, or red for a reason
			// other than your own change request).
			var approve func() error
			if key.Kind != trackerkey.Personal && pr.ApprovableByViewer() && acts.approve != nil {
				approve = acts.approve(pr)
			}

			prList = append(prList, components.PullRequest{
				Title:       pr.Title,
				Number:      pr.Number,
				HeadRefName: pr.HeadRefName,
				BaseRefName: pr.BaseRefName,
				Permalink:   pr.Permalink,
				Status:      status,
				Approve:     approve,
			}.Build())

			if key.Kind == trackerkey.Personal {
				state.green = state.green || status.ShowGreenIcon
				state.red = state.red || status.ShowRedIcon
			}
		}

		result = append(result, ui.SystraySubmenu{
			Title: title,
			Submenu: &ui.SystrayMenu{
				Items: prList,
			},
		})
	}

	result = append(result, quitItems(acts.quit)...)

	return result, state, errors.Join(errs...)
}

func (page *IndexPage) quit() error {
	page.systray.Quit()

	return nil
}

func (page *IndexPage) actions() treeActions {
	//exhaustruct:ignore
	acts := treeActions{quit: page.quit}

	if page.Approver != nil {
		acts.approve = func(pr gitter.PullRequest) func() error {
			return func() error {
				// Menu actions run on the UI thread: don't block it on the network.
				go func() {
					err := approvePR(page.ctx, page.Approver, page.systray.Notify, pr)
					if err != nil {
						page.Logger.Error("approving PR", "number", pr.Number, "err", err)
					}
				}()

				return nil
			}
		}
	}

	return acts
}

// approvalComment is posted with every approval made from the tray.
const approvalComment = "LGTM!"

// approvePR approves pr and reports the outcome as a desktop notification.
func approvePR(
	ctx context.Context,
	approver Approver,
	notify func(title, content string) error,
	pr gitter.PullRequest,
) error {
	if approver == nil {
		return errors.New("no approver configured")
	}

	_, err := approver.Approve(ctx, pr.Id, approvalComment)
	if err != nil {
		notifyErr := notify(fmt.Sprintf("Couldn't approve #%d", pr.Number), err.Error())

		return errors.Join(err, notifyErr)
	}

	return notify(fmt.Sprintf("Approved #%d", pr.Number), pr.Title)
}

func (page *IndexPage) applyIcon(state iconState) {
	icon, err := iconFor(state, page.Darkmode)
	if err != nil {
		page.Logger.Error("choosing tray icon", "err", err)
		return
	}

	err = page.systray.SetIcon(icon)
	if err != nil {
		page.Logger.Error("setting tray icon", "err", err)
	}
}

// loop applies poller updates to the menu and icon until ctx is cancelled.
// It is the only goroutine that touches the page's state.
func (page *IndexPage) loop(ctx context.Context) {
	//exhaustruct:ignore
	state := iconState{}

	for {
		select {
		case <-ctx.Done():
			return

		case prs := <-page.PullRequests:
			items, prState, err := buildTree(prs, page.Trackers, page.actions())
			if err != nil {
				page.Logger.Error("building menu", "err", err)
			}

			err = page.systray.SetMenu(items)
			if err != nil {
				page.Logger.Error("updating menu", "err", err)
			}

			state.green, state.red = prState.green, prState.red

		case hasErr := <-page.HasErr:
			state.err = hasErr
		}

		page.applyIcon(state)
	}
}

// Run shows the tray and blocks until the app quits, either from the Quit
// button or because ctx was cancelled.
func (page *IndexPage) Run(ctx context.Context) error {
	//exhaustruct:ignore
	initialIcon, err := iconFor(iconState{}, page.Darkmode)
	if err != nil {
		return err
	}

	page.systray = ui.MakeSystray("Appletini", initialIcon, page.Logger)

	err = page.systray.Setup()
	if err != nil {
		return fmt.Errorf("setting up systray: %w", err)
	}

	// Usable (Quit works) before the first poll comes back.
	err = page.systray.SetMenu(quitItems(page.quit))
	if err != nil {
		return fmt.Errorf("setting initial menu: %w", err)
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Approvals started from the menu are cancelled when the app quits.
	page.ctx = ctx

	go page.loop(ctx)

	go func() {
		<-ctx.Done()
		page.systray.Quit()
	}()

	return page.systray.Run()
}
