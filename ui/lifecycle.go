package ui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
)

func (s *Systray) createApp() {
	// create app
	app := app.NewWithID("appletini")

	s.fyneApp = &app

	var ok bool

	// get desktop app
	desk, ok := (*s.fyneApp).(desktop.App)
	if !ok {
		s.Logger.Fatal("could not create desktop app")
	}

	// create main menu
	s.MainMenu.fyneMenu = fyne.NewMenu(s.title)

	// set main menu onto the desktop app created
	desk.SetSystemTrayMenu(s.MainMenu.fyneMenu)

	(*s.fyneApp).Lifecycle().SetOnStarted(func() {
		// set default icon on main menu
		desk.SetSystemTrayIcon(s.icon)
	})
}

func (s *Systray) SetIcon(icon fyne.Resource) {
	s.icon = icon

	desk, ok := (*s.fyneApp).(desktop.App)
	if !ok {
		s.Logger.Fatal("could not get desktop app")
	}

	desk.SetSystemTrayIcon(s.icon)
}

func (s *Systray) Setup() {
	s.createApp()
}

func (s *Systray) Sync() {
	s.MainMenu.sync(s.Logger)
}

func (s Systray) Run() {
	if s.fyneApp == nil {
		s.Logger.Fatal("systray not setup yet")
	}

	(*s.fyneApp).Run()
}

func (menu *SystrayMenu) sync(logger *log.Logger) {
	menu.fyneMenu.Items = []*fyne.MenuItem{}

	for _, item := range menu.Items {
		switch (item).Type() {
		case Button:
			item := (item).(SystrayButton)

			menuItem := fyne.NewMenuItem(item.Title, func() {
				err := item.Action()
				if err != nil {
					logger.Fatal(err)
				}
			})
			menu.fyneMenu.Items = append(menu.fyneMenu.Items, menuItem)

		case Separator:
			menuItem := fyne.NewMenuItemSeparator()
			menu.fyneMenu.Items = append(menu.fyneMenu.Items, menuItem)

		case Submenu:
			item := (item).(SystraySubmenu)

			menuItem := fyne.NewMenuItem(item.Title, func() {})
			menuItem.ChildMenu = fyne.NewMenu(item.Title)
			menu.fyneMenu.Items = append(menu.fyneMenu.Items, menuItem)

			item.Submenu.fyneMenu = menuItem.ChildMenu
			item.Submenu.sync(logger)

		}
	}

	menu.fyneMenu.Refresh()
}
