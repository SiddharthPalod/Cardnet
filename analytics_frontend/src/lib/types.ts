export interface NetworkStats {
  total_tx: number;
  approved_tx: number;
  approval_rate: number;
}

export interface MerchantStats {
  merchant_id: string;
  tx_count: number;
  approved_tx: number;
  total_amount: number;
  approval_rate: number;
}

export interface BinStats {
  bin: string;
  tx_count: number;
  approved_tx: number;
  total_amount: number;
  approval_rate: number;
}

export interface RuleEffectiveness {
  rule_id: string;
  triggered_count: number;
  decline_count: number;
  decline_rate: number;
}

// Network-wide transaction history item for analytics UI
export interface NetworkHistoryItem {
  auth_id: string;
  merchant_id: string;
  amount: number;
  currency: string;
  status: string;
  reason?: string;
  created_at: string;
}
