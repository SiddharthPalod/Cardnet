"use client";

import { useState, useEffect } from 'react';
import Link from 'next/link';
import StatusBadge from '../../components/StatusBadge';
import Navbar from '../../components/Navbar';
import type { AuthHistoryItem } from '../lib/types';

export default function HistoryPage() {
  const [history, setHistory] = useState<AuthHistoryItem[]>([]);
  const [filter, setFilter] = useState({ merchant: '', status: '' });

  useEffect(() => {
    const stored = localStorage.getItem('auth_history');
    if (stored) {
      setHistory(JSON.parse(stored));
    }
  }, []);

  const filteredHistory = history.filter((item) => {
    if (filter.merchant && !item.merchant_id.includes(filter.merchant)) return false;
    if (filter.status && item.status !== filter.status) return false;
    return true;
  });

  return (
    <div className="min-h-screen bg-[#B5FF37] text-black px-6 md:px-16">
      <Navbar />
      <main className="max-w-6xl mx-auto py-10">
        <div className="flex items-center justify-between gap-4 mb-10">
          <div>
            <h1 className="text-4xl md:text-5xl font-extrabold font-['Poppins'] uppercase">
              Authorization History
            </h1>
            <p className="mt-3 text-base md:text-lg font-['Poppins'] text-black/70">
              Browse past authorization attempts and outcomes stored in your browser.
            </p>
          </div>

        </div>

        <section className="bg-white border-2 border-black rounded-2xl shadow-[10px_10px_0_0_#000] p-6 md:p-8 mb-8">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label className="block text-xs md:text-sm font-semibold font-['Poppins'] uppercase mb-2">
                Filter by Merchant
              </label>
              <input
                type="text"
                className="w-full border-2 border-black rounded-lg px-3 py-2 font-['Poppins'] text-sm bg-white"
                placeholder="Merchant ID..."
                value={filter.merchant}
                onChange={(e) => setFilter({ ...filter, merchant: e.target.value })}
              />
            </div>
            <div>
              <label className="block text-xs md:text-sm font-semibold font-['Poppins'] uppercase mb-2">
                Filter by Status
              </label>
              <select
                className="w-full border-2 border-black rounded-lg px-3 py-2 font-['Poppins'] text-sm bg-white"
                value={filter.status}
                onChange={(e) => setFilter({ ...filter, status: e.target.value })}
              >
                <option value="">All Statuses</option>
                <option value="APPROVED">Approved</option>
                <option value="SOFT_DECLINE">Soft Decline</option>
                <option value="HARD_DECLINE">Hard Decline</option>
                <option value="ERROR">Error</option>
              </select>
            </div>
          </div>
        </section>

        <section className="bg-white border-2 border-black rounded-2xl shadow-[10px_10px_0_0_#000] overflow-hidden">
          {filteredHistory.length === 0 ? (
            <div className="p-8 text-center text-black/60 font-['Poppins']">
              No authorization history found. Send some requests from the{" "}
              <Link href="/auth" className="underline font-semibold">
                Auth page
              </Link> 
              .
            </div>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full border-t border-black/10">
                <thead className="bg-black text-white font-['Poppins'] text-xs md:text-sm uppercase">
                  <tr>
                    <th className="px-4 py-3 text-left">Time</th>
                    <th className="px-4 py-3 text-left">Merchant</th>
                    <th className="px-4 py-3 text-left">Amount</th>
                    <th className="px-4 py-3 text-left">Status</th>
                    <th className="px-4 py-3 text-left">Reason</th>
                    <th className="px-4 py-3 text-left">Auth ID</th>
                  </tr>
                </thead>
                <tbody>
                  {filteredHistory.map((item, index) => (
                    <tr
                      key={item.auth_id ? `${item.auth_id}-${index}` : `${item.created_at}-${index}`}
                      className="border-t border-black/10 hover:bg-[#f5f5f5]"
                    >
                      <td className="px-4 py-3 text-xs md:text-sm">
                        {new Date(item.created_at).toLocaleString()}
                      </td>
                      <td className="px-4 py-3 text-xs md:text-sm font-mono">{item.merchant_id}</td>
                      <td className="px-4 py-3 text-xs md:text-sm">
                        {(item.amount / 100).toFixed(2)} {item.currency}
                      </td>
                      <td className="px-4 py-3">
                        <StatusBadge status={item.status} />
                      </td>
                      <td className="px-4 py-3 text-xs md:text-sm">{item.reason || "-"}</td>
                      <td className="px-4 py-3 text-[10px] md:text-xs font-mono break-all">
                        {item.auth_id}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </section>

      </main>
    </div>
  );
}
