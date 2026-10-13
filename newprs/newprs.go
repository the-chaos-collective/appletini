// Package newprs spots pull requests that are new to review since earlier polls.
package newprs

import (
	"errors"
	"slices"
	"strings"

	"appletini/gitter"
	"appletini/trackerkey"
)

// Found is a new PR to review, with the result key of the tracker it was found in.
type Found struct {
	Key string
	PR  gitter.PullRequest
}

// Detector remembers every PR seen this session. It is not safe for
// concurrent use; the UI loop is its only caller.
type Detector struct {
	seen   map[string]struct{}
	primed bool
}

// Update records the PRs from one poll and returns the ones to notify about:
// PRs in label, repo or author trackers that haven't been seen before this
// session, excluding your own PRs and PRs you've already approved.
//
// The first call only records what's there (nothing is "new" on startup).
// PRs stay remembered even if they vanish from a later poll, so a partial
// poll after an error can't make old PRs look new again. Unparseable result
// keys are skipped and reported in the returned error.
func (d *Detector) Update(prs map[string][]gitter.PullRequest) ([]Found, error) {
	if d.seen == nil {
		d.seen = make(map[string]struct{})
	}

	var (
		found []Found
		errs  []error
	)

	for _, rawKey := range sortedKeys(prs) {
		key, err := trackerkey.Parse(rawKey)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		for _, pr := range prs[rawKey] {
			if pr.Id == "" {
				continue
			}

			_, known := d.seen[pr.Id]
			d.seen[pr.Id] = struct{}{}

			if known || !d.primed || !worthNotifying(key, pr) {
				continue
			}

			found = append(found, Found{Key: rawKey, PR: pr})
		}
	}

	d.primed = true

	return found, errors.Join(errs...)
}

// worthNotifying is whether a PR you haven't seen before is one to review.
func worthNotifying(key trackerkey.Key, pr gitter.PullRequest) bool {
	return key.Kind != trackerkey.Personal && !pr.ViewerAuthor && !pr.ApprovedByViewer()
}

// sortedKeys gives a stable order, so the same PR found in two trackers is
// always attributed to the same one.
func sortedKeys(prs map[string][]gitter.PullRequest) []string {
	keys := make([]string, 0, len(prs))
	for key := range prs {
		keys = append(keys, key)
	}

	slices.SortFunc(keys, strings.Compare)

	return keys
}
