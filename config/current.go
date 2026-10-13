package config

import (
	"errors"
	"fmt"
	"io/fs"
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

var ErrUnexpectedConfigType = errors.New("migrated config has unexpected type")

func (l Loader) Load(filename string) (Config, error) {
	_, err := os.Stat(filename)

	switch {
	case errors.Is(err, fs.ErrNotExist):
		l.Logger.Info("Generating default config", "path", filename)
		defaultConfig := v1.Default() // Load default v1

		err = defaultConfig.Save(filename)
		if err != nil {
			return Config{}, fmt.Errorf("saving default config: %w", err)
		}
	case err != nil:
		// Exists but can't be inspected (e.g. permissions): never overwrite it.
		return Config{}, fmt.Errorf("checking config file: %w", err)
	}

	migrated, err := l.Migrator.MigrateTo(filename, LATEST_VERSION)
	if err != nil {
		return Config{}, err
	}

	config, ok := migrated.(Config)
	if !ok {
		return Config{}, fmt.Errorf("%w: %T", ErrUnexpectedConfigType, migrated)
	}

	return config, nil
}
