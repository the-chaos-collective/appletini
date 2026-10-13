package config_test

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/charmbracelet/log"

	"appletini/config"
	"appletini/config/migration"
	"appletini/logging"
)

const tokenEnv = "APPLETINI_TEST_TOKEN"

func newLoader() config.Loader {
	logger := logging.NewCharm(log.New(io.Discard))

	return config.Loader{
		Migrator: migration.Migrator{DumpMigrations: false, Logger: logger},
		Logger:   logger,
	}
}

func TestLoad(t *testing.T) {
	v2Config := `{
	  "__version": 2,
	  "github": {"gqlAPI": "https://api.github.com/graphql", "token": "` + tokenEnv + `"},
	  "poll": {"frequencySeconds": 60},
	  "tracking": {"personal": true, "byLabel": [], "byRepo": [], "byAuthor": []},
	  "itemCount": 5,
	  "darkMode": false
	}`

	v1Config := `{
	  "github": {"gqlAPI": "https://api.github.com/graphql", "token": "` + tokenEnv + `"},
	  "poll": {"frequencySeconds": 45},
	  "tracking": {"byLabel": [{"title": "Bugs", "owner": "o", "repo": "r", "label": "bug"}], "byRepo": []},
	  "itemCount": 7,
	  "darkMode": true
	}`

	testCases := []struct {
		name        string
		contents    *string // nil: file does not exist
		setToken    bool
		errContains string
		check       func(t *testing.T, conf config.Config, path string)
	}{
		{
			name:     "existing v2 config",
			contents: &v2Config,
			setToken: true,
			check: func(t *testing.T, conf config.Config, _ string) {
				if conf.Version != 2 || conf.ItemCount != 5 || conf.Poll.Frequency != 60 || !conf.Tracking.Personal {
					t.Fatalf("unexpected config %+v", conf)
				}

				if conf.Computed.GithubToken != "the-token" {
					t.Fatalf("token not loaded: %q", conf.Computed.GithubToken)
				}
			},
		},
		{
			name:     "v1 config is migrated and rewritten",
			contents: &v1Config,
			setToken: true,
			check: func(t *testing.T, conf config.Config, path string) {
				if conf.Version != 2 || conf.ItemCount != 7 || conf.Poll.Frequency != 45 || !conf.Darkmode {
					t.Fatalf("unexpected config %+v", conf)
				}

				if len(conf.Tracking.ByLabel) != 1 || conf.Tracking.ByLabel[0].Label != "bug" {
					t.Fatalf("label trackers not migrated: %+v", conf.Tracking.ByLabel)
				}

				written, err := os.ReadFile(path)
				if err != nil {
					t.Fatalf("reading migrated file: %v", err)
				}

				if !strings.Contains(string(written), `"__version": 2`) {
					t.Fatalf("migrated file was not saved as v2:\n%s", written)
				}
			},
		},
		{
			name:     "missing file creates the default config",
			contents: nil,
			setToken: true,
			check: func(t *testing.T, conf config.Config, path string) {
				if conf.Version != 2 || conf.Github.GraphQL != "https://api.github.com/graphql" {
					t.Fatalf("unexpected default config %+v", conf)
				}

				_, err := os.Stat(path)
				if err != nil {
					t.Fatalf("default config was not written: %v", err)
				}
			},
		},
		{
			name:        "token env var not set",
			contents:    &v2Config,
			setToken:    false,
			errContains: "token not present in " + tokenEnv,
		},
		{
			name:        "invalid json",
			contents:    ptr("{not json"),
			setToken:    true,
			errContains: "unmarshaling config",
		},
		{
			name:        "unknown future version",
			contents:    ptr(`{"__version": 99}`),
			setToken:    true,
			errContains: "no loader for config v99",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "config.json")

			if tc.contents != nil {
				err := os.WriteFile(path, []byte(*tc.contents), 0o600)
				if err != nil {
					t.Fatal(err)
				}
			}

			if tc.setToken {
				t.Setenv(tokenEnv, "the-token")
				t.Setenv("GITHUB_ACCESS_TOKEN", "the-token") // name used by the default config
			} else {
				t.Setenv(tokenEnv, "")
				os.Unsetenv(tokenEnv)
			}

			conf, err := newLoader().Load(path)

			if tc.errContains != "" {
				if err == nil || !strings.Contains(err.Error(), tc.errContains) {
					t.Fatalf("got error %v, expected it to contain %q", err, tc.errContains)
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			tc.check(t, conf, path)
		})
	}
}

func TestLoadDoesNotOverwriteUninspectableConfig(t *testing.T) {
	// A path whose parent is a file: Stat fails with something other than
	// "not exist", and Load must not try to write a default config there.
	dir := t.TempDir()
	parent := filepath.Join(dir, "not-a-dir")

	err := os.WriteFile(parent, []byte("x"), 0o600)
	if err != nil {
		t.Fatal(err)
	}

	_, err = newLoader().Load(filepath.Join(parent, "config.json"))
	if err == nil || !strings.Contains(err.Error(), "checking config file") {
		t.Fatalf("got error %v, expected a 'checking config file' error", err)
	}

	if errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("error should not be treated as a missing file: %v", err)
	}
}

func ptr(s string) *string { return &s }
