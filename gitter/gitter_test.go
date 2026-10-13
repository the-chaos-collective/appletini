package gitter_test

import (
	"encoding/json"
	"testing"

	"appletini/gitter"
	"appletini/status"
)

func TestPRInfo(t *testing.T) {
	testCases := []struct {
		name           string
		reviewDecision string
		mergeable      string
		expected       status.PRInfo
	}{
		{
			name:           "approved and mergeable",
			reviewDecision: "APPROVED",
			mergeable:      "MERGEABLE",
			expected:       status.PRInfo{Review: status.ReviewState_Approved, Mergeable: status.MergeableState_Mergeable},
		},
		{
			name:           "review required and conflicting",
			reviewDecision: "REVIEW_REQUIRED",
			mergeable:      "CONFLICTING",
			expected:       status.PRInfo{Review: status.ReviewState_RequiresReview, Mergeable: status.MergeableState_Conflict},
		},
		{
			name:           "changes requested, unknown mergeability",
			reviewDecision: "CHANGES_REQUESTED",
			mergeable:      "UNKNOWN",
			expected:       status.PRInfo{Review: status.ReviewState_ChangesRequested, Mergeable: status.MergeableState_Unknown},
		},
		{
			name:           "no review decision means none required",
			reviewDecision: "",
			mergeable:      "MERGEABLE",
			expected:       status.PRInfo{Review: status.ReviewState_NoReviewRequired, Mergeable: status.MergeableState_Mergeable},
		},
		{
			name:           "unrecognised values fall back to unknown",
			reviewDecision: "SOMETHING_NEW",
			mergeable:      "ALSO_NEW",
			expected:       status.PRInfo{Review: status.ReviewState_Unknown, Mergeable: status.MergeableState_Unknown},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//exhaustruct:ignore
			pr := gitter.PullRequest{ReviewDecision: tc.reviewDecision, Mergeable: tc.mergeable}

			if got := pr.PRInfo(); got != tc.expected {
				t.Fatalf("got %+v, expected %+v", got, tc.expected)
			}
		})
	}
}

func TestPRInfoChecks(t *testing.T) {
	testCases := []struct {
		state    string
		expected status.ChecksState
	}{
		{state: "", expected: status.ChecksState_None},
		{state: "SUCCESS", expected: status.ChecksState_Passing},
		{state: "FAILURE", expected: status.ChecksState_Failing},
		{state: "ERROR", expected: status.ChecksState_Failing},
		{state: "PENDING", expected: status.ChecksState_Pending},
		{state: "EXPECTED", expected: status.ChecksState_Pending},
		{state: "SOMETHING_NEW", expected: status.ChecksState_Unknown},
	}

	for _, tc := range testCases {
		t.Run("state "+tc.state, func(t *testing.T) {
			//exhaustruct:ignore
			pr := gitter.PullRequest{ReviewDecision: "APPROVED", Mergeable: "MERGEABLE", ChecksState: tc.state}

			if got := pr.PRInfo().Checks; got != tc.expected {
				t.Fatalf("got %v, expected %v", got, tc.expected)
			}
		})
	}
}

func TestLastCommitChecksState(t *testing.T) {
	testCases := []struct {
		name     string
		json     string
		expected string
	}{
		{name: "state present", json: `{"nodes":[{"commit":{"statusCheckRollup":{"state":"SUCCESS"}}}]}`, expected: "SUCCESS"},
		{name: "no checks ran (null rollup)", json: `{"nodes":[{"commit":{"statusCheckRollup":null}}]}`, expected: ""},
		{name: "no commits", json: `{"nodes":[]}`, expected: ""},
		{name: "field missing entirely", json: `{}`, expected: ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var c gitter.LastCommit

			err := json.Unmarshal([]byte(tc.json), &c)
			if err != nil {
				t.Fatalf("unmarshal: %v", err)
			}

			if got := c.ChecksState(); got != tc.expected {
				t.Fatalf("got %q, expected %q", got, tc.expected)
			}
		})
	}
}
