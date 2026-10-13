package status

// ShowGreenIcon is true when the PR is ready to merge: approved (or no review
// required), not conflicting, and CI passing or absent. Running CI is not ready.
func ShowGreenIcon(input PRInfo) bool {
	return (input.Mergeable == MergeableState_Mergeable || input.Mergeable == MergeableState_Unknown) &&
		(input.Review == ReviewState_Approved || input.Review == ReviewState_NoReviewRequired) &&
		(input.Checks == ChecksState_Passing || input.Checks == ChecksState_None)
}

// ShowRedIcon is true when the PR needs the author's attention.
func ShowRedIcon(input PRInfo) bool {
	return input.Review == ReviewState_ChangesRequested ||
		input.Mergeable == MergeableState_Conflict ||
		input.Checks == ChecksState_Failing
}
