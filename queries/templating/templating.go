// Package templating renders the GraphQL query templates used by the tracker
// queries. It uses text/template (not html/template, which HTML-escapes values
// such as "won't fix" into "won&#39;t fix") and exposes a `quote` function that
// turns any value into a valid GraphQL string literal.
package templating

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"text/template"
)

var ErrEmptyTemplate = errors.New("template must not be empty")

// Quote returns s as a double-quoted GraphQL string literal.
// GraphQL string escapes are a superset of JSON's, so JSON encoding is safe.
func Quote(s string) (string, error) {
	var buf bytes.Buffer

	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)

	err := enc.Encode(s)
	if err != nil {
		return "", fmt.Errorf("quoting %q: %w", s, err)
	}

	return strings.TrimSuffix(buf.String(), "\n"), nil
}

// Render parses tpl under name and executes it against data.
func Render(name, tpl string, data any) (string, error) {
	if strings.TrimSpace(tpl) == "" {
		return "", ErrEmptyTemplate
	}

	parsed, err := template.New(name).
		Option("missingkey=error").
		Funcs(template.FuncMap{"quote": Quote}).
		Parse(tpl)
	if err != nil {
		return "", fmt.Errorf("parsing template %s: %w", name, err)
	}

	var out bytes.Buffer

	err = parsed.Execute(&out, data)
	if err != nil {
		return "", fmt.Errorf("executing template %s: %w", name, err)
	}

	return out.String(), nil
}
