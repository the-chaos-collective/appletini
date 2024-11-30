package status

// ReviewState
type ReviewState int

const (
	ReviewState_Approved ReviewState = iota
	ReviewState_RequiresReview
)

// MergeableState
type MergeableState int

const (
	MergeableState_Conflict MergeableState = iota
	MergeableState_Mergeable
)

type PRInfo struct {
	Review    ReviewState
	Mergeable MergeableState
}

type Status struct {
	ShowGreenIcon bool
	ShowRedIcon   bool
	Emoji         string
	Message       string
}

func Classify(input PRInfo) Status {
	return Status{
		ShowGreenIcon: ShowGreenIcon(input),
		ShowRedIcon:   ShowRedIcon(input),
		Emoji:         Emoji(input),
		Message:       Message(input),
	}
}
