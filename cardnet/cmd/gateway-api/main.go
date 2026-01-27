package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	analyticspb "cardnet/internal/api/analytics"
	authpb "cardnet/internal/api/grpc"
)

type Gateway struct {
	authClient      authpb.AuthGatewayClient
	analyticsClient analyticspb.AnalyticsServiceClient
}

func (g *Gateway) CORSMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers for all responses
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		// Note: Browsers normalize custom header names to lowercase for CORS preflight checks
		// Use lowercase for custom headers (x-correlation-id) to match browser behavior
		w.Header().Set("Access-Control-Allow-Headers", "accept, content-type, content-length, accept-encoding, x-csrf-token, authorization, x-correlation-id")
		w.Header().Set("Access-Control-Max-Age", "3600")

		// Handle preflight OPTIONS request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next(w, r)
	}
}

func (g *Gateway) handleAuthorize(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	var req struct {
		MerchantID string `json:"merchant_id"`
		CardToken  string `json:"card_token"`
		Amount     int64  `json:"amount"`
		Currency   string `json:"currency"`
		MCC        string `json:"mcc"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	resp, err := g.authClient.Authorize(context.Background(), &authpb.AuthRequest{
		RequestId:  uuid.New().String(),
		MerchantId: req.MerchantID,
		CardToken:  req.CardToken,
		Amount:     req.Amount,
		Currency:   req.Currency,
		Mcc:        req.MCC,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (g *Gateway) handleAnalyticsOverview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	resp, err := g.analyticsClient.GetNetworkStats(context.Background(), &analyticspb.GetNetworkStatsRequest{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (g *Gateway) handleMerchantStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	merchantID := r.URL.Query().Get("merchant_id")
	if merchantID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "merchant_id required"})
		return
	}

	resp, err := g.analyticsClient.GetMerchantStats(context.Background(), &analyticspb.GetMerchantStatsRequest{
		MerchantId: merchantID,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (g *Gateway) handleBinStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	bin := r.URL.Query().Get("bin")
	if bin == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "bin required"})
		return
	}

	resp, err := g.analyticsClient.GetBinStats(context.Background(), &analyticspb.GetBinStatsRequest{
		Bin: bin,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (g *Gateway) handleRuleEffectiveness(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	ruleID := r.URL.Query().Get("rule_id")
	if ruleID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "rule_id required"})
		return
	}

	resp, err := g.analyticsClient.GetRuleEffectiveness(context.Background(), &analyticspb.GetRuleEffectivenessRequest{
		RuleId: ruleID,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// handleAnalyticsHistory is a placeholder for a future network-wide transaction
// history API.
func (g *Gateway) handleAnalyticsHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	merchantID := r.URL.Query().Get("merchant_id")
	status := r.URL.Query().Get("status")

	resp, err := g.analyticsClient.GetNetworkHistory(context.Background(), &analyticspb.GetNetworkHistoryRequest{
		MerchantId: merchantID,
		Status:     status,
		Limit:      200,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// Shape response for frontend: match NetworkHistoryItem in analytics_frontend
	type httpItem struct {
		AuthID     string `json:"auth_id"`
		MerchantID string `json:"merchant_id"`
		Amount     int64  `json:"amount"`
		Currency   string `json:"currency"`
		Status     string `json:"status"`
		Reason     string `json:"reason,omitempty"`
		CreatedAt  string `json:"created_at"`
	}

	items := make([]httpItem, 0, len(resp.Items))
	for _, it := range resp.Items {
		items = append(items, httpItem{
			AuthID:     it.GetAuthId(),
			MerchantID: it.GetMerchantId(),
			// Convert dollars back to integer minor units for UI (cents)
			Amount:    int64(it.GetAmount() * 100),
			Currency:  it.GetCurrency(),
			Status:    it.GetStatus(),
			Reason:    it.GetReason(),
			CreatedAt: time.Unix(it.GetCreatedAt(), 0).UTC().Format(time.RFC3339),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(items)
}

func main() {
	log.Println("Starting Gateway API...")

	authAddr := os.Getenv("AUTH_GATEWAY_ADDR")
	if authAddr == "" {
		authAddr = "localhost:50051"
	}

	analyticsAddr := os.Getenv("ANALYTICS_SERVICE_ADDR")
	if analyticsAddr == "" {
		analyticsAddr = "localhost:60057"
	}

	log.Printf("Connecting to Auth Gateway at %s", authAddr)
	authConn, err := grpc.NewClient(authAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to auth gateway: %v", err)
	}
	// defer authConn.Close()

	log.Printf("Connecting to Analytics Service at %s", analyticsAddr)
	analyticsConn, err := grpc.NewClient(analyticsAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to analytics service: %v", err)
	}
	// defer analyticsConn.Close()

	gateway := &Gateway{
		authClient:      authpb.NewAuthGatewayClient(authConn),
		analyticsClient: analyticspb.NewAnalyticsServiceClient(analyticsConn),
	}

	http.HandleFunc("/api/auth/authorize", gateway.CORSMiddleware(gateway.handleAuthorize))
	http.HandleFunc("/api/analytics/overview", gateway.CORSMiddleware(gateway.handleAnalyticsOverview))
	http.HandleFunc("/api/analytics/merchant", gateway.CORSMiddleware(gateway.handleMerchantStats))
	http.HandleFunc("/api/analytics/bin", gateway.CORSMiddleware(gateway.handleBinStats))
	http.HandleFunc("/api/analytics/rules", gateway.CORSMiddleware(gateway.handleRuleEffectiveness))
	http.HandleFunc("/api/analytics/history", gateway.CORSMiddleware(gateway.handleAnalyticsHistory))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Gateway listening on :%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
