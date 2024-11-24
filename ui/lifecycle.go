package ui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
)

func (systray *Systray) createApp() {
	// create app
	app := app.NewWithID("git_appletini")

	systray.fyneApp = &app

	var ok bool

	// get desktop app
	desk, ok := (*systray.fyneApp).(desktop.App)
	if !ok {
		log.Fatal("could not create desktop app")
	}

	// create main menu
	systray.MainMenu.fyneMenu = fyne.NewMenu(systray.title)

	// set main menu onto the desktop app created
	desk.SetSystemTrayMenu(systray.MainMenu.fyneMenu)

	(*systray.fyneApp).Lifecycle().SetOnStarted(func() {

		// set default icon on main menu
		desk.SetSystemTrayIcon(systray.icon)

	})
}

func (systray *Systray) SetIcon(icon fyne.Resource) {
	systray.icon = icon

	desk, ok := (*systray.fyneApp).(desktop.App)
	if !ok {
		log.Fatal("could not get desktop app")
	}

	desk.SetSystemTrayIcon(systray.icon)
}

func (systray *Systray) Setup() {
	systray.createApp()
}

func (systray *Systray) Sync() {
	systray.MainMenu.sync()
}

func (systray Systray) Run() {
	if systray.fyneApp == nil {
		log.Fatal("systray not setup yet")
	}

	(*systray.fyneApp).Run()
}

// TODO: Only update if something changed (hashing)
func (menu *SystrayMenu) sync() {
	menu.fyneMenu.Items = []*fyne.MenuItem{}

	for _, item := range menu.Items {
		switch (item).Type() {
		case Button:
			item := (item).(SystrayButton)

			menuItem := fyne.NewMenuItem(item.Title, item.Action)
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
			item.Submenu.sync()

		}
	}

	menu.fyneMenu.Refresh()
}
