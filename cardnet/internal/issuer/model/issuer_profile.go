package model

type ChaosConfig struct {
	ForceTimeout   bool
	ForceDecline   bool
	ExtraLatencyMs int
}

type IssuerProfile struct {
	Name        string
	LatencyMs   int
	DeclineRate float64
	TimeoutRate float64
	Chaos       ChaosConfig
}
