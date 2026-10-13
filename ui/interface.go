package ui

import (
	"sync"

	"appletini/logging"

	"fyne.io/fyne/v2"
)

type Itemable interface {
	Type() ItemType
}

type Component interface {
	Build() *Itemable
}

type Lifecycle interface {
	Setup() error
	Run() error
	Quit()
}

const (
	Separator = iota
	Button
	Submenu
)

type ItemType int

// Systray wraps the fyne tray app. The menu and icon are updated from
// different goroutines, so mu guards both (a single mutex, so there is no
// lock ordering to get wrong).
type Systray struct {
	mu       sync.Mutex
	quitOnce sync.Once
	fyneApp  fyne.App
	title    string
	icon     fyne.Resource
	MainMenu *SystrayMenu
	Logger   logging.Logger
}

type SystrayMenu struct {
	fyneMenu *fyne.Menu
	Items    []Itemable
}

type SystrayButton struct {
	Title  string
	Action func() error
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
