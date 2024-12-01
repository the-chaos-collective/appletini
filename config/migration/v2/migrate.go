package v2

import (
	"git_applet/config/migration/migration_types"
)

func (config Config) ToNext() (migration_types.Migratable, error) {
	// TODO: Implement to migrate to v3
	return config, nil
}
