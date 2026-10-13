package status

import "fmt"

func (state ReviewState) String() string {
	switch state {
	case ReviewState_Unknown:
		return "Unknown"
	case ReviewState_Approved:
		return "Approved"
	case ReviewState_RequiresReview:
		return "Review required"
	case ReviewState_ChangesRequested:
		return "Changes requested"
	case ReviewState_NoReviewRequired:
		return "No review required"
	}
	return fmt.Sprintf("Undefined(%d)", state)
}

func (state MergeableState) String() string {
	switch state {
	case MergeableState_Unknown:
		return "Unknown"
	case MergeableState_Conflict:
		return "Conflict"
	case MergeableState_Mergeable:
		return "Can be merged"
	}
	return fmt.Sprintf("Undefined(%d)", state)
}

func (state ChecksState) String() string {
	switch state {
	case ChecksState_None:
		return "No checks"
	case ChecksState_Passing:
		return "Checks passing"
	case ChecksState_Failing:
		return "Checks failing"
	case ChecksState_Pending:
		return "Checks running"
	case ChecksState_Unknown:
		return "Checks unknown"
	}
	return fmt.Sprintf("Undefined(%d)", state)
}
