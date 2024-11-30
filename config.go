package main

import (
	"encoding/json"
	"fmt"
	"git_applet/global_types"
	"os"
)

func loadConfig() (global_types.Config, error) {
	file_contents, err := os.ReadFile(CONFIG_FILE)
	if err != nil {
		return global_types.Config{}, fmt.Errorf("reading config file: %w", err)
	}
	config := global_types.Config{}
	err = json.Unmarshal(file_contents, &config)
	if err != nil {
		return global_types.Config{}, fmt.Errorf("unmarshaling config: %w", err)
	}
	// fmt.Printf("Config: %v\n", Config.Github.Token)
	return config, nil
}
