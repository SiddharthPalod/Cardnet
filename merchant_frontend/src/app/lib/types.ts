export interface AuthRequest {
  merchant_id: string;
  card_token: string;
  amount: number;
  currency: string;
  mcc: string;
}

export interface AuthResponse {
  auth_id: string;
  status: 'APPROVED' | 'SOFT_DECLINE' | 'HARD_DECLINE' | 'ERROR';
  reason?: string;
  correlation_id?: string;
  risk_score?: number;
  issuer_latency_ms?: number;
}

export interface AuthHistoryItem {
  auth_id: string;
  merchant_id: string;
  amount: number;
  currency: string;
  status: string;
  reason?: string;
  created_at: string;
  issuer_latency_ms?: number;
}
