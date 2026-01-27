package feedback

import "time"

type ThresholdUpdate struct {
	RuleName       string    `json:"rule_name"`
	EntityType     string    `json:"entity_type"` // merchant | bin | global
	EntityID       string    `json:"entity_id"`
	SuggestedValue int       `json:"suggested_value"`
	Confidence     float64   `json:"confidence"`
	ModelVersion   string    `json:"model_version"`
	CreatedAt      time.Time `json:"created_at"`
}
