package config

import (
	"encoding/json"
	"fmt"
	"os"
)

func Load(filename string) (Config, error) {
	file_contents, err := os.ReadFile(filename)
	if err != nil {
		return Config{}, fmt.Errorf("reading config file: %w", err)
	}

	config := Config{}

	err = json.Unmarshal(file_contents, &config)
	if err != nil {
		return Config{}, fmt.Errorf("unmarshaling config: %w", err)
	}

	err = config.loadGithubToken()
	if err != nil {
		return Config{}, fmt.Errorf("loading GitHub token: %w", err)
	}

	return config, nil
}
