package ui

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/charmbracelet/log"

	"appletini/logging"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
)

// badItem claims a type it doesn't have.
type badItem struct{ t ItemType }

func (b badItem) Type() ItemType { return b.t }

func bufferLogger() (logging.Logger, *bytes.Buffer) {
	var buf bytes.Buffer

	return logging.NewCharm(log.New(&buf)), &buf
}

func TestBuildItems(t *testing.T) {
	noop := func() error { return nil }

	testCases := []struct {
		name        string
		items       []Itemable
		expected    []string // labels, "-" for separators, "Title>[...]" for submenus
		expectedErr error
	}{
		{
			name:     "empty",
			items:    nil,
			expected: []string{},
		},
		{
			name: "buttons, separator and nested submenus",
			items: []Itemable{
				SystraySubmenu{Title: "Group", Submenu: &SystrayMenu{Items: []Itemable{
					SystraySubmenu{Title: "PR", Submenu: &SystrayMenu{Items: []Itemable{
						SystrayButton{Title: "Open", Action: noop},
					}}},
				}}},
				SystraySeparator{},
				SystrayButton{Title: "Quit", Action: noop},
			},
			expected: []string{"Group>[PR>[Open]]", "-", "Quit"},
		},
		{name: "nil item", items: []Itemable{nil}, expectedErr: ErrNilItem},
		{name: "nil submenu", items: []Itemable{SystraySubmenu{Title: "x"}}, expectedErr: ErrNilSubmenu},
		{
			name: "error inside a submenu",
			items: []Itemable{SystraySubmenu{Title: "x", Submenu: &SystrayMenu{Items: []Itemable{
				badItem{t: 42},
			}}}},
			expectedErr: ErrBadItem,
		},
		{name: "fake button", items: []Itemable{badItem{t: Button}}, expectedErr: ErrBadItem},
		{name: "fake submenu", items: []Itemable{badItem{t: Submenu}}, expectedErr: ErrBadItem},
		{name: "unknown type", items: []Itemable{badItem{t: 42}}, expectedErr: ErrBadItem},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			logger, _ := bufferLogger()
			menu := &SystrayMenu{Items: tc.items}

			got, err := menu.buildItems(logger)
			if !errors.Is(err, tc.expectedErr) {
				t.Fatalf("got error %v, expected %v", err, tc.expectedErr)
			}

			if tc.expectedErr != nil {
				return
			}

			labels := describe(got)
			if strings.Join(labels, ",") != strings.Join(tc.expected, ",") {
				t.Fatalf("got %v, expected %v", labels, tc.expected)
			}
		})
	}
}

func describe(items []*fyne.MenuItem) []string {
	out := make([]string, 0, len(items))

	for _, item := range items {
		switch {
		case item.IsSeparator:
			out = append(out, "-")
		case item.ChildMenu != nil:
			out = append(out, item.Label+">["+strings.Join(describe(item.ChildMenu.Items), ",")+"]")
		default:
			out = append(out, item.Label)
		}
	}

	return out
}

func TestButtonActionErrorIsLoggedNotFatal(t *testing.T) {
	logger, buf := bufferLogger()
	called := 0

	menu := &SystrayMenu{Items: []Itemable{
		SystrayButton{Title: "Open in browser", Action: func() error {
			called++
			return errors.New("no browser")
		}},
		SystrayButton{Title: "No action"},
	}}

	items, err := menu.buildItems(logger)
	if err != nil {
		t.Fatal(err)
	}

	items[0].Action() // would have exited the process before
	items[1].Action() // nil action is a no-op

	if called != 1 {
		t.Fatalf("action called %d times", called)
	}

	if !strings.Contains(buf.String(), "no browser") {
		t.Fatalf("error was not logged: %q", buf.String())
	}
}

func TestNotSetUp(t *testing.T) {
	logger, _ := bufferLogger()
	s := MakeSystray("t", nil, logger)

	testCases := []struct {
		name string
		call func() error
	}{
		{name: "SetMenu", call: func() error { return s.SetMenu(nil) }},
		{name: "SetIcon", call: func() error { return s.SetIcon(nil) }},
		{name: "Run", call: s.Run},
		{name: "Notify", call: func() error { return s.Notify("t", "c") }},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.call(); !errors.Is(err, ErrNotSetup) {
				t.Fatalf("got error %v, expected %v", err, ErrNotSetup)
			}
		})
	}

	s.Quit() // must not panic before setup
}

func TestSetMenu(t *testing.T) {
	test.NewApp()

	logger, _ := bufferLogger()
	s := MakeSystray("Appletini", nil, logger)
	s.MainMenu.fyneMenu = fyne.NewMenu("Appletini")

	testCases := []struct {
		name        string
		items       []Itemable
		expected    []string
		expectedErr error
	}{
		{
			name:     "sets items",
			items:    []Itemable{SystrayButton{Title: "A"}, SystraySeparator{}, SystrayButton{Title: "B"}},
			expected: []string{"A", "-", "B"},
		},
		{
			name:     "replaces previous items",
			items:    []Itemable{SystrayButton{Title: "C"}},
			expected: []string{"C"},
		},
		{
			name:        "invalid items leave the menu unchanged",
			items:       []Itemable{nil},
			expected:    []string{"C"},
			expectedErr: ErrNilItem,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := s.SetMenu(tc.items)
			if !errors.Is(err, tc.expectedErr) {
				t.Fatalf("got error %v, expected %v", err, tc.expectedErr)
			}

			got := describe(s.MainMenu.fyneMenu.Items)
			if strings.Join(got, ",") != strings.Join(tc.expected, ",") {
				t.Fatalf("got %v, expected %v", got, tc.expected)
			}
		})
	}
}
