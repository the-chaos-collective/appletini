package status

func Emoji(input PRInfo) []string {
	if input.Mergeable == MergeableState_Mergeable && (input.Review == ReviewState_Approved || input.Review == ReviewState_NoReviewRequired) {
		return []string{"🟢"}
	}
	return []string{"🔴", "🚧"}
}
