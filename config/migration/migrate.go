package migration

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"appletini/config/migration/migration_types"
	v1 "appletini/config/migration/v1"
	v2 "appletini/config/migration/v2"
	"appletini/logging"
)

type Migrator struct {
	DumpMigrations bool
	Logger         logging.Logger
}

func (m Migrator) loadAsVersion(filename string, version int) (migration_types.Migratable, error) {
	switch version {
	case 0:
		fallthrough
	case 1:
		return v1.Load(filename)
	case 2:
		return v2.Load(filename)
	}

	return nil, fmt.Errorf("no loader for config v%d", version)
}

func (m Migrator) MigrateTo(filename string, targetVersion int) (migration_types.Migratable, error) {
	logged := false

	for {
		version, err := m.readVersion(filename)
		if err != nil {
			return nil, fmt.Errorf("determining current version: %w", err)
		}

		current, err := m.loadAsVersion(filename, version)
		if err != nil {
			return nil, fmt.Errorf("loading config (as v%d): %w", version, err)
		}

		if version >= targetVersion {
			return current, nil
		}

		if !logged {
			m.Logger.Info("Migrating config", "from", fmt.Sprintf("v%d", version), "to", fmt.Sprintf("v%d", targetVersion))
			logged = true
		}

		new, err := current.ToNext()
		if err != nil {
			return nil, fmt.Errorf("migrating config: %w", err)
		}

		if m.DumpMigrations {
			err = os.Rename(filename, strings.ReplaceAll(filename, ".json", fmt.Sprintf(".v%d.json", version)))
			if err != nil {
				return nil, fmt.Errorf("renaming old config: %w", err)
			}
		}

		err = new.Save(filename)
		if err != nil {
			return nil, fmt.Errorf("saving new config: %w", err)
		}

	}
}

func (m Migrator) readVersion(filename string) (int, error) {
	type ConfigVersion struct {
		Version int `json:"__version"`
	}

	file_contents, err := os.ReadFile(filename)
	if err != nil {
		return -1, fmt.Errorf("reading config file: %w", err)
	}

	var config ConfigVersion

	err = json.Unmarshal(file_contents, &config)
	if err != nil {
		return -1, fmt.Errorf("unmarshaling config: %w", err)
	}

	if config.Version == 0 {
		return 1, nil
	}

	return config.Version, nil
}
