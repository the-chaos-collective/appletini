// Package trackerkey defines the keys that identify each tracker's results.
//
// The poller builds keys ("label_0", "repo_2", "author_1", "personal") and
// the UI parses them back to look up the tracker's title, so the format lives
// here in one place.
package trackerkey

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

type Kind string

const (
	Personal Kind = "personal"
	Label    Kind = "label"
	Repo     Kind = "repo"
	Author   Kind = "author"
)

// order is the display order of tracker kinds in the menu.
var order = []Kind{Personal, Label, Repo, Author}

var (
	ErrNegativeIndex = errors.New("tracker index must not be negative")
	ErrUnknownKind   = errors.New("unknown tracker kind")
	ErrMalformed     = errors.New("malformed tracker key (expected [kind]_[index] or personal)")
)

// Key identifies one tracker's result set.
type Key struct {
	Kind  Kind
	Index int
}

// New returns the key for the idx-th tracker of kind.
func New(kind Kind, idx int) (string, error) {
	if idx < 0 {
		return "", ErrNegativeIndex
	}

	switch kind {
	case Personal:
		return string(Personal), nil
	case Label, Repo, Author:
		return fmt.Sprintf("%s_%d", kind, idx), nil
	}

	return "", fmt.Errorf("%w: %q", ErrUnknownKind, kind)
}

// Parse is the inverse of New.
func Parse(raw string) (Key, error) {
	if raw == string(Personal) {
		return Key{Kind: Personal, Index: 0}, nil
	}

	kind, idxStr, found := strings.Cut(raw, "_")
	if !found || kind == "" || idxStr == "" {
		return Key{}, fmt.Errorf("%w: %q", ErrMalformed, raw)
	}

	switch Kind(kind) {
	case Label, Repo, Author:
	case Personal: // "personal" takes no index, so "personal_0" is malformed
		return Key{}, fmt.Errorf("%w: %q", ErrMalformed, raw)
	default:
		return Key{}, fmt.Errorf("%w: %q", ErrUnknownKind, kind)
	}

	idx, err := strconv.Atoi(idxStr)
	if err != nil {
		return Key{}, fmt.Errorf("%w: %q: %w", ErrMalformed, raw, err)
	}

	if idx < 0 {
		return Key{}, fmt.Errorf("%w: %q", ErrNegativeIndex, raw)
	}

	return Key{Kind: Kind(kind), Index: idx}, nil
}

// Compare orders keys for display: personal first, then label, repo and
// author trackers, each in the order they appear in the config.
func Compare(a, b Key) int {
	if a.Kind != b.Kind {
		return slices.Index(order, a.Kind) - slices.Index(order, b.Kind)
	}

	return a.Index - b.Index
}
