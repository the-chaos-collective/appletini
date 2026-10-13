package status

// ReviewState
type ReviewState uint

const (
	ReviewState_Unknown ReviewState = iota
	ReviewState_Approved
	ReviewState_RequiresReview
	ReviewState_ChangesRequested
	ReviewState_NoReviewRequired
)

// MergeableState
type MergeableState uint

const (
	MergeableState_Unknown MergeableState = iota
	MergeableState_Conflict
	MergeableState_Mergeable
)

// ChecksState is the combined CI result (all checks and commit statuses) on
// the PR's latest commit. The zero value means the PR has no checks.
type ChecksState uint

const (
	ChecksState_None ChecksState = iota
	ChecksState_Passing
	ChecksState_Failing
	ChecksState_Pending
	ChecksState_Unknown
)

type PRInfo struct {
	Review    ReviewState
	Mergeable MergeableState
	Checks    ChecksState
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
