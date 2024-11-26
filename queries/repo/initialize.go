package repo

import (
	"errors"
	"fmt"
)

func MakeRepoQuery(config Config) (RepoQuery, error) {
	err := validateConfig(config)
	if err != nil {
		return RepoQuery{}, fmt.Errorf("invalid config: %w", err)
	}

	query, err := generateQuery(config)
	if err != nil {
		return RepoQuery{}, fmt.Errorf("error generating query: %w", err)
	}

	return RepoQuery{
		generatedQuery: query,
	}, nil
}

func validateConfig(config Config) error {
	for i, repoConfig := range config.Trackers {
		if repoConfig.Id == "" {
			return fmt.Errorf("Trackers[%v].Identifier must not be empty", i)
		}

		if repoConfig.Repo == "" {
			return fmt.Errorf("Trackers[%v].Name must not be empty", i)
		}

		if repoConfig.Owner == "" {
			return fmt.Errorf("Trackers[%v].Owner must not be empty", i)
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
