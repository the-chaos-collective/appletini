package status

func ShowGreenIcon(input PRInfo) bool {
	return (input.Mergeable == MergeableState_Mergeable || input.Mergeable == MergeableState_Unknown) &&
		(input.Review == ReviewState_Approved || input.Review == ReviewState_NoReviewRequired)
}

func ShowRedIcon(input PRInfo) bool {
	return input.Review == ReviewState_ChangesRequested || input.Mergeable == MergeableState_Conflict
}
