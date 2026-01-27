"use client";

import { useEffect, useState } from "react";
import Navbar from "../../components/Navbar";
import type { NetworkHistoryItem } from "../../lib/types";
import { getNetworkHistory } from "../../lib/analyticsApi";

function StatusPill({ status }: { status: string }) {
  const normalized = status.toUpperCase();

  let bg = "bg-gray-200 text-gray-900 border-gray-400";
  if (normalized === "APPROVED") bg = "bg-emerald-200 text-emerald-900 border-emerald-500";
  else if (normalized.includes("DECLINE") || normalized === "DECLINED") bg = "bg-red-200 text-red-900 border-red-500";
  else if (normalized.includes("ERROR")) bg = "bg-yellow-200 text-yellow-900 border-yellow-500";

  return (
    <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-[10px] md:text-xs font-semibold border ${bg}`}>
      {normalized}
    </span>
  );
}

export default function AnalyticsHistoryPage() {
  const [items, setItems] = useState<NetworkHistoryItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [filter, setFilter] = useState({ merchant: "", status: "" });

  useEffect(() => {
    let cancelled = false;

    getNetworkHistory()
      .then((data) => {
        if (!cancelled) {
          setItems(data || []);
        }
      })
      .catch((err) => {
        if (!cancelled) {
          console.error("Failed to load network history:", err);
          setError("Unable to load history yet. Backend API not implemented.");
        }
      })
      .finally(() => {
        if (!cancelled) {
          setLoading(false);
        }
      });

    return () => {
      cancelled = true;
    };
  }, []);

  const filtered = items.filter((item) => {
    if (filter.merchant && !item.merchant_id.toLowerCase().includes(filter.merchant.toLowerCase())) {
      return false;
    }
    if (filter.status && item.status.toUpperCase() !== filter.status.toUpperCase()) {
      return false;
    }
    return true;
  });

  return (
    <div className="min-h-screen bg-[#B5FF37] text-black">
      <Navbar />

      <main className="max-w-6xl mx-auto px-6 md:px-0 py-10">
        <header className="flex flex-col md:flex-row md:items-end md:justify-between gap-4 mb-10">
          <div>
            <h1 className="text-4xl md:text-5xl font-extrabold font-['Poppins'] uppercase">
              Network Transaction History
            </h1>
            <p className="mt-3 text-base md:text-lg font-['Poppins'] text-black/70">
              Cross-merchant view of recent authorizations across the network.
            </p>
          </div>
        </header>

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
                <option value="DECLINED">Declined</option>
                <option value="SOFT_DECLINE">Soft Decline</option>
                <option value="HARD_DECLINE">Hard Decline</option>
                <option value="ERROR">Error</option>
              </select>
            </div>
          </div>
        </section>

        <section className="bg-white border-2 border-black rounded-2xl shadow-[10px_10px_0_0_#000] overflow-hidden">
          {loading ? (
            <div className="p-8 text-center font-['Poppins'] text-black/70">
              Loading network history...
            </div>
          ) : error ? (
            <div className="p-8 text-center font-['Poppins'] text-red-600">
              {error}
            </div>
          ) : filtered.length === 0 ? (
            <div className="p-8 text-center font-['Poppins'] text-black/60">
              No transactions found yet. Once the backend exposes network history, it will appear here.
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
                  {filtered.map((item) => (
                    <tr
                      key={item.auth_id}
                      className="border-t border-black/10 hover:bg-[#f5f5f5]"
                    >
                      <td className="px-4 py-3 text-xs md:text-sm">
                        {new Date(item.created_at).toLocaleString()}
                      </td>
                      <td className="px-4 py-3 text-xs md:text-sm font-mono">
                        {item.merchant_id}
                      </td>
                      <td className="px-4 py-3 text-xs md:text-sm">
                        {(item.amount / 100).toFixed(2)} {item.currency}
                      </td>
                      <td className="px-4 py-3">
                        <StatusPill status={item.status} />
                      </td>
                      <td className="px-4 py-3 text-xs md:text-sm">
                        {item.reason || "-"}
                      </td>
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
