package ui

import "fyne.io/fyne/v2"

func MakeSystray(title string, icon fyne.Resource) Systray {
	return Systray{
		title: title,
		icon:  icon,
		MainMenu: &SystrayMenu{
			Items: make([]Itemable, 0),
		},
	}
}
