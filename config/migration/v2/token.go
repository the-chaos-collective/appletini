package v2

import (
	"fmt"
	"os"
)

func (c *Config) loadGithubToken() error {
	token, present := os.LookupEnv(c.Github.Token)
	if !present {
		return fmt.Errorf("token not present in %s env var", c.Github.Token)
	}

	c.Computed.GithubToken = token

	return nil
}
