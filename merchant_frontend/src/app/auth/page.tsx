"use client";

import { useState } from 'react';
import Link from 'next/link';
import AuthForm from '../../components/AuthForm';
import AuthResult from '../../components/AuthResult';
import Navbar from '../../components/Navbar';
import { authorizePayment } from '../lib/api';
import type { AuthRequest, AuthResponse } from '../lib/types';

export default function AuthPage() {
  const [response, setResponse] = useState<AuthResponse | null>(null);
  const [request, setRequest] = useState<AuthRequest | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (req: AuthRequest) => {
    setLoading(true);
    setError(null);
    setRequest(req);

    try {
      const res = await authorizePayment(req);
      setResponse(res);
      
      // Store in localStorage for history
      const history = JSON.parse(localStorage.getItem('auth_history') || '[]');
      history.unshift({
        ...res,
        merchant_id: req.merchant_id,
        amount: req.amount,
        currency: req.currency,
        created_at: new Date().toISOString(),
      });
      localStorage.setItem('auth_history', JSON.stringify(history.slice(0, 100))); // Keep last 100
    } catch (err: any) {
      setError(err.message || 'Failed to authorize payment');
      setResponse(null);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-[#B5FF37] text-black px-6 md:px-16">
      <Navbar />
      <main className="max-w-5xl mx-auto py-10">
        <div className="flex items-center justify-between gap-4 mb-10">
          <div>
            <h1 className="text-4xl md:text-5xl font-extrabold font-['Poppins'] uppercase">
              Simulate Authorization
            </h1>
            <p className="mt-3 text-base md:text-lg font-['Poppins'] text-black/70">
              Submit merchant, card, and transaction details and see how CardNet responds.
            </p>
          </div>

        </div>

        <section className="bg-white border-2 border-black rounded-2xl shadow-[10px_10px_0_0_#000] p-6 md:p-8 mb-10">
          <h2 className="text-2xl md:text-3xl font-extrabold font-['Poppins'] uppercase mb-6">
            Transaction Details
          </h2>
          <AuthForm onSubmit={handleSubmit} loading={loading} />
        </section>

        {(response || error) && (
          <section className="bg-white border-2 border-black rounded-2xl shadow-[10px_10px_0_0_#000] p-6 md:p-8">
            <AuthResult response={response} request={request} error={error || undefined} />
          </section>
        )}

        <div className="mt-10 md:hidden">
          <Link
            href="/history"
            className="inline-flex w-full items-center justify-center bg-black px-6 py-3 text-white text-sm font-bold font-['Poppins'] uppercase"
          >
            View History
          </Link>
        </div>
      </main>
    </div>
  );
}
