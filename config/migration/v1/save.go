package v1

import (
	"encoding/json"
	"fmt"
	"os"
)

func (config Config) Save(filename string) error {
	file_contents, err := json.MarshalIndent(&config, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}

	err = os.WriteFile(filename, file_contents, 0o600)
	if err != nil {
		return fmt.Errorf("writing config file: %w", err)
	}

	return nil
}
