package v1

import (
	"fmt"

	"git_applet/config/migration/migration_types"
	v2 "git_applet/config/migration/v2"
)

func (config Config) ToNext() (migration_types.Migratable, error) {
	new := v2.Config{
		Version: 2,
		Github: v2.GithubConfig{
			GraphQL: config.Github.GraphQL,
			Token:   config.Github.Token,
		},
		Poll: v2.PollConfig{
			Frequency: config.Poll.Frequency,
		},
		Tracking: v2.Tracking{
			Personal: true,
			ByLabel:  config.Tracking.ByLabel.ToNext(),
			ByRepo:   config.Tracking.ByRepo.ToNext(),
			ByAuthor: v2.AuthorSet{
				v2.Author{
					Title: "The Appletini Team",
					Authors: []string{
						"darvoid",
						"tobyselway",
					},
					RepoName: "appletini",
					Owner:    "the-chaos-collective",
				},
			},
		},
		ItemCount: config.ItemCount,
		Darkmode:  config.Darkmode,
	}

	err := new.Setup()
	if err != nil {
		return new, fmt.Errorf("from v1 to v2: %w", err)
	}

	return new, nil
}

func (labelTracker LabeledRepo) ToNext() v2.Labeled {
	return v2.Labeled{
		Title:    labelTracker.Title,
		Label:    labelTracker.Label,
		Owner:    labelTracker.Owner,
		RepoName: labelTracker.RepoName,
	}
}

func (labelTrackerSet LabeledRepoSet) ToNext() v2.LabeledSet {
	new := make(v2.LabeledSet, 0)
	for _, tracker := range labelTrackerSet {
		new = append(new, tracker.ToNext())
	}
	return new
}

func (repoTracker Repo) ToNext() v2.Repo {
	return v2.Repo{
		Title:    repoTracker.Title,
		Owner:    repoTracker.Owner,
		RepoName: repoTracker.RepoName,
	}
}

func (repoTrackerSet RepoSet) ToNext() v2.RepoSet {
	new := make(v2.RepoSet, 0)
	for _, tracker := range repoTrackerSet {
		new = append(new, tracker.ToNext())
	}
	return new
}
