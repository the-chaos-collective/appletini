package v2_test

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	v2 "appletini/config/migration/v2"
)

const tokenEnv = "APPLETINI_V2_TEST_TOKEN"

func sample() v2.Config {
	//exhaustruct:ignore
	return v2.Config{
		Version:   2,
		Github:    v2.GithubConfig{GraphQL: "https://api.github.com/graphql", Token: tokenEnv},
		Poll:      v2.PollConfig{Frequency: 30},
		ItemCount: 10,
		Darkmode:  true,
		Tracking: v2.Tracking{
			Personal: true,
			ByLabel:  v2.LabeledSet{{Title: "Bugs", Owner: "o", RepoName: "r", Label: "won't fix"}},
			ByRepo:   v2.RepoSet{{Title: "Repo", Owner: "o", RepoName: "r"}},
			ByAuthor: v2.AuthorSet{{Title: "Team", Owner: "o", RepoName: "r", Authors: []string{"a", "b"}}},
		},
	}
}

func TestSaveLoadRoundTrip(t *testing.T) {
	t.Setenv(tokenEnv, "tok")

	path := filepath.Join(t.TempDir(), "config.json")
	original := sample()

	err := original.Save(path)
	if err != nil {
		t.Fatalf("Save: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}

	if perm := info.Mode().Perm(); perm != 0o600 {
		t.Fatalf("config saved with permissions %o, expected 600", perm)
	}

	loaded, err := v2.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	original.Computed.GithubToken = "tok" // filled from the env var on load
	if !reflect.DeepEqual(loaded, original) {
		t.Fatalf("round trip changed config:\n got %+v\nwant %+v", loaded, original)
	}
}

func TestComputedTokenIsNotSaved(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	conf := sample()
	conf.Computed.GithubToken = "super-secret"

	err := conf.Save(path)
	if err != nil {
		t.Fatalf("Save: %v", err)
	}

	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if strings.Contains(string(contents), "super-secret") {
		t.Fatal("the token value was written to the config file")
	}
}

func TestLoadErrors(t *testing.T) {
	dir := t.TempDir()

	write := func(name, contents string) string {
		path := filepath.Join(dir, name)

		err := os.WriteFile(path, []byte(contents), 0o600)
		if err != nil {
			t.Fatal(err)
		}

		return path
	}

	testCases := []struct {
		name        string
		path        string
		errContains string
	}{
		{name: "missing file", path: filepath.Join(dir, "missing.json"), errContains: "reading config file"},
		{name: "invalid json", path: write("invalid.json", "{"), errContains: "unmarshaling config"},
		{
			name:        "token env var missing",
			path:        write("notoken.json", `{"__version": 2, "github": {"token": "APPLETINI_V2_UNSET_VAR"}}`),
			errContains: "token not present in APPLETINI_V2_UNSET_VAR env var",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := v2.Load(tc.path)
			if err == nil || !strings.Contains(err.Error(), tc.errContains) {
				t.Fatalf("got error %v, expected it to contain %q", err, tc.errContains)
			}
		})
	}
}

func TestToNextIsIdentityForLatest(t *testing.T) {
	next, err := sample().ToNext()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(next, sample()) {
		t.Fatalf("got %+v, expected unchanged config", next)
	}
}
