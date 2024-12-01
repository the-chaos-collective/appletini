package config

import (
	"fmt"
	"os"
	"strings"

	"git_applet/config/migration"
	v1 "git_applet/config/migration/v1"
	v2 "git_applet/config/migration/v2"
)

type (
	Config   = v2.Config
	Tracking = v2.Tracking
)

func Load(filename string, dumpMigrations bool) (Config, error) {
	_, err := os.ReadFile(filename)
	if err != nil {
		defaultConfig := v1.Default()
		err = defaultConfig.Save(filename)
		if err != nil {
			return Config{}, fmt.Errorf("saving default config: %w", err)
		}
	}

	currentVersion, err := migration.ReadVersion(filename)
	if err != nil {
		return Config{}, fmt.Errorf("determining current version: %w", err)
	}

	switch currentVersion {
	case 0:
		fallthrough
	case 1:
		old, err := v1.Load(filename)
		if err != nil {
			return Config{}, fmt.Errorf("reading old config: %w", err)
		}

		new := old.ToNext()

		if dumpMigrations {
			err = os.Rename(filename, strings.ReplaceAll(filename, ".json", ".v1.json"))
			if err != nil {
				return Config{}, fmt.Errorf("renaming old config: %w", err)
			}
		}

		err = new.Save(filename)
		if err != nil {
			return Config{}, fmt.Errorf("saving new config: %w", err)
		}

		return new, nil

	case 2:
		return v2.Load(filename)
	}

	return Config{}, fmt.Errorf("unable to load config (v%d)", currentVersion)
}
