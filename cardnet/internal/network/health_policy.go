package network

const (
	MinIssuerHealth = 60
)

func IsHealthy(score int) bool {
	return score >= MinIssuerHealth
}
