package actions

import (
	"errors"
	"fmt"
	"net/url"
	"os/exec"
	"runtime"
)

var (
	ErrEmptyURL          = errors.New("url must not be empty")
	ErrUnsupportedScheme = errors.New("only http and https links can be opened")
	ErrUnsupportedOS     = errors.New("opening links is not supported on this OS")
)

// openCommand returns the command that opens link in the default browser on goos.
func openCommand(goos, link string) (string, []string, error) {
	if link == "" {
		return "", nil, ErrEmptyURL
	}

	parsed, err := url.Parse(link)
	if err != nil {
		return "", nil, fmt.Errorf("parsing URL %q: %w", link, err)
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return "", nil, fmt.Errorf("%w: %q", ErrUnsupportedScheme, link)
	}

	switch goos {
	case "darwin":
		return "open", []string{link}, nil
	case "windows":
		return "rundll32", []string{"url.dll,FileProtocolHandler", link}, nil
	case "linux", "freebsd", "openbsd", "netbsd":
		return "xdg-open", []string{link}, nil
	}

	return "", nil, fmt.Errorf("%w: %s", ErrUnsupportedOS, goos)
}

// OpenLink opens link in the user's default browser.
func OpenLink(link string) error {
	name, args, err := openCommand(runtime.GOOS, link)
	if err != nil {
		return err
	}

	//nolint:gosec // the command is fixed per OS and the URL is validated above
	cmd := exec.Command(name, args...)

	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("error opening URL (%v): %w", link, err)
	}

	// Reap the child so it doesn't linger as a zombie; its exit status is irrelevant.
	go func() {
		_ = cmd.Wait()
	}()

	return nil
}
