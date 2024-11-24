package repo

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
		config.CommentsAmount,
	)
}

func (tracker Tracker) String() string {
	return fmt.Sprintf(`{
		Owner: %v,
		Name: %v,
		Identifier: %v
	}`,
		tracker.Owner,
		tracker.Name,
		tracker.Identifier)
}
