package auth

import (
	"context"
	"encoding/json"
	"time"

	ratepb "cardnet/internal/api/grpc"
	issuerpb "cardnet/internal/api/issuer"
	ledgerpb "cardnet/internal/api/ledger"
	riskpb "cardnet/internal/api/risk"
	"cardnet/internal/gates"
	"cardnet/internal/issuer/routing"
	"cardnet/internal/middleware"
	"cardnet/internal/network"
	"cardnet/internal/transaction"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"
)

type Service struct {
	repo            Repository
	txStateRepo     transaction.Repository
	asyncStateWriter *transaction.AsyncStateWriter
	idempotencyCache *IdempotencyCache
	gates           *gates.AtomicGates
	healthTracker   *IssuerHealthTracker
	rateClient      ratepb.RateLimiterClient
	riskClient      riskpb.RiskEngineClient
	issuerClient    issuerpb.IssuerServiceClient
	ledgerClient    ledgerpb.EventLedgerServiceClient
	networkMetrics  *network.Metrics
	rateCB          *middleware.CircuitBreaker
	riskCB          *middleware.CircuitBreaker
	issuerCB        *middleware.CircuitBreaker
}

func NewService(
	repo Repository,
	txStateRepo transaction.Repository,
	gatesConfig gates.Config,
	networkMetrics *network.Metrics,
	rateClient ratepb.RateLimiterClient,
	riskClient riskpb.RiskEngineClient,
	issuerClient issuerpb.IssuerServiceClient,
	ledgerClient ledgerpb.EventLedgerServiceClient,
) *Service {
	healthTracker := NewIssuerHealthTracker()
	
	// Fast path: In-memory idempotency cache (O(1) lookups)
	idempotencyCache := NewIdempotencyCache(5 * time.Minute)
	
	// Async state writer (durable path off critical path)
	// Increase buffer size to handle high load (10000 to prevent queue full errors)
	asyncStateWriter := transaction.NewAsyncWriter(txStateRepo, 10000)
	
	// Atomic gates (lock-free, updated periodically)
	gatesInstance := gates.NewAtomicGates(gatesConfig, networkMetrics, healthTracker)
	
	return &Service{
		repo:             repo,
		txStateRepo:      txStateRepo,
		asyncStateWriter: asyncStateWriter,
		idempotencyCache: idempotencyCache,
		gates:             gatesInstance,
		healthTracker:     healthTracker,
		rateClient:     rateClient,
		riskClient:     riskClient,
		issuerClient:   issuerClient,
		ledgerClient:   ledgerClient,
		networkMetrics: networkMetrics,
		rateCB:         middleware.NewCircuitBreaker(5, 5*time.Second),
		riskCB:         middleware.NewCircuitBreaker(5, 5*time.Second),
		issuerCB:       middleware.NewCircuitBreaker(5, 5*time.Second),
	}
}

// Payload structures
type AuthRequestPayload struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
}

type RiskDecisionPayload struct {
	Score   int32    `json:"score"`
	Reasons []string `json:"reasons"`
}

type IssuerResponsePayload struct {
	Code string `json:"code"`
}

type FinalOutcomePayload struct {
	Status   string `json:"status"`
	Reason   string `json:"reason"`
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
}

func (s *Service) Authorize(
	ctx context.Context,
	requestID string,
	merchantID string,
	cardToken string,
	amount int64,
	currency string,
	mcc string,
) (string, string, error) {
	// FAST PATH: O(1) in-memory idempotency check (no DB call)
	existingAuthID, existingStatus, found := s.idempotencyCache.Get(requestID)
	if found {
		// Request already processed - return cached result (idempotent)
		return existingAuthID, existingStatus, nil
	}

	// FALLBACK: DB lookup if not in cache (non-blocking)
	existingAuthID, existingStatus, err := s.repo.GetByRequestID(ctx, requestID)
	if err == nil {
		// Found in DB - cache it and return
		s.idempotencyCache.Set(requestID, existingAuthID, existingStatus)
		return existingAuthID, existingStatus, nil
	}
	// If DB lookup fails, proceed optimistically (don't block auth path)

	// Extract BIN from card token (handles dashes, spaces, etc.)
	bin := routing.ExtractBIN(cardToken)

	// Generate AuthID early for tracking (idempotency key for state transitions)
	authID := uuid.New().String()

	// ASYNC: Record INITIATED state (non-blocking, queued for background write)
	s.asyncStateWriter.RecordStateAsync(authID, transaction.INITIATED)

	// 1. EVENT: AUTH_REQUEST
	authPayload, _ := json.Marshal(AuthRequestPayload{
		Amount:   amount,
		Currency: currency,
	})
	s.publishEvent(
		authID, merchantID, cardToken,
		ledgerpb.EventType_AUTH_REQUEST,
		ledgerpb.Decision_DECISION_UNSPECIFIED,
		authPayload,
	)

	rateCtx, cancelRate := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancelRate()

	var rlResp *ratepb.RateLimitResponse
	err = s.rateCB.Execute(func() error {
		var callErr error
		rlResp, callErr = s.rateClient.CheckLimit(rateCtx, &ratepb.RateLimitRequest{
			MerchantId: merchantID,
			Bin:        bin,
			Mcc:        mcc,
		})
		return callErr
	})
	if err != nil {
		return "", "", err
	}

	if !rlResp.Allowed {
		return "", "RATE_LIMITED", nil
	}

	riskCtx, cancelRisk := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancelRisk()

	var riskResp *riskpb.RiskResponse
	err = s.riskCB.Execute(func() error {
		var callErr error
		riskResp, callErr = s.riskClient.Evaluate(riskCtx, &riskpb.RiskRequest{
			Transaction: &riskpb.Transaction{
				RequestId:  requestID,
				MerchantId: merchantID,
				Bin:        bin,
				Amount:     amount,
				Currency:   currency,
				Mcc:        mcc,
			},
		})
		return callErr
	})
	if err != nil {
		return "", "RISK_UNAVAILABLE", nil
	}

	// ASYNC: Save risk audit (non-blocking, could be queued)
	_ = s.repo.SaveRiskAudit(
		ctx,
		requestID,
		int(riskResp.RiskScore),
		riskResp.Reasons,
	)

	// ASYNC: Record RISK_EVALUATED state (non-blocking, queued for background write)
	s.asyncStateWriter.RecordStateAsync(authID, transaction.RISK_EVALUATED)

	// 2. EVENT: RISK_DECISION
	riskDecision := ledgerpb.Decision_DECISION_UNSPECIFIED
	switch riskResp.Decision {
	case riskpb.Decision_APPROVE:
		riskDecision = ledgerpb.Decision_APPROVED
	case riskpb.Decision_SOFT_DECLINE:
		riskDecision = ledgerpb.Decision_REVIEW
	case riskpb.Decision_HARD_DECLINE:
		riskDecision = ledgerpb.Decision_DECLINED
	}

	riskPayload, _ := json.Marshal(RiskDecisionPayload{
		Score:   riskResp.RiskScore,
		Reasons: riskResp.Reasons,
	})
	s.publishEvent(
		authID, merchantID, cardToken,
		ledgerpb.EventType_RISK_DECISION,
		riskDecision,
		riskPayload,
	)

	switch riskResp.Decision {
	case riskpb.Decision_APPROVE:
		// continue

	case riskpb.Decision_SOFT_DECLINE:
		return "", "CHALLENGE_REQUIRED", nil

	case riskpb.Decision_HARD_DECLINE:
		return "", "DECLINED_FRAUD", nil

	default:
		return "", "RISK_UNKNOWN", nil
	}

	// =====================================================
	// 5️⃣ ISSUER AUTHORIZATION
	// =====================================================
	
	// Resolve issuer to check health gate
	// Note: If BIN resolution fails here, we still try the issuer call
	// The card-network service will handle BIN resolution and return appropriate errors
	profile, ok := routing.ResolveIssuer(cardToken)
	issuerKey := ""
	if ok {
		issuerKey = profile.Name
		// Check issuer health gate - skip if unhealthy
		if s.gates.ShouldSkipIssuer(issuerKey) {
			// Fast fail - issuer is unhealthy, skip issuer call
			// For now, we'll decline if issuer is skipped
			return "", "ISSUER_UNAVAILABLE", nil
		}
	}
	// If BIN resolution failed, issuerKey will be empty but we still proceed
	// The card-network will handle the BIN resolution and return an error if needed

	// ASYNC: Record ISSUER_REQUESTED state (non-blocking, queued for background write)
	s.asyncStateWriter.RecordStateAsync(authID, transaction.ISSUER_REQUESTED)

	issuerCtx, cancelIssuer := context.WithTimeout(ctx, 2*time.Second) // Higher timeout for banks
	defer cancelIssuer()

	var issuerResp *issuerpb.IssuerAuthResponse

	err = s.issuerCB.Execute(func() error {
		var callErr error
		issuerResp, callErr = s.issuerClient.Authorize(issuerCtx, &issuerpb.IssuerAuthRequest{
			AuthId:     authID,
			CardNumber: cardToken,
			Amount:     amount,
			Currency:   currency,
		})
		return callErr
	})

	// Track issuer health metrics
	if err != nil {
		if err == context.DeadlineExceeded {
			s.healthTracker.RecordTimeout(issuerKey)
		} else {
			s.healthTracker.RecordFailure(issuerKey)
		}
	} else {
		s.healthTracker.RecordSuccess(issuerKey)
	}

	// ASYNC: Record ISSUER_RESPONDED state (non-blocking, queued for background write)
	s.asyncStateWriter.RecordStateAsync(authID, transaction.ISSUER_RESPONDED)

	// 3. EVENT: ISSUER_RESPONSE
	issuerDecision := ledgerpb.Decision_DECLINED
	issuerCode := "error"
	if err == nil {
		if issuerResp.Approved {
			issuerDecision = ledgerpb.Decision_APPROVED
		}
		issuerCode = issuerResp.DeclineCode
	}

	issuerPayload, _ := json.Marshal(IssuerResponsePayload{
		Code: issuerCode,
	})
	s.publishEvent(
		authID, merchantID, cardToken,
		ledgerpb.EventType_ISSUER_RESPONSE,
		issuerDecision,
		issuerPayload,
	)

	status := "APPROVED"
	if err != nil {
		// Check if it's an "unknown BIN" error from card-network
		if st, ok := grpcstatus.FromError(err); ok {
			if st.Code() == codes.InvalidArgument && st.Message() == "unknown BIN" {
				return "", "UNKNOWN_BIN", nil
			}
		}
		// Differentiate timeout vs failure if needed, for now centralized:
		status = "ERROR_ISSUER"
		// If issuer is down (CB open or timeout), we might decline
		return "", "ISSUER_UNAVAILABLE", nil
	}

	if !issuerResp.Approved {
		status = "DECLINED"
		return "", issuerResp.DeclineCode, nil
	}

	// =====================================================
	// 6️⃣ PERSISTENCE
	// =====================================================
	// NON-BLOCKING: If DB save fails, log but don't block response
	// Durability != latency - we've already processed the authorization
	err = s.repo.SaveAuthorization(
		ctx,
		authID,
		requestID,
		merchantID,
		cardToken,
		amount,
		status,
	)

	if err != nil {
		if err == ErrDuplicateRequest {
			// Duplicate request - return existing result (non-blocking lookup)
			existingAuthID, existingStatus, lookupErr := s.repo.GetByRequestID(ctx, requestID)
			if lookupErr == nil {
				return existingAuthID, existingStatus, nil
			}
			// If lookup also fails, proceed optimistically
		}
		// Log but don't block - authorization was already processed
		// In production, use structured logging
		_ = err // Silently proceed - durability != latency
	}

	// SYNC: Record FINALIZED state (only final state is synchronous for correctness)
	_ = s.asyncStateWriter.RecordStateSync(ctx, authID, transaction.FINALIZED)
	
	// Cache result for idempotency (fast path for future requests)
	s.idempotencyCache.Set(requestID, authID, status)
	
	finalPayload, _ := json.Marshal(FinalOutcomePayload{
		Status:   status,
		Reason:   "authorized",
		Amount:   amount,
		Currency: currency,
	})
	s.publishEvent(authID, merchantID, cardToken, ledgerpb.EventType_FINAL_OUTCOME, ledgerpb.Decision_APPROVED, finalPayload)
	return authID, status, nil
}

func (s *Service) publishEvent(
	authID, merchantID, cardHash string,
	eventType ledgerpb.EventType,
	decision ledgerpb.Decision,
	payload []byte,
) {
	// Fire and forget - minimal latency impact
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, _ = s.ledgerClient.RecordEvent(ctx, &ledgerpb.AuthEvent{
			EventId:    uuid.New().String(),
			AuthId:     authID,
			MerchantId: merchantID,
			CardHash:   cardHash,
			EventType:  eventType,
			Decision:   decision,
			Payload:    payload,
			CreatedAt:  time.Now().Unix(),
		})
	}()
}
