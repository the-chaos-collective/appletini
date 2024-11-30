package ui

import (
	"log"

	"fyne.io/fyne/v2"
)

func MakeSystray(title string, icon fyne.Resource, logger *log.Logger) Systray {
	return Systray{
		title: title,
		icon:  icon,
		MainMenu: &SystrayMenu{
			Items: make([]Itemable, 0),
		},
		Logger: logger,
	}
}
