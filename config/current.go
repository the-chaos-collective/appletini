package config

import (
	"fmt"
	"log"
	"os"

	"git_applet/config/migration"
	v1 "git_applet/config/migration/v1"
	v2 "git_applet/config/migration/v2"
)

type (
	Config   = v2.Config
	Tracking = v2.Tracking
)

var LatestVersion int = 2

func Load(filename string, dumpMigrations bool, logger *log.Logger) (Config, error) {
	_, err := os.ReadFile(filename)
	if err != nil { // If no config exists
		logger.Println("Generating default config")
		defaultConfig := v1.Default() // Load default v1
		err = defaultConfig.Save(filename)
		if err != nil {
			return Config{}, fmt.Errorf("saving default config: %w", err)
		}
	}

	config, err := migration.MigrateTo(filename, LatestVersion, dumpMigrations, logger)
	if err != nil {
		return Config{}, err
	}

	return config.(Config), err
}
