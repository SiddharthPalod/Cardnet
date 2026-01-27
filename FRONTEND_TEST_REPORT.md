# Frontend Testing Report
**Date:** January 27, 2026  
**Tested URLs:**
- Analytics Frontend: `http://localhost:3001/`
- Merchant Frontend: `http://localhost:3000/`

---

## 1. Analytics Frontend (`http://localhost:3001/`)

### Overall Layout
Based on code analysis, the analytics frontend displays:
- **Navbar** at the top
- **Hero section** with title/description
- **KPI Cards Row** (4 cards):
  - Total Transactions
  - Approval Rate (%)
  - Declined Transactions
  - TPS (Load)
- **Live Insights Panel** with bullet points
- **System Health Panel** (4 health cards):
  - System Load (TPS)
  - Error Rate (%)
  - Latency (Simulated) - **BUG IDENTIFIED**
  - API Status (Healthy/Degraded)
- **Charts Section**:
  - Transaction Volume (Last 15 Minutes) - Line chart
  - Approval vs Decline Trends - Line chart
  - Transaction Breakdown by Status - Bar chart

### Health Metrics Behavior

#### API Data (from `http://localhost:8081/api/analytics/overview`)
**Current Values (monitored over 15 seconds):**
- Total TX: **59** (static - did not change)
- Approved TX: **27** (static - did not change)
- Approval Rate: **45.76%** (static - did not change)

**Calculated Metrics:**
- Declined TX: **32**
- Error Rate: **54.24%**
- TPS (Load): **0.98**
- API Status: **Degraded** (approval rate < 85%)

**Observation:** API values remained **completely static** over 15 seconds, suggesting:
- No new transactions are being processed
- Backend metrics may be cached or not updating
- Or system is idle

#### Frontend Behavior
- **Refresh Interval:** Every 5 seconds (`setInterval(fetchStats, 5000)`)
- **Data Source:** Real API calls to `/api/analytics/overview`
- **Network Activity:** Expected API calls every 5 seconds

### 🐛 **BUG IDENTIFIED: Latency Metric**

**Location:** `analytics_frontend/src/app/page.tsx:89`

**Issue:**
```tsx
<HealthCard label="Latency (Simulated)" value={`${Math.floor(120 + Math.random() * 40)} ms`} />
```

**Problem:**
- Uses `Math.random()` which generates a **new random value on EVERY render**
- This means the latency value changes constantly, even when the component re-renders for unrelated reasons
- Not tied to the 5-second API refresh cycle
- Creates visual "flickering" and inconsistent behavior

**Expected Behavior:**
- Should either:
  1. Come from the API (real latency data)
  2. Use a stable random seed that only changes when stats update
  3. Be memoized to only update when stats change

**Impact:** Medium - Visual inconsistency, misleading users about actual latency

### Visual Issues (Expected)
- Latency value will appear to "flicker" randomly
- Other metrics should update smoothly every 5 seconds (if API data changes)

---

## 2. Merchant Frontend (`http://localhost:3000/`)

### Overall Layout
Based on code analysis, the merchant frontend is a **landing/marketing page** with:
- **Navbar** with navigation links (Home, Simulate, History, Send Authorization)
- **Hero Section** (yellow/green background `#B5FF37`)
- **CTA Section**
- **Ad Banner**
- **Hero Transactions Section** (white background):
  - "Best Partner For Your Transactions" heading
  - Feature list with expandable descriptions
  - Image display
  - History card stack

### Health Metrics
**NONE** - This is a marketing/landing page, not a dashboard. No system health metrics are displayed on the main page.

### Functionality
- Landing page for merchants
- Navigation to `/auth` (simulation page) and `/history` (transaction history)
- No real-time metrics or health indicators

### Status
- ✅ **200 OK** - Page loads successfully
- ✅ **No errors detected** in basic connectivity test
- ✅ **No health metrics expected** - this is intentional

---

## 3. Network & API Analysis

### Analytics Frontend API Calls
- **Endpoint:** `http://localhost:8081/api/analytics/overview`
- **Frequency:** Every 5 seconds
- **Response Time:** Fast (< 100ms observed)
- **Data Stability:** Static (no changes observed over 15 seconds)

### Network Panel Observations (Expected)
When monitoring in browser DevTools:
- **Regular API calls** every 5 seconds to `/api/analytics/overview`
- **Status:** 200 OK
- **Response:** JSON with `total_tx`, `approved_tx`, `approval_rate`
- **No CORS errors** (same-origin or properly configured)

---

## 4. Summary of Issues

### Critical Issues
**None identified**

### Medium Issues
1. **Latency Metric Bug** (`analytics_frontend/src/app/page.tsx:89`)
   - Random value regenerates on every render
   - Should be memoized or come from API

### Low Issues
1. **Static API Data**
   - Metrics not updating (may be expected if system is idle)
   - Consider adding timestamp or "last updated" indicator

### Recommendations
1. **Fix Latency Metric:**
   ```tsx
   // Option 1: Memoize with stats dependency
   const latency = useMemo(() => 
     Math.floor(120 + Math.random() * 40), 
     [stats?.total_tx]
   );
   
   // Option 2: Use API data if available
   // Option 3: Use stable seed based on stats
   ```

2. **Add Loading States:**
   - Show loading indicator during API calls
   - Display "Last updated" timestamp

3. **Add Error Handling:**
   - Show error message if API calls fail
   - Retry logic for failed requests

4. **Consider Real Latency Data:**
   - If backend tracks latency, expose it via API
   - Replace simulated latency with real data

---

## 5. Test Results Summary

| Metric | Analytics Frontend | Merchant Frontend |
|--------|-------------------|-------------------|
| **Page Load** | ✅ 200 OK | ✅ 200 OK |
| **Health Metrics** | ✅ Present (4 cards) | ❌ N/A (landing page) |
| **API Integration** | ✅ Working (5s refresh) | ❌ N/A |
| **Data Updates** | ⚠️ Static (no new data) | ❌ N/A |
| **Visual Bugs** | ⚠️ Latency flickers | ✅ None detected |
| **Console Errors** | ❓ Not tested (browser access) | ❓ Not tested |
| **Network Errors** | ✅ None detected | ✅ None detected |

---