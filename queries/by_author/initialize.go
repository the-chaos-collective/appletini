package by_author

import (
	"errors"
	"fmt"
)

func MakeQuery(config Config) (Query, error) {
	if len(config.Trackers) == 0 {
		return Query{
			shouldBeExecuted: false,
			generatedQuery:   "",
		}, nil
	}

	err := validateConfig(config)
	if err != nil {
		return Query{}, fmt.Errorf("invalid config: %w", err)
	}

	query, err := generateQuery(config)
	if err != nil {
		return Query{}, fmt.Errorf("error generating query: %w", err)
	}

	return Query{
		shouldBeExecuted: true,
		generatedQuery:   query,
	}, nil
}

func validateConfig(config Config) error {
	for i, tracker := range config.Trackers {
		if tracker.Id == "" {
			return fmt.Errorf("Trackers[%v].Identifier must not be empty", i)
		}

		if tracker.Repo == "" {
			return fmt.Errorf("Trackers[%v].Repo must not be empty", i)
		}

		if tracker.Owner == "" {
			return fmt.Errorf("Trackers[%v].Owner must not be empty", i)
		}

		if len(tracker.Authors) == 0 {
			return fmt.Errorf("Trackers[%v].Authors must have at least one author", i)
		}
	}

	if config.ReviewAmount <= 0 {
		return errors.New("ReviewAmount must not be zero or less")
	}

	if config.PrAmount <= 0 {
		return errors.New("PrAmount must not be zero or less")
	}

	if config.CommentsAmount <= 0 {
		return errors.New("CommentsAmount must not be zero or less")
	}

	return nil
}
