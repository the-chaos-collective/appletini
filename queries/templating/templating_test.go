package templating_test

import (
	"errors"
	"strings"
	"testing"

	"appletini/queries/templating"
)

func TestQuote(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "plain", input: "bug", expected: `"bug"`},
		{name: "empty", input: "", expected: `""`},
		{name: "apostrophe", input: "won't fix", expected: `"won't fix"`},
		{name: "ampersand and angle brackets", input: "R&D <ui>", expected: `"R&D <ui>"`},
		{name: "double quote", input: `say "hi"`, expected: `"say \"hi\""`},
		{name: "backslash", input: `a\b`, expected: `"a\\b"`},
		{name: "newline", input: "a\nb", expected: `"a\nb"`},
		{name: "unicode", input: "🐛 bug", expected: `"🐛 bug"`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := templating.Quote(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tc.expected {
				t.Fatalf("got %s, expected %s", got, tc.expected)
			}
		})
	}
}

func TestRender(t *testing.T) {
	type data struct {
		Label string
		Count int
	}

	testCases := []struct {
		name        string
		tpl         string
		data        any
		expected    string
		expectedErr error
		errContains string
	}{
		{
			name:     "quotes values",
			tpl:      `label(name: {{ quote .Label }}) first: {{ .Count }}`,
			data:     data{Label: `won't "fix" & R&D`, Count: 3},
			expected: `label(name: "won't \"fix\" & R&D") first: 3`,
		},
		{
			name:     "quote of printf",
			tpl:      `{{ quote (printf "repo:%s/%s" .Label .Label) }}`,
			data:     data{Label: "x"},
			expected: `"repo:x/x"`,
		},
		{
			name:        "empty template",
			tpl:         "  \n",
			data:        data{},
			expectedErr: templating.ErrEmptyTemplate,
		},
		{
			name:        "parse error",
			tpl:         "{{ .Label ",
			data:        data{},
			errContains: "parsing template",
		},
		{
			name:        "unknown field",
			tpl:         "{{ .Missing }}",
			data:        data{},
			errContains: "executing template",
		},
		{
			name:        "missing map key",
			tpl:         "{{ .Missing }}",
			data:        map[string]string{},
			errContains: "executing template",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := templating.Render("test", tc.tpl, tc.data)

			switch {
			case tc.expectedErr != nil:
				if !errors.Is(err, tc.expectedErr) {
					t.Fatalf("got error %v, expected %v", err, tc.expectedErr)
				}
			case tc.errContains != "":
				if err == nil || !strings.Contains(err.Error(), tc.errContains) {
					t.Fatalf("got error %v, expected it to contain %q", err, tc.errContains)
				}
			default:
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				if got != tc.expected {
					t.Fatalf("got %s, expected %s", got, tc.expected)
				}
			}
		})
	}
}
