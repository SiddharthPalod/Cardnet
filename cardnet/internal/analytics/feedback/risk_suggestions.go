package feedback

import "time"

type RiskSuggestion struct {
	EntityType     string    `json:"entity_type"` // merchant | bin | transaction
	EntityID       string    `json:"entity_id"`
	SuggestedScore float64   `json:"suggested_score"`
	Confidence     float64   `json:"confidence"` // 0..1
	ModelVersion   string    `json:"model_version"`
	Reason         string    `json:"reason"`
	CreatedAt      time.Time `json:"created_at"`
}
