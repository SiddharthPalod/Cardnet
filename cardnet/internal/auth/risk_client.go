package auth

import (
	"bytes"
	"cardnet/internal/risk/models"
	"encoding/json"
	"net/http"
)

func CallRiskEngine(tx models.Transaction, ctx models.RiskContext) (models.RiskResult, error) {

	payload := map[string]interface{}{
		"tx":  tx,
		"ctx": ctx,
	}

	body, _ := json.Marshal(payload)

	resp, err := http.Post(
		"http://risk-engine:8082/evaluate",
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return models.RiskResult{}, err
	}
	defer resp.Body.Close()

	var result models.RiskResult
	json.NewDecoder(resp.Body).Decode(&result)

	return result, nil
}
