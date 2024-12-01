package v2

import (
	"encoding/json"
	"fmt"
	"os"
)

func (config *Config) Setup() error {
	err := config.loadGithubToken()
	if err != nil {
		return fmt.Errorf("loading GitHub token: %w", err)
	}
	return nil
}

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

	err = config.Setup()
	if err != nil {
		return Config{}, fmt.Errorf("config setup: %w", err)
	}

	return config, nil
}
