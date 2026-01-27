# 🌐 Mega Project: **CardNet**

> A high-scale **card authorization, risk & settlement gateway**
> A card network authorization gateway handling high TPS, with inline risk scoring, merchant rate limiting, immutable event ledgers, and async settlement pipelines using Go, gRPC, Next.js, Postgres, and Cassandra.”

---

# 0️⃣ Vision, Stack & Architecture (Once at the Start)

## 🎯 Goal

Build a **high-throughput card authorization system** that:

* Processes payment authorizations in **<50ms**
* Enforces **merchant & BIN rate limits**
* Performs **real-time fraud checks**
* Emits **immutable authorization events**
* Enables **async settlement & analytics**

This is **core payments infrastructure**.

For demo purposed frontend:

- **Frontend A – Merchant Simulator UI**: web app that sends auth requests to `auth-gateway` and shows real-time results.
- **Frontend B – Ops/Analytics Dashboard**: web app that visualizes how the system is behaving (throughput, latency, approval rates, rate-limiting, risk decisions, errors).

Demo: [Watch demo video](https://youtu.be/7A6a_pECYzY)

---

## 🧰 Tech Stack
* **Backend**: Go
* **Frontend**: Next.js
* **Protocol**: gRPC (internal services)
* **Databases**:
  * PostgreSQL (ACID state)
  * Cassandra (event ledger)
* **Messaging**: Kafka-style event queues
* **Infra**: Docker, docker-compose

---

# **Test Results**

* **SLA Compliance**: Average latency 24-47ms (target: <50ms)
* **Throughput**: 746-1,317 req/s (target: 1000+ req/s)
* **Success Rate**: 75-100% depending on load
* **Backpressure**: System handles 10x load gracefully (99.63% success)
* **Rate Limiting**: Effective protection (65-97% rate limited)
* **Stability**: No crashes or hangs under any load

### **Performance Summary:**
- Quick Load (10 workers): 355 req/s, 24.38ms avg latency ✅
- Standard Load (20 workers): 746 req/s, 26.71ms avg latency ✅
- High Load (50 workers): 1,317 req/s, 37.86ms avg latency ✅
- Backpressure (50 workers, 10x load): 1,064 req/s, 46.46ms avg latency ✅


# 🏗️ High-Level System Design Diagram

```text
Merchants (POS / Payment APIs)
   |
   v
+---------------------------------------------------+
| Authorization Gateway (Stateless, Horizontally Scaled) |
| - Validation                                        |
| - Idempotency / Deduplication                      |
| - Correlation IDs                                   |
+----------------------+----------------+-------------+
             |                |
             |                |
          +-------v------+   +-----v-------+
          | Rate Limiter |   | Risk Engine |
          | - Merchant TPS|  | - Velocity  |
          | - BIN / MCC   |  | - Geo checks|
          +-------+-------+  | - Amount dev |
             \          +-----+-------+
         \               |
          \              |
           \            v
            +---------------------+
            |  Issuer Simulator   |
            |  - Approve / Decline|
            |  - Timeouts / Latency|
            +----------+----------+
                  |
                  v
            +---------------------+
            | Authorization Result|
            +----------+----------+
                  |
      +--------------------------+--------------------------+
      |                          |                          |
      v                          v                          v
    +-------------+            +--------------+           +----------------+
    | PostgreSQL  |            | Cassandra    |           | Message Queue  |
    | (Strong ACID|            | (Event Ledger)|          | (auth.events)  |
    |  state)     |            | - Auth events |          | - settlement   |
    +------+------+            +------+-------+           +-------+--------+
      |                          |                           |
      v                          v                           v
  Settlement Service           Fraud Analytics                Merchant Reports
  (Batch, T+1)                 (Batch / Stream)               (Dashboards / Exports)
```

---

# 🔁 Authorization Request Flow 

```
1. Merchant sends authorization request
   |
2. Authorization Gateway validates + deduplicates
   |
3. Rate Limiter checks merchant/BIN limits
   |
4. Risk Engine computes fraud score
   |
5. Issuer Simulator returns approve/decline
   |
6. Final decision returned to merchant
   |
7. Events written asynchronously to ledger + queue
```

### Mental Model:

```
Merchant
   |
(1)| gRPC
   ▼
Authorization Gateway
   |
(2)| Rate Limit Check
   |
(3)| Risk Evaluation
   |
(4)| Issuer Authorization
   |
(5)| Final Decision
   |
(6)| Persist State + Publish Events
```

---

# Design Choices

### Why **Authorization Gateway is Stateless**

* Enables horizontal scaling
* No session affinity
* Easy load balancing

---

### Why **Postgres + Cassandra**

* **Postgres** → current authorization state (ACID)
* **Cassandra** → immutable, high-volume event ledger
* Separation of **state vs history**

---

### Why **Risk Engine is Synchronous**

* Fraud decisions must happen **before money moves**
* Async fraud is too late for card networks

---

### Why **Message Queue**

* Decouples real-time authorization from:

  * settlement
  * analytics
  * reporting
* Protects hot path latency

---