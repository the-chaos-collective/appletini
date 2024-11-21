package actions

import (
	"fmt"
	"os/exec"
)

func OpenLink(url string) error {
	a := exec.Command("open", url)

	err := a.Start()
	if err != nil {
		return fmt.Errorf("error opening URL (%v): %w", url, err)
	}

	return nil
}
