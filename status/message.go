package status

import "fmt"

func Message(input PRInfo) string {
	if input.Checks == ChecksState_None {
		return fmt.Sprintf("%s | %s", input.Review, input.Mergeable)
	}

	return fmt.Sprintf("%s | %s | %s", input.Review, input.Mergeable, input.Checks)
}
