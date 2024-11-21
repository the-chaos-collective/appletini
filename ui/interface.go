package ui

import "fyne.io/fyne/v2"

type Itemable interface {
	Type() ItemType
}

type Component interface {
	Build() *Itemable
}

type Lifecycle interface {
	Setup()
	Run()
}

const (
	Separator = iota
	Button
	Submenu
)

type ItemType int

type Systray struct {
	fyneApp  *fyne.App
	title    string
	icon     fyne.Resource
	MainMenu *SystrayMenu
}

type SystrayMenu struct {
	fyneMenu *fyne.Menu
	Items    []Itemable
}

type SystrayButton struct {
	Title  string
	Action func()
}

func (SystrayButton) Type() ItemType {
	return Button
}

type SystraySubmenu struct {
	Title   string
	Submenu *SystrayMenu
}

func (SystraySubmenu) Type() ItemType {
	return Submenu
}

type SystraySeparator struct{}

func (SystraySeparator) Type() ItemType {
	return Separator
}
