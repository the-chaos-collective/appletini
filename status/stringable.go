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
