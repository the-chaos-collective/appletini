package actions

import (
	"errors"
	"slices"
	"testing"
)

func TestOpenCommand(t *testing.T) {
	const link = "https://github.com/o/r/pull/1"

	testCases := []struct {
		name         string
		goos         string
		link         string
		expectedName string
		expectedArgs []string
		expectedErr  error
		expectAnyErr bool
	}{
		{name: "macOS", goos: "darwin", link: link, expectedName: "open", expectedArgs: []string{link}},
		{name: "linux", goos: "linux", link: link, expectedName: "xdg-open", expectedArgs: []string{link}},
		{name: "freebsd", goos: "freebsd", link: link, expectedName: "xdg-open", expectedArgs: []string{link}},
		{name: "windows", goos: "windows", link: link, expectedName: "rundll32", expectedArgs: []string{"url.dll,FileProtocolHandler", link}},
		{name: "http is allowed", goos: "linux", link: "http://x", expectedName: "xdg-open", expectedArgs: []string{"http://x"}},
		{name: "unsupported os", goos: "plan9", link: link, expectedErr: ErrUnsupportedOS},
		{name: "empty url", goos: "linux", link: "", expectedErr: ErrEmptyURL},
		{name: "file url", goos: "linux", link: "file:///etc/passwd", expectedErr: ErrUnsupportedScheme},
		{name: "flag-like value", goos: "linux", link: "--help", expectedErr: ErrUnsupportedScheme},
		{name: "unparseable", goos: "linux", link: "http://[::1", expectAnyErr: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			name, args, err := openCommand(tc.goos, tc.link)

			if tc.expectAnyErr {
				if err == nil {
					t.Fatal("expected an error")
				}

				return
			}

			if !errors.Is(err, tc.expectedErr) {
				t.Fatalf("got error %v, expected %v", err, tc.expectedErr)
			}

			if name != tc.expectedName || !slices.Equal(args, tc.expectedArgs) {
				t.Fatalf("got %s %v, expected %s %v", name, args, tc.expectedName, tc.expectedArgs)
			}
		})
	}
}

func TestOpenLinkValidatesBeforeRunning(t *testing.T) {
	err := OpenLink("javascript:alert(1)")
	if !errors.Is(err, ErrUnsupportedScheme) {
		t.Fatalf("got error %v, expected %v", err, ErrUnsupportedScheme)
	}
}
