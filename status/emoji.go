package status

func Emoji(input PRInfo) []string {
	emoji := make([]string, 0)

	if ShowGreenIcon(input) {
		emoji = append(emoji, "🟢")
	} else {
		emoji = append(emoji, "🔴")
	}

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

	switch input.Checks {
	case ChecksState_None, ChecksState_Passing:
		break
	case ChecksState_Failing:
		emoji = append(emoji, "❌")
	case ChecksState_Pending:
		emoji = append(emoji, "⏳")
	case ChecksState_Unknown:
		emoji = append(emoji, "❔")
	}

	return emoji
}
