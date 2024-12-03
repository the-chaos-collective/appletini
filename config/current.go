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

const LATEST_VERSION int = 2

type Loader struct {
	Migrator migration.Migrator
	Logger   logging.Logger
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

	config, err := l.Migrator.MigrateTo(filename, LATEST_VERSION)
	if err != nil {
		return Config{}, err
	}

	return config.(Config), err
}
