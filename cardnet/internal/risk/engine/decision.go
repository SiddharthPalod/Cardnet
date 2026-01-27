package engine

func mapDecision(score int) string {
	switch {
	case score <= 30:
		return "APPROVE"
	case score <= 60:
		return "SOFT_DECLINE"
	default:
		return "HARD_DECLINE"
	}
}
