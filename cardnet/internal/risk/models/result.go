package models

type RiskResult struct {
	Decision  string
	RiskScore int
	Reasons   []string
}
