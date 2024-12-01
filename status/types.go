package status

// ReviewState
type ReviewState int

const (
	ReviewState_Unknown ReviewState = iota
	ReviewState_Approved
	ReviewState_RequiresReview
	ReviewState_ChangesRequested
	ReviewState_NoReviewRequired
)

// MergeableState
type MergeableState int

const (
	MergeableState_Unknown MergeableState = iota
	MergeableState_Conflict
	MergeableState_Mergeable
)

type PRInfo struct {
	Review    ReviewState
	Mergeable MergeableState
}

type Status struct {
	ShowGreenIcon bool
	ShowRedIcon   bool
	Emoji         []string
	Message       string
}

func (input PRInfo) Classify() Status {
	return Status{
		ShowGreenIcon: ShowGreenIcon(input),
		ShowRedIcon:   ShowRedIcon(input),
		Emoji:         Emoji(input),
		Message:       Message(input),
	}
}
