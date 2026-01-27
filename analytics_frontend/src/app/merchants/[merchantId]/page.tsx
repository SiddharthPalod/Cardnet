"use client";

import { use, useState, useEffect } from 'react';
import Link from 'next/link';
import MetricTable from '../../../components/MetricTable';
import { getMerchantStats } from '../../../lib/analyticsApi';
import type { MerchantStats } from '../../../lib/types';

export default function MerchantDetailPage({
  params,
}: {
  params: Promise<{ merchantId: string }>;
}) {
  const { merchantId } = use(params);
  const [stats, setStats] = useState<MerchantStats | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchStats = async () => {
      try {
        const data = await getMerchantStats(merchantId);
        setStats(data);
      } catch (err) {
        console.error('Failed to fetch merchant stats:', err);
      } finally {
        setLoading(false);
      }
    };
    fetchStats();
  }, [merchantId]);

  if (loading) {
    return (
      <div className="min-h-screen bg-[#B5FF37] text-black px-6 md:px-16 py-10">
        <div className="max-w-6xl mx-auto">
          <div className="text-center font-['Poppins'] text-xl">Loading...</div>
        </div>
      </div>
    );
  }

  if (!stats) {
    return (
      <div className="min-h-screen bg-[#B5FF37] text-black px-6 md:px-16 py-10">
        <div className="max-w-6xl mx-auto">
          <p className="font-['Poppins'] text-xl">Merchant not found</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-[#B5FF37] text-black px-6 md:px-16 py-10">
      <main className="max-w-6xl mx-auto">
        <div className="flex items-center justify-between gap-4 mb-10">
          <div>
            <h1 className="text-4xl md:text-5xl font-extrabold font-['Poppins'] uppercase">
              Merchant: {stats.merchant_id}
            </h1>
            <p className="mt-3 text-base md:text-lg font-['Poppins'] text-black/70">
              Detailed analytics and metrics for this merchant.
            </p>
          </div>
          <Link
            href="http://localhost:3000"
            className="hidden md:inline-flex items-center justify-center bg-black px-8 py-3 text-white text-sm md:text-base font-bold font-['Poppins'] uppercase hover:scale-[1.02] active:scale-95 transition-transform"
          >
            Merchant Portal
          </Link>
        </div>

        <div className="mb-6">
          <Link 
            href="/" 
            className="text-black font-['Poppins'] font-semibold uppercase hover:opacity-70 transition-opacity"
          >
            ← Back to Search
          </Link>
        </div>

        <MetricTable
          title="Merchant Metrics"
          data={[
            { label: 'Transaction Count', value: stats.tx_count ?? 0 },
            { label: 'Approved Transactions', value: stats.approved_tx ?? 0 },
            { label: 'Total Volume', value: `$${(stats.total_amount ?? 0).toFixed(2)}` },
            {
              label: 'Approval Rate',
              value:
                stats.approval_rate != null && !isNaN(stats.approval_rate)
                  ? `${(stats.approval_rate * 100).toFixed(2)}%`
                  : '0%',
            },
          ]}
        />

        <div className="mt-10 md:hidden">
          <Link
            href="http://localhost:3000"
            className="inline-flex w-full items-center justify-center bg-black px-6 py-3 text-white text-sm font-bold font-['Poppins'] uppercase hover:scale-[1.02] active:scale-95 transition-transform"
          >
            Merchant Portal
          </Link>
        </div>
      </main>
    </div>
  );
}
