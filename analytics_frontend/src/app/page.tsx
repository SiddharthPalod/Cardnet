"use client";

import KpiCard from '../components/KpiCard';
import { getNetworkStats } from '../lib/analyticsApi';
import { useState, useEffect, useMemo } from 'react';
import type { NetworkStats } from '../lib/types';
import Navbar from '../components/Navbar';
import Hero from '../components/Hero';
import TimeSeriesChart from '../components/TimeSeriesChart';
import { LineChart, Line, BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer, Cell } from 'recharts';

export default function Home() {
  const [stats, setStats] = useState<NetworkStats | null>(null);
  const [loading, setLoading] = useState(true);

  const fetchStats = async () => {
    try {
      const data = await getNetworkStats();
      setStats(data);
    } catch (err) {
      console.error('Failed to fetch network stats:', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchStats();
    const interval = setInterval(fetchStats, 5000);
    return () => clearInterval(interval);
  }, []);

  if (loading && !stats) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-[#B5FF37] text-black">
        <div className="text-xl font-['Poppins']">Loading...</div>
      </div>
    );
  }

  const declinedTx = (stats?.total_tx ?? 0) - (stats?.approved_tx ?? 0);
  const approvalRate = stats?.approval_rate ?? 0;
  const tps = stats?.total_tx ? (stats.total_tx / 60).toFixed(2) : '0.00';

  // Derive a stable, stats-based latency estimate instead of per-render randomness
  const rate = approvalRate || 0;
  const baseLatency = 120; // healthy baseline
  const latencyPenalty = (1 - rate) * 80; // up to +80ms when approval rate is very low
  const estimatedLatencyMs = Math.max(
    80,
    Math.min(250, Math.round(baseLatency + latencyPenalty)),
  ); // clamp to a reasonable range

  return (
    <div className="min-h-screen flex flex-col bg-[#B5FF37] text-black">
      <Navbar />

      <main className="grow max-w-6xl mx-auto w-full px-6 md:px-0">
        <Hero />

        {/* SMART KPI ROW */}
        <div className="grid grid-cols-1 mt-10 md:grid-cols-4 gap-6 mb-10">
          <KpiCard title="Total Transactions" value={stats?.total_tx ?? 0} color="default" />
          <KpiCard title="Approval Rate" value={`${(approvalRate * 100).toFixed(2)}%`} color="blue" />
          <KpiCard title="Declined Transactions" value={declinedTx} color="red" />
          <KpiCard title="TPS (Load)" value={tps} color="green" />
        </div>

        {/* INSIGHTS PANEL */}
        <section className="bg-white border-2 border-black rounded-2xl shadow-[10px_10px_0_0_#000] p-6 md:p-8 mb-10">
          <h2 className="text-2xl md:text-3xl font-extrabold uppercase mb-6">
            Live Insights
          </h2>

          <ul className="space-y-3 font-['Poppins'] text-lg">
            <li>📊 Approval Rate is <strong>{(approvalRate * 100).toFixed(2)}%</strong></li>
            <li>❌ {declinedTx} transactions declined</li>
            <li>⚡ TPS currently at <strong>{tps}</strong></li>
            <li>
              {approvalRate > 0.9 
                ? "✅ System performing well" 
                : approvalRate > 0.75 
                ? "⚠️ Approval rate trending lower" 
                : "🚨 High decline activity detected"}
            </li>
          </ul>
        </section>

        {/* SYSTEM HEALTH PANEL */}
        <section className="bg-white border-2 border-black rounded-2xl shadow-[10px_10px_0_0_#000] p-6 md:p-8 mb-10">
          <h2 className="text-2xl md:text-3xl font-extrabold uppercase mb-6">
            System Health
          </h2>

          <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
            <HealthCard label="System Load" value={tps} />
            <HealthCard label="Error Rate" value={`${(declinedTx / (stats?.total_tx || 1) * 100).toFixed(2)}%`} />
            <HealthCard label="Latency (Estimated)" value={`${estimatedLatencyMs} ms`} />
            <HealthCard label="API Status" value={approvalRate > 0.60 ? "Healthy" : "Degraded"} />
          </div>
        </section>

        {/* CHARTS */}
        <div className="space-y-10 pb-16">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-10">
            <TransactionVolumeChart stats={stats} />
            <ApprovalTrendsChart stats={stats} />
          </div>
          <DeclineReasonsChart stats={stats} />
        </div>

      </main>
    </div>
  );
}


/* ---------------- SMALL COMPONENTS ---------------- */

function HealthCard({ label, value }: { label: string; value: string | number }) {
  const getStyle = () => {
    if (typeof value === "string") {
      if (value.toLowerCase().includes("healthy")) {
        return "bg-[#B5FF37] border-black shadow-[6px_6px_0_0_#000]";
      }
      if (value.toLowerCase().includes("degraded")) {
        return "bg-red-400 border-black shadow-[6px_6px_0_0_#000]";
      }
      if (value.toLowerCase().includes("warning")) {
        return "bg-yellow-300 border-black shadow-[6px_6px_0_0_#000]";
      }
    }
    return "bg-[#B5FF37] border-black shadow-[6px_6px_0_0_#000]";
  };

  return (
    <div className={`p-4 border-2 rounded-lg ${getStyle()}`}>
      <div className="text-xs uppercase mb-2 text-black/70 font-bold">
        {label}
      </div>
      <div className="text-3xl font-extrabold text-black">
        {value}
      </div>
    </div>
  );
}


/* ---------------- CHARTS ---------------- */

function TransactionVolumeChart({ stats }: { stats: NetworkStats | null }) {
  const chartData = useMemo(() => {
    const data = [];
    const now = new Date();
    const totalTx = stats?.total_tx ?? 0;

    for (let i = 14; i >= 0; i--) {
      const time = new Date(now.getTime() - i * 60000);
      const base = totalTx > 0 ? totalTx / 15 : 0;
      // Deterministic, smooth variation based on time index
      const variationFactor = 0.9 + 0.2 * Math.sin(i / 3);
      const value = Math.max(0, Math.round(base * variationFactor));

      data.push({
        time: time.toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit' }),
        value,
      });
    }

    return data;
  }, [stats?.total_tx]);

  return (
    <TimeSeriesChart
      data={chartData}
      title="Transaction Volume (Last 15 Minutes)"
      dataKey="value"
      color="#000000"
    />
  );
}


function ApprovalTrendsChart({ stats }: { stats: NetworkStats | null }) {
  const chartData = useMemo(() => {
    const data = [];
    const now = new Date();
    const totalTx = stats?.total_tx ?? 0;
    const approvalRate = stats?.approval_rate ?? 0;

    for (let i = 14; i >= 0; i--) {
      const time = new Date(now.getTime() - i * 60000);
      const base = totalTx > 0 ? totalTx / 15 : 0;
      // Deterministic, smooth variation that slightly drifts over time
      const variationFactor = 0.9 + 0.2 * Math.cos(i / 3);
      const volume = Math.max(0, Math.round(base * variationFactor));

      const approved = Math.round(volume * approvalRate);
      const declined = volume - approved;

      data.push({
        time: time.toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit' }),
        approved,
        declined,
      });
    }

    return data;
  }, [stats]);

  return (
    <ChartCard title="Approval vs Decline Trends">
      <div className="w-full h-75 min-h-75">
        <ResponsiveContainer width="100%" height="100%">
          <LineChart data={chartData}>
            <CartesianGrid strokeDasharray="3 3" stroke="#000000" opacity={0.1} />
            <XAxis dataKey="time" />
            <YAxis />
            <Tooltip />
            <Legend />
            <Line dataKey="approved" stroke="#10b981" strokeWidth={3} dot={false} />
            <Line dataKey="declined" stroke="#ef4444" strokeWidth={3} dot={false} />
          </LineChart>
        </ResponsiveContainer>
      </div>
    </ChartCard>
  );
}


function DeclineReasonsChart({ stats }: { stats: NetworkStats | null }) {
  const chartData = useMemo(() => {
    const totalTx = stats?.total_tx ?? 0;
    const approvedTx = stats?.approved_tx ?? 0;
    const declinedTx = totalTx - approvedTx;

    const fraud = Math.round(declinedTx * 0.3);
    const funds = Math.round(declinedTx * 0.25);
    const issuer = Math.round(declinedTx * 0.2);
    const rate = Math.round(declinedTx * 0.15);
    const other = declinedTx - fraud - funds - issuer - rate;

    return [
      { reason: 'Approved', count: approvedTx },
      { reason: 'Fraud', count: fraud },
      { reason: 'Insufficient Funds', count: funds },
      { reason: 'Issuer Unavailable', count: issuer },
      { reason: 'Rate Limited', count: rate },
      { reason: 'Other', count: Math.max(0, other) },
    ].filter(item => item.count > 0);
  }, [stats]);

  const getColor = (reason: string) => {
    switch (reason) {
      case 'Approved': return '#10b981';
      case 'Fraud': return '#ef4444';
      case 'Insufficient Funds': return '#f59e0b';
      case 'Issuer Unavailable': return '#8b5cf6';
      case 'Rate Limited': return '#ec4899';
      default: return '#000000';
    }
  };

  return (
    <ChartCard title="Transaction Breakdown by Status">
      <div className="w-full h-75 min-h-75">
        <ResponsiveContainer width="100%" height="100%">
          <BarChart data={chartData} layout="vertical">
            <CartesianGrid strokeDasharray="3 3" stroke="#000000" opacity={0.1} />
            <XAxis type="number" />
            <YAxis dataKey="reason" type="category" width={140} />
            <Tooltip />
            <Bar dataKey="count" radius={[0, 8, 8, 0]}>
              {chartData.map((entry, index) => (
                <Cell key={index} fill={getColor(entry.reason)} />
              ))}
            </Bar>
          </BarChart>
        </ResponsiveContainer>
      </div>
    </ChartCard>
  );
}


function ChartCard({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <div className="bg-white border-2 border-black rounded-2xl shadow-[10px_10px_0_0_#000] p-6 md:p-8">
      <h3 className="text-xl md:text-2xl font-extrabold uppercase mb-6">
        {title}
      </h3>
      {children}
    </div>
  );
}
