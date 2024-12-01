package v1

func Default() Config {
	return Config{
		Github: GithubConfig{
			GraphQL: "https://api.github.com/graphql",
			Token:   "GITHUB_ACCESS_TOKEN",
		},
		Poll: PollConfig{
			Frequency: 30,
		},
		Tracking: Tracking{
			ByLabel: LabeledRepoSet{
				LabeledRepo{
					Title:    "Appletini Bugs",
					Label:    "bug",
					RepoName: "appletini",
					Owner:    "the-chaos-collective",
				},
			},
			ByRepo: RepoSet{
				Repo{
					Title:    "Appletini",
					RepoName: "appletini",
					Owner:    "the-chaos-collective",
				},
			},
		},
		Darkmode:  false,
		ItemCount: 10,
	}
}
