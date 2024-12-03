package main

type FeatureFlags struct {
	DumpMigrations bool
	MockQueries    bool
}

func featureFlags() FeatureFlags {
	return FeatureFlags{
		DumpMigrations: false,
		MockQueries:    false,
	}
}
