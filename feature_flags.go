package main

type FeatureFlags struct {
	DumpMigrations bool
	MockQueries    bool
}

func LoadFeatureFlags() FeatureFlags {
	return FeatureFlags{
		DumpMigrations: false,
		MockQueries:    false,
	}
}
