"use client";

import type { AuthResponse } from '../app/lib/types';
import StatusBadge from './StatusBadge';
import JsonViewer from './JsonViewer';

interface AuthResultProps {
  response: AuthResponse | null;
  request: any;
  error?: string;
}

export default function AuthResult({ response, request, error }: AuthResultProps) {
  if (error) {
    return (
      <div className="bg-red-50 border border-red-200 dark:bg-red-900/20 dark:border-red-800 p-6 rounded-lg shadow-md">
        <h2 className="text-xl font-bold mb-2 text-red-800 dark:text-red-200">Error</h2>
        <p className="text-red-600 dark:text-red-400">{error}</p>
      </div>
    );
  }

  if (!response) {
    return null;
  }

  const isApproved = response.status === 'APPROVED';
  const isError = response.status === 'ERROR';

  return (
    <div
      className={`p-6 rounded-lg shadow-md border ${
        isApproved
          ? 'bg-green-50 border-green-200 dark:bg-green-900/20 dark:border-green-800'
          : isError
          ? 'bg-red-50 border-red-200 dark:bg-red-900/20 dark:border-red-800'
          : 'bg-yellow-50 border-yellow-200 dark:bg-yellow-900/20 dark:border-yellow-800'
      }`}
    >
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-xl font-bold">Authorization Result</h2>
        <StatusBadge status={response.status} />
      </div>

      <div className="space-y-2 mb-4">
        <div>
          <span className="font-semibold">Auth ID:</span>{' '}
          <span className="font-mono text-sm">{response.auth_id}</span>
        </div>
        {response.correlation_id && (
          <div>
            <span className="font-semibold">Correlation ID:</span>{' '}
            <span className="font-mono text-sm">{response.correlation_id}</span>
          </div>
        )}
        {response.reason && (
          <div>
            <span className="font-semibold">Reason:</span> {response.reason}
          </div>
        )}
        {response.risk_score !== undefined && (
          <div>
            <span className="font-semibold">Risk Score:</span> {response.risk_score}
          </div>
        )}
        {response.issuer_latency_ms !== undefined && (
          <div>
            <span className="font-semibold">Issuer Latency:</span> {response.issuer_latency_ms}ms
          </div>
        )}
      </div>

      <JsonViewer data={{ request, response }} title="Raw JSON" />
    </div>
  );
}
