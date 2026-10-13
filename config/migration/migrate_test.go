package migration

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/charmbracelet/log"

	v2 "appletini/config/migration/v2"
	"appletini/logging"
)

const tokenEnv = "APPLETINI_MIGRATION_TEST_TOKEN"

func writeConfig(t *testing.T, contents string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "config.json")

	err := os.WriteFile(path, []byte(contents), 0o600)
	if err != nil {
		t.Fatal(err)
	}

	return path
}

func TestReadVersion(t *testing.T) {
	testCases := []struct {
		name        string
		contents    *string
		expected    int
		errContains string
	}{
		{name: "no version field means v1", contents: ptr(`{"itemCount": 1}`), expected: 1},
		{name: "explicit zero means v1", contents: ptr(`{"__version": 0}`), expected: 1},
		{name: "v2", contents: ptr(`{"__version": 2}`), expected: 2},
		{name: "future version", contents: ptr(`{"__version": 7}`), expected: 7},
		{name: "invalid json", contents: ptr(`{`), expected: -1, errContains: "unmarshaling config"},
		{name: "wrong type", contents: ptr(`{"__version": "2"}`), expected: -1, errContains: "unmarshaling config"},
		{name: "missing file", contents: nil, expected: -1, errContains: "reading config file"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "missing.json")
			if tc.contents != nil {
				path = writeConfig(t, *tc.contents)
			}

			//exhaustruct:ignore
			got, err := Migrator{}.readVersion(path)

			if tc.errContains != "" && (err == nil || !strings.Contains(err.Error(), tc.errContains)) {
				t.Fatalf("got error %v, expected it to contain %q", err, tc.errContains)
			}

			if tc.errContains == "" && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tc.expected {
				t.Fatalf("got version %d, expected %d", got, tc.expected)
			}
		})
	}
}

func TestMigrateTo(t *testing.T) {
	v1 := `{"github": {"gqlAPI": "https://x", "token": "` + tokenEnv + `"}, "poll": {"frequencySeconds": 10},
	        "tracking": {"byLabel": [], "byRepo": []}, "itemCount": 3, "darkMode": false}`
	v2Contents := `{"__version": 2, "github": {"gqlAPI": "https://x", "token": "` + tokenEnv + `"},
	        "poll": {"frequencySeconds": 10}, "tracking": {}, "itemCount": 3}`

	testCases := []struct {
		name           string
		contents       string
		dump           bool
		setToken       bool
		expectedBackup bool
		errContains    string
	}{
		{name: "v1 to v2", contents: v1, setToken: true},
		{name: "v1 to v2 keeping a backup", contents: v1, dump: true, setToken: true, expectedBackup: true},
		{name: "already latest", contents: v2Contents, setToken: true},
		{name: "migration fails without token", contents: v1, setToken: false, errContains: "migrating config"},
		{name: "unknown version", contents: `{"__version": 99}`, setToken: true, errContains: "no loader for config v99"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setToken {
				t.Setenv(tokenEnv, "token")
			} else {
				t.Setenv(tokenEnv, "")
				os.Unsetenv(tokenEnv)
			}

			path := writeConfig(t, tc.contents)
			m := Migrator{DumpMigrations: tc.dump, Logger: logging.NewCharm(log.New(io.Discard))}

			got, err := m.MigrateTo(path, 2)

			if tc.errContains != "" {
				if err == nil || !strings.Contains(err.Error(), tc.errContains) {
					t.Fatalf("got error %v, expected it to contain %q", err, tc.errContains)
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			conf, ok := got.(v2.Config)
			if !ok {
				t.Fatalf("got %T, expected v2.Config", got)
			}

			if conf.Version != 2 || conf.ItemCount != 3 {
				t.Fatalf("unexpected config %+v", conf)
			}

			_, err = os.Stat(strings.ReplaceAll(path, ".json", ".v1.json"))
			if backup := err == nil; backup != tc.expectedBackup {
				t.Fatalf("backup exists = %v, expected %v", backup, tc.expectedBackup)
			}
		})
	}
}

func ptr(s string) *string { return &s }
