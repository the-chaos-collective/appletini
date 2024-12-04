package status

func ShowGreenIcon(input PRInfo) bool {
	return false
}

func ShowRedIcon(input PRInfo) bool {
	return input.Review == ReviewState_RequiresReview
}
