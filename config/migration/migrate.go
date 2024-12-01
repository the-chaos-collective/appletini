package migration

import (
	"encoding/json"
	"fmt"
	"os"
)

type ConfigVersion struct {
	Version int `json:"__version"`
}

func ReadVersion(filename string) (int, error) {
	file_contents, err := os.ReadFile(filename)
	if err != nil {
		return -1, fmt.Errorf("reading config file: %w", err)
	}

	config := ConfigVersion{}

	err = json.Unmarshal(file_contents, &config)
	if err != nil {
		return -1, fmt.Errorf("unmarshaling config: %w", err)
	}

	return config.Version, nil
}
