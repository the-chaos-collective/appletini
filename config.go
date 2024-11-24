package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func loadConfig() error {
	file_contents, err := os.ReadFile(CONFIG_FILE)
	if err != nil {
		return fmt.Errorf("reading config file: %w", err)
	}
	err = json.Unmarshal(file_contents, &Config)
	if err != nil {
		return fmt.Errorf("unmarshaling config: %w", err)
	}
	fmt.Printf("%v\n", Config.Tracking)

	return nil
}
