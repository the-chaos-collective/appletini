package ui

import (
	"errors"
	"fmt"

	"appletini/logging"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

var (
	ErrNotDesktop = errors.New("fyne app is not a desktop app (no system tray support)")
	ErrNotSetup   = errors.New("systray not set up yet")
	ErrNilItem    = errors.New("menu item must not be nil")
	ErrNilSubmenu = errors.New("submenu must not be nil")
	ErrBadItem    = errors.New("unsupported menu item")
	ErrNilAction  = errors.New("confirm action must not be nil")
)

func (s *Systray) createApp() error {
	fyneApp := app.NewWithID("appletini")

	desk, ok := fyneApp.(desktop.App)
	if !ok {
		return ErrNotDesktop
	}

	s.fyneApp = fyneApp

	// create main menu and set it onto the desktop app
	s.MainMenu.fyneMenu = fyne.NewMenu(s.title)
	desk.SetSystemTrayMenu(s.MainMenu.fyneMenu)

	// set default icon on main menu
	desk.SetSystemTrayIcon(s.icon)
	fyneApp.Lifecycle().SetOnStarted(func() {
		s.mu.Lock()
		defer s.mu.Unlock()

		desk.SetSystemTrayIcon(s.icon)
	})

	return nil
}

func (s *Systray) SetIcon(icon fyne.Resource) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.fyneApp == nil {
		return ErrNotSetup
	}

	desk, ok := s.fyneApp.(desktop.App)
	if !ok {
		return ErrNotDesktop
	}

	s.icon = icon
	desk.SetSystemTrayIcon(s.icon)

	return nil
}

func (s *Systray) Setup() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.createApp()
}

// SetMenu replaces the main menu's items and redraws it.
func (s *Systray) SetMenu(items []Itemable) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.MainMenu.fyneMenu == nil {
		return ErrNotSetup
	}

	s.MainMenu.Items = items

	return s.MainMenu.sync(s.Logger)
}

// Notify shows a desktop notification.
func (s *Systray) Notify(title, content string) error {
	s.mu.Lock()
	fyneApp := s.fyneApp
	s.mu.Unlock()

	if fyneApp == nil {
		return ErrNotSetup
	}

	fyneApp.SendNotification(fyne.NewNotification(title, content))

	return nil
}

// Confirm opens a small window asking the user to confirm an action. onConfirm
// runs only if they press the confirm button; cancelling or closing the window
// does nothing.
func (s *Systray) Confirm(title, message, confirmLabel string, onConfirm func()) error {
	s.mu.Lock()
	fyneApp := s.fyneApp
	s.mu.Unlock()

	if fyneApp == nil {
		return ErrNotSetup
	}

	if onConfirm == nil {
		return ErrNilAction
	}

	win := fyneApp.NewWindow(title)
	win.Resize(fyne.NewSize(380, 180))
	win.CenterOnScreen()

	confirm := dialog.NewConfirm(title, message, func(ok bool) {
		win.Close()

		if ok {
			onConfirm()
		}
	}, win)
	confirm.SetDismissText("Cancel")
	confirm.SetConfirmText(confirmLabel)
	confirm.SetConfirmImportance(widget.DangerImportance)

	win.Show()
	confirm.Show()

	return nil
}

// Run blocks until the app quits.
func (s *Systray) Run() error {
	s.mu.Lock()
	fyneApp := s.fyneApp
	s.mu.Unlock()

	if fyneApp == nil {
		return ErrNotSetup
	}

	fyneApp.Run()

	return nil
}

// Quit stops the app, which makes Run return. It is safe to call before
// Setup and more than once; only the first call after Setup has an effect.
func (s *Systray) Quit() {
	s.mu.Lock()
	fyneApp := s.fyneApp
	s.mu.Unlock()

	if fyneApp == nil {
		return
	}

	s.quitOnce.Do(fyneApp.Quit)
}

// buildItems converts menu items into fyne menu items, recursing into submenus.
func (menu *SystrayMenu) buildItems(logger logging.Logger) ([]*fyne.MenuItem, error) {
	items := make([]*fyne.MenuItem, 0, len(menu.Items))

	for i, item := range menu.Items {
		if item == nil {
			return nil, fmt.Errorf("item %d: %w", i, ErrNilItem)
		}

		switch item.Type() {
		case Button:
			button, ok := item.(SystrayButton)
			if !ok {
				return nil, fmt.Errorf("item %d: %w: %T claims to be a button", i, ErrBadItem, item)
			}

			items = append(items, fyne.NewMenuItem(button.Title, func() {
				if button.Action == nil {
					return
				}

				// A failed action (e.g. no browser) must not take the app down.
				err := button.Action()
				if err != nil {
					logger.Error("menu action failed", "item", button.Title, "err", err)
				}
			}))

		case Separator:
			items = append(items, fyne.NewMenuItemSeparator())

		case Submenu:
			sub, ok := item.(SystraySubmenu)
			if !ok {
				return nil, fmt.Errorf("item %d: %w: %T claims to be a submenu", i, ErrBadItem, item)
			}

			if sub.Submenu == nil {
				return nil, fmt.Errorf("item %d (%s): %w", i, sub.Title, ErrNilSubmenu)
			}

			childItems, err := sub.Submenu.buildItems(logger)
			if err != nil {
				return nil, fmt.Errorf("submenu %s: %w", sub.Title, err)
			}

			menuItem := fyne.NewMenuItem(sub.Title, func() {})
			menuItem.ChildMenu = fyne.NewMenu(sub.Title, childItems...)
			sub.Submenu.fyneMenu = menuItem.ChildMenu
			items = append(items, menuItem)

		default:
			return nil, fmt.Errorf("item %d: %w: type %d", i, ErrBadItem, item.Type())
		}
	}

	return items, nil
}

func (menu *SystrayMenu) sync(logger logging.Logger) error {
	if menu.fyneMenu == nil {
		return ErrNotSetup
	}

	items, err := menu.buildItems(logger)
	if err != nil {
		return err
	}

	menu.fyneMenu.Items = items
	menu.fyneMenu.Refresh()

	return nil
}
