package gitter_test

import (
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
