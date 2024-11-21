package labeled

import (
	"fmt"
)

func (config Config) String() string {
	return fmt.Sprintf(`{
	Trackers: %v,
	PrAmount: %v,
	ReviewAmount: %v,
	CommentsAmount: %v
	}`,
		config.Trackers,
		config.PrAmount,
		config.ReviewAmount,
		config.CommentsAmount)
}

func (tracker Tracker) String() string {
	return fmt.Sprintf(`{
		Owner: %v,
		Repo: %v,
		Identifier: %v,
		Label: %v
	}`,
		tracker.Owner,
		tracker.Repo,
		tracker.Identifier,
		tracker.Label)
}
