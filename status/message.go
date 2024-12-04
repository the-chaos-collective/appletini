package status

import "fmt"

func Message(input PRInfo) string {
	return fmt.Sprintf("%s | %s", input.Review, input.Mergeable)
}
