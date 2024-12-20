package by_label

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
		ID: %v,
		Owner: %v,
		Repo: %v,
		Identifier: %v,
		Label: %v
	}`,
		tracker.Id,
		tracker.Owner,
		tracker.Repo,
		tracker.Title,
		tracker.Label)
}
