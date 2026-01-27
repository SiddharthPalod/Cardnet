### High-level goal

- **Frontend A – Merchant Simulator UI**: web app that sends auth requests to `auth-gateway` and shows real-time results.
- **Frontend B – Ops/Analytics Dashboard**: web app that visualizes how the system is behaving (throughput, latency, approval rates, rate-limiting, risk decisions, errors).

---

## Phase 0 – Constraints & Architecture decisions

- **Frontend tech**:  
  - **Option**: React + TypeScript (Vite/Next.js) for both UIs to share components.
- **Browser → Go backend**:
  - Browsers can’t directly call gRPC; plan to add a **REST/HTTP façade** in front of the gRPC services:
    - Either:
      - **`gateway-api` service (Go)** exposing JSON/REST and internally calling `auth-gateway` and `analytics-service` via gRPC, or
      - Use **grpc-gateway** to auto-generate REST endpoints from existing `.proto`.
  - All UI traffic goes through this HTTP layer.
- **Analytics data source**:
  - Use existing **Postgres/Cassandra/queues** plus an **`analytics-service` gRPC API** (already in your plan) to provide:
    - TPS, p95 latency, error/decline breakdown, rate-limit counts, rule hits, etc.

---

## Phase 1 – Stabilize & document current backend

- **Inventory services**:
  - Confirm running services: `auth-gateway`, `rate-limiter`, `risk-engine`, `issuer-simulator`, `event-ledger`, `settlement-service`, `analytics-service`.
- **Define external contracts**:
  - For **Merchant UI**:
    - `POST /api/auth/authorize` → `{ merchant_id, card_token, amount, currency, ... }` → `{ status, reason, auth_id, ... }`
  - For **Analytics UI**:
    - `GET /api/analytics/overview`
    - `GET /api/analytics/merchants/:id`
    - `GET /api/analytics/timeseries?...`
- **Create example flows**:
  - A few **canonical test merchants / cards / scenarios** to be used from the UI (e.g. “high-risk merchant”, “rate-limited BIN”, “geo mismatch”).

---

## Phase 2 – HTTP/REST façade for the auth path

- **Design `gateway-api` (or grpc-gateway) layer**:
  - **Responsibilities**:
    - Terminate HTTP/JSON.
    - Validate + map requests to `auth-gateway.Authorize`.
    - Translate gRPC responses/errors to friendly HTTP status codes.
    - Add correlation IDs to responses (to later cross-link in Analytics UI).
- **Endpoints for Merchant UI**:
  - **`POST /api/auth/authorize`**
    - Request: merchant/card/amount info.
    - Response: auth decision + metadata (risk score, decision reason, issuer latency).
  - **`GET /api/auth/sample-config`**
    - Optional: return sample merchants/cards for demo.
- **Security**:
  - For now, keep simple:
    - No auth or a static API key in headers for local testing.
    - Plan for future: JWT-based merchant auth.

**Deliverable**: You can drive the existing gRPC hot path from curl/Postman using plain HTTP/JSON.

---

## Phase 3 – Frontend A: Merchant Simulator UI

### 3.1 UX & pages

- **Page 1 – “Send Authorization”**
  - Form fields:
    - Merchant selector (dropdown),
    - Card token selector (or free text),
    - Amount, currency, MCC,
    - Optional: geo/location.
  - “Send Request” button.
  - Display:
    - **Result panel**: APPROVE / SOFT_DECLINE / HARD_DECLINE, reason, auth_id.
    - Raw JSON of request/response (collapsible for debugging).

- **Page 2 – “Live Stream / History”**
  - Table of last N auths:
    - Time, merchant, amount, status, decision reason, issuer latency.
  - Simple filters (merchant, status).
  - Auto-refresh or WebSocket to stream events (optional enhancement).

### 3.2 Implementation tasks

- **Setup app skeleton**:
  - `merchant-ui/` with React+TS.
  - Basic routing: `/auth`, `/history`.
- **Integrate with backend**:
  - API client module calling `POST /api/auth/authorize`.
  - Error handling and loading states.
- **Developer tooling**:
  - `.env` for `VITE_API_BASE_URL` (or similar).
- **Stretch**:
  - WebSocket or SSE endpoint from backend to push new auth results to the UI in real time.

**Deliverable**: A user can sit on the Merchant UI, send real-time auth requests, and see responses immediately.

---

## Phase 4 – Analytics data API

### 4.1 Define analytics metrics

Backed by `analytics-service` and/or DBs, provide:

- **Overview metrics** (for dashboard header):
  - Current/smoothed **TPS**.
  - p50/p95 **latency**.
  - **Approval rate** (%).
  - **Decline reasons breakdown** (risk vs issuer vs rate-limit).
- **Time-series**:
  - TPS over time (per minute).
  - Latency percentiles over time.
  - Approvals/declines over time.
- **Per-dimension views**:
  - By merchant: TPS, approval rate, risk rule hits.
  - By BIN: TPS, decline rate.

### 4.2 Analytics HTTP API

Exposed by the same `gateway-api` or by a dedicated `analytics-api`:

- `GET /api/analytics/overview`
- `GET /api/analytics/timeseries?metric=tps&window=15m`
- `GET /api/analytics/merchants/:merchantId`
- `GET /api/analytics/rules` (optional: rule hit counts)
- `GET /api/analytics/rate-limits` (optional: counters from rate-limiter)

**Implementation**:
- Add or reuse **gRPC methods** in `analytics-service`:
  - `GetOverview`, `GetTimeSeries`, `GetMerchantMetrics`, etc.
- Implement the **HTTP → gRPC** adapters mirroring those methods.

---

## Phase 5 – Frontend B: Ops/Analytics Dashboard

### 5.1 UX & pages

- **Page 1 – Global Dashboard**
  - KPI cards: TPS, p95 latency, approval %, error rate.
  - Graphs:
    - TPS over last 15m/1h.
    - p95 latency over last 15m/1h.
    - Stacked bar for approvals/declines by reason.
- **Page 2 – Merchant Detail**
  - Selector/search for merchant.
  - Per-merchant stats: TPS, approval rate, most common decline reasons.
  - Time-series of that merchant’s activity.
- **Page 3 – System Health (optional)**
  - Rate-limiter stats (how many requests throttled).
  - Risk rule hit frequency.
  - Issuer timeouts count.

### 5.2 Implementation tasks

- **Setup `analytics-ui/`**:
  - Same stack as Merchant UI to share components.
- **Charts**:
  - Pick a charting library (e.g. Recharts, Chart.js, ECharts).
  - Implement reusable `TimeSeriesChart`, `KpiCard`.
- **Integrate APIs**:
  - API client for `GET /api/analytics/...`.
  - Poll every 5–10s or use WebSockets/SSE for live updates.

**Deliverable**: An operator can see how the **first UI’s traffic** impacts the system in real time.

---

## Phase 6 – Wiring, correlation & storytelling

- **Correlation IDs**:
  - Include a `correlation_id` or `trace_id` in:
    - Merchant UI response.
    - Analytics and logs.
  - Allow Analytics UI to search by that ID to show the path through rate-limiter → risk-engine → issuer.
- **Cross-link UIs**:
  - In Merchant UI, for each auth result, add a link “View in Analytics” that opens Analytics UI with filters for that merchant/time.
- **Demo scenarios**:
  - Preconfigure:
    - Normal merchant (mostly approved).
    - Risky merchant (high fraud declines).
    - Overloaded merchant/BIN (rate-limited).

---

## Phase 7 – Hardening & polish

- **Non-functional**:
  - Basic auth for Analytics UI (even simple password) to mimic real ops tools.
  - Rate-limiting on HTTP façade, if desired.
- **DX**:
  - `docker-compose` profiles:
    - `core-backend`,
    - `merchant-ui`,
    - `analytics-ui`.
  - One-command dev startup script.

---