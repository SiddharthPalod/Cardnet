import { v4 as uuidv4 } from 'uuid';
import type { AuthRequest, AuthResponse } from './types';

// Ensure API_BASE doesn't end with a slash and doesn't include /api
const getApiBase = () => {
  const base = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8081';
  // Remove trailing slash and any /api suffix to avoid duplication
  return base.replace(/\/+$/, '').replace(/\/api$/, '');
};

const API_BASE = getApiBase();

export async function apiFetch<T>(
  path: string,
  options: RequestInit = {}
): Promise<T> {
  const correlationId = uuidv4();

  // Ensure path starts with /
  const normalizedPath = path.startsWith('/') ? path : `/${path}`;
  const url = `${API_BASE}${normalizedPath}`;

  const res = await fetch(url, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      'X-Correlation-Id': correlationId,
      ...(options.headers || {}),
    },
  });

  if (!res.ok) {
    const errorText = await res.text();
    throw new Error(errorText || `HTTP ${res.status}`);
  }

  return res.json();
}

export async function authorizePayment(request: AuthRequest): Promise<AuthResponse> {
  return apiFetch<AuthResponse>('/api/auth/authorize', {
    method: 'POST',
    body: JSON.stringify(request),
  });
}
