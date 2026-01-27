import type { NetworkStats, MerchantStats, BinStats, RuleEffectiveness, NetworkHistoryItem } from './types';

const API_BASE = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8081';

async function fetchApi<T>(path: string): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`);
  if (!res.ok) {
    throw new Error(`HTTP ${res.status}: ${await res.text()}`);
  }
  return res.json();
}

export async function getNetworkStats(): Promise<NetworkStats> {
  return fetchApi<NetworkStats>('/api/analytics/overview');
}

export async function getMerchantStats(merchantId: string): Promise<MerchantStats> {
  return fetchApi<MerchantStats>(`/api/analytics/merchant?merchant_id=${merchantId}`);
}

export async function getBinStats(bin: string): Promise<BinStats> {
  return fetchApi<BinStats>(`/api/analytics/bin?bin=${bin}`);
}

export async function getRuleEffectiveness(ruleId: string): Promise<RuleEffectiveness> {
  return fetchApi<RuleEffectiveness>(`/api/analytics/rules?rule_id=${ruleId}`);
}

export async function getNetworkHistory(): Promise<NetworkHistoryItem[]> {
  return fetchApi<NetworkHistoryItem[]>('/api/analytics/history');
}
