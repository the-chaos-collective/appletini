package status

func Emoji(input PRInfo) []string {
	if input.Mergeable == MergeableState_Mergeable && (input.Review == ReviewState_Approved || input.Review == ReviewState_NoReviewRequired) {
		return []string{"🟢"}
	}

	emoji := []string{"🔴"}

	switch input.Mergeable {
	case MergeableState_Mergeable:
		break
	case MergeableState_Conflict:
		emoji = append(emoji, "🚧")
	case MergeableState_Unknown:
		emoji = append(emoji, "❓")
	}

	switch input.Review {
	case ReviewState_Approved:
		break
	case ReviewState_NoReviewRequired:
		break
	case ReviewState_ChangesRequested:
		emoji = append(emoji, "📝")
	case ReviewState_RequiresReview:
		emoji = append(emoji, "👁‍🗨")
	case ReviewState_Unknown:
		emoji = append(emoji, "❔")
	}

	return emoji
}
