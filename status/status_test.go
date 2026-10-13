package status_test

import (
	"reflect"
	"testing"

	"appletini/status"
)

func TestClassify(t *testing.T) {
	testCases := []struct {
		name     string
		input    status.PRInfo
		expected status.Status
	}{
		{
			name:  "approved and mergeable is green",
			input: status.PRInfo{Review: status.ReviewState_Approved, Mergeable: status.MergeableState_Mergeable},
			expected: status.Status{
				ShowGreenIcon: true, ShowRedIcon: false,
				Emoji: []string{"🟢"}, Message: "Approved | Can be merged",
			},
		},
		{
			name:  "no review required and unknown mergeability is still green",
			input: status.PRInfo{Review: status.ReviewState_NoReviewRequired, Mergeable: status.MergeableState_Unknown},
			expected: status.Status{
				ShowGreenIcon: true, ShowRedIcon: false,
				Emoji: []string{"🟢", "❓"}, Message: "No review required | Unknown",
			},
		},
		{
			name:  "changes requested is red",
			input: status.PRInfo{Review: status.ReviewState_ChangesRequested, Mergeable: status.MergeableState_Mergeable},
			expected: status.Status{
				ShowGreenIcon: false, ShowRedIcon: true,
				Emoji: []string{"🔴", "📝"}, Message: "Changes requested | Can be merged",
			},
		},
		{
			name:  "conflict is red even when approved",
			input: status.PRInfo{Review: status.ReviewState_Approved, Mergeable: status.MergeableState_Conflict},
			expected: status.Status{
				ShowGreenIcon: false, ShowRedIcon: true,
				Emoji: []string{"🔴", "🚧"}, Message: "Approved | Conflict",
			},
		},
		{
			name:  "review required is neither green nor red",
			input: status.PRInfo{Review: status.ReviewState_RequiresReview, Mergeable: status.MergeableState_Mergeable},
			expected: status.Status{
				ShowGreenIcon: false, ShowRedIcon: false,
				Emoji: []string{"🔴", "👁‍🗨"}, Message: "Review required | Can be merged",
			},
		},
		{
			name:  "everything unknown",
			input: status.PRInfo{Review: status.ReviewState_Unknown, Mergeable: status.MergeableState_Unknown},
			expected: status.Status{
				ShowGreenIcon: false, ShowRedIcon: false,
				Emoji: []string{"🔴", "❓", "❔"}, Message: "Unknown | Unknown",
			},
		},
		{
			name:  "changes requested and conflict",
			input: status.PRInfo{Review: status.ReviewState_ChangesRequested, Mergeable: status.MergeableState_Conflict},
			expected: status.Status{
				ShowGreenIcon: false, ShowRedIcon: true,
				Emoji: []string{"🔴", "🚧", "📝"}, Message: "Changes requested | Conflict",
			},
		},
		{
			name: "ready with checks passing is green",
			input: status.PRInfo{
				Review: status.ReviewState_Approved, Mergeable: status.MergeableState_Mergeable, Checks: status.ChecksState_Passing,
			},
			expected: status.Status{
				ShowGreenIcon: true, ShowRedIcon: false,
				Emoji: []string{"🟢"}, Message: "Approved | Can be merged | Checks passing",
			},
		},
		{
			name: "failing checks are red even when approved and mergeable",
			input: status.PRInfo{
				Review: status.ReviewState_Approved, Mergeable: status.MergeableState_Mergeable, Checks: status.ChecksState_Failing,
			},
			expected: status.Status{
				ShowGreenIcon: false, ShowRedIcon: true,
				Emoji: []string{"🔴", "❌"}, Message: "Approved | Can be merged | Checks failing",
			},
		},
		{
			name: "running checks block ready to merge but are not red",
			input: status.PRInfo{
				Review: status.ReviewState_Approved, Mergeable: status.MergeableState_Mergeable, Checks: status.ChecksState_Pending,
			},
			expected: status.Status{
				ShowGreenIcon: false, ShowRedIcon: false,
				Emoji: []string{"🔴", "⏳"}, Message: "Approved | Can be merged | Checks running",
			},
		},
		{
			name: "unknown checks state is not ready",
			input: status.PRInfo{
				Review: status.ReviewState_NoReviewRequired, Mergeable: status.MergeableState_Mergeable, Checks: status.ChecksState_Unknown,
			},
			expected: status.Status{
				ShowGreenIcon: false, ShowRedIcon: false,
				Emoji: []string{"🔴", "❔"}, Message: "No review required | Can be merged | Checks unknown",
			},
		},
		{
			name: "everything wrong at once",
			input: status.PRInfo{
				Review: status.ReviewState_ChangesRequested, Mergeable: status.MergeableState_Conflict, Checks: status.ChecksState_Failing,
			},
			expected: status.Status{
				ShowGreenIcon: false, ShowRedIcon: true,
				Emoji: []string{"🔴", "🚧", "📝", "❌"}, Message: "Changes requested | Conflict | Checks failing",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.input.Classify(); !reflect.DeepEqual(got, tc.expected) {
				t.Fatalf("got %+v, expected %+v", got, tc.expected)
			}
		})
	}
}

func TestStrings(t *testing.T) {
	testCases := []struct {
		name     string
		value    interface{ String() string }
		expected string
	}{
		{name: "review unknown", value: status.ReviewState_Unknown, expected: "Unknown"},
		{name: "review approved", value: status.ReviewState_Approved, expected: "Approved"},
		{name: "review required", value: status.ReviewState_RequiresReview, expected: "Review required"},
		{name: "changes requested", value: status.ReviewState_ChangesRequested, expected: "Changes requested"},
		{name: "no review required", value: status.ReviewState_NoReviewRequired, expected: "No review required"},
		{name: "undefined review", value: status.ReviewState(99), expected: "Undefined(99)"},
		{name: "mergeable unknown", value: status.MergeableState_Unknown, expected: "Unknown"},
		{name: "conflict", value: status.MergeableState_Conflict, expected: "Conflict"},
		{name: "mergeable", value: status.MergeableState_Mergeable, expected: "Can be merged"},
		{name: "undefined mergeable", value: status.MergeableState(42), expected: "Undefined(42)"},
		{name: "no checks", value: status.ChecksState_None, expected: "No checks"},
		{name: "checks passing", value: status.ChecksState_Passing, expected: "Checks passing"},
		{name: "checks failing", value: status.ChecksState_Failing, expected: "Checks failing"},
		{name: "checks running", value: status.ChecksState_Pending, expected: "Checks running"},
		{name: "checks unknown", value: status.ChecksState_Unknown, expected: "Checks unknown"},
		{name: "undefined checks", value: status.ChecksState(7), expected: "Undefined(7)"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.value.String(); got != tc.expected {
				t.Fatalf("got %q, expected %q", got, tc.expected)
			}
		})
	}
}
