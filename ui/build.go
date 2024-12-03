package ui

import (
	"log"

	"fyne.io/fyne/v2"
)

func MakeSystray(title string, icon fyne.Resource, logger *log.Logger) Systray {
	return Systray{
		fyneApp: nil,
		title:   title,
		icon:    icon,
		MainMenu: &SystrayMenu{
			fyneMenu: nil,
			Items:    make([]Itemable, 0),
		},
		Logger: logger,
	}
}
