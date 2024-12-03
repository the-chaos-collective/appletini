package config

import (
	"fmt"
	"os"

	"appletini/config/migration"
	v1 "appletini/config/migration/v1"
	v2 "appletini/config/migration/v2"
	"appletini/logging"
)

type (
	Config   = v2.Config
	Tracking = v2.Tracking
)

var LatestVersion int = 2

type Loader struct {
	DumpMigrations bool
	Logger         logging.Logger
}

func (l Loader) Load(filename string) (Config, error) {
	_, err := os.ReadFile(filename)
	if err != nil { // If no config exists
		l.Logger.Info("Generating default config")
		defaultConfig := v1.Default() // Load default v1
		err = defaultConfig.Save(filename)
		if err != nil {
			return Config{}, fmt.Errorf("saving default config: %w", err)
		}
	}

	config, err := migration.MigrateTo(l.Logger, filename, LatestVersion, l.DumpMigrations)
	if err != nil {
		return Config{}, err
	}

	return config.(Config), err
}
