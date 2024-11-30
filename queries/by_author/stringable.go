package by_author

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
		Authors: %v
	}`,
		tracker.Id,
		tracker.Owner,
		tracker.Repo,
		tracker.Title,
		tracker.Authors)
}
