package v1_test

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	v1 "appletini/config/migration/v1"
	v2 "appletini/config/migration/v2"
)

func TestSaveLoadRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	original := v1.Default()

	err := original.Save(path)
	if err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := v1.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if !reflect.DeepEqual(loaded, original) {
		t.Fatalf("round trip changed config:\n got %+v\nwant %+v", loaded, original)
	}
}

func TestLoadErrors(t *testing.T) {
	dir := t.TempDir()
	invalid := filepath.Join(dir, "invalid.json")

	err := os.WriteFile(invalid, []byte("{"), 0o600)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name        string
		path        string
		errContains string
	}{
		{name: "missing file", path: filepath.Join(dir, "missing.json"), errContains: "reading config file"},
		{name: "invalid json", path: invalid, errContains: "unmarshaling config"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := v1.Load(tc.path)
			if err == nil || !strings.Contains(err.Error(), tc.errContains) {
				t.Fatalf("got error %v, expected it to contain %q", err, tc.errContains)
			}
		})
	}
}

func TestSaveError(t *testing.T) {
	err := v1.Default().Save(filepath.Join(t.TempDir(), "no-such-dir", "config.json"))
	if err == nil || !strings.Contains(err.Error(), "writing config file") {
		t.Fatalf("got error %v, expected a write error", err)
	}
}

func TestToNext(t *testing.T) {
	const tokenEnv = "APPLETINI_V1_TEST_TOKEN"

	conf := v1.Default()
	conf.Github.Token = tokenEnv

	testCases := []struct {
		name        string
		setToken    bool
		errContains string
	}{
		{name: "migrates", setToken: true},
		{name: "fails without token", setToken: false, errContains: "from v1 to v2"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setToken {
				t.Setenv(tokenEnv, "tok")
			} else {
				t.Setenv(tokenEnv, "")
				os.Unsetenv(tokenEnv)
			}

			next, err := conf.ToNext()

			if tc.errContains != "" {
				if err == nil || !strings.Contains(err.Error(), tc.errContains) {
					t.Fatalf("got error %v, expected it to contain %q", err, tc.errContains)
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			got, ok := next.(v2.Config)
			if !ok {
				t.Fatalf("got %T, expected v2.Config", next)
			}

			if got.Version != 2 || !got.Tracking.Personal || got.Computed.GithubToken != "tok" {
				t.Fatalf("unexpected v2 config %+v", got)
			}

			if got.Poll.Frequency != conf.Poll.Frequency || got.ItemCount != conf.ItemCount || got.Darkmode != conf.Darkmode {
				t.Fatalf("scalar fields not carried over: %+v", got)
			}

			if len(got.Tracking.ByLabel) != len(conf.Tracking.ByLabel) || len(got.Tracking.ByRepo) != len(conf.Tracking.ByRepo) {
				t.Fatalf("trackers not carried over: %+v", got.Tracking)
			}

			if got.Tracking.ByLabel[0].Label != conf.Tracking.ByLabel[0].Label {
				t.Fatalf("label tracker changed: %+v", got.Tracking.ByLabel[0])
			}
		})
	}
}
