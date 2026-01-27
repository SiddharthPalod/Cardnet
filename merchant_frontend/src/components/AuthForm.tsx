"use client";

import { useState } from 'react';
import type { AuthRequest } from '../app/lib/types';

interface AuthFormProps {
  onSubmit: (request: AuthRequest) => void;
  loading: boolean;
}

const SAMPLE_MERCHANTS = [
  { id: 'amazon', name: 'Amazon' },
  { id: 'walmart', name: 'Walmart' },
  { id: 'target', name: 'Target' },
  { id: 'high-risk-merchant', name: 'High Risk Merchant' },
];

const SAMPLE_CARDS = [
  '1234-5678-9012-3456',
  '4111-1111-1111-1111',
  '5555-5555-5555-4444',
  '4000-0000-0000-0002', // Rate-limited BIN
];

export default function AuthForm({ onSubmit, loading }: AuthFormProps) {
  const [formData, setFormData] = useState<AuthRequest>({
    merchant_id: 'amazon',
    card_token: '1234-5678-9012-3456',
    amount: 10000,
    currency: 'USD',
    mcc: '5411',
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit(formData);
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div>
        <label className="block text-xs md:text-sm font-semibold font-['Poppins'] uppercase mb-2">
          Merchant ID
        </label>
        <select
          className="w-full border-2 border-black rounded-lg px-3 py-2 font-['Poppins'] text-sm bg-white"
          value={formData.merchant_id}
          onChange={(e) => setFormData({ ...formData, merchant_id: e.target.value })}
        >
          {SAMPLE_MERCHANTS.map((m) => (
            <option key={m.id} value={m.id}>
              {m.name} ({m.id})
            </option>
          ))}
        </select>
      </div>

      <div>
        <label className="block text-xs md:text-sm font-semibold font-['Poppins'] uppercase mb-2">
          Card Token
        </label>
        <select
          className="w-full border-2 border-black rounded-lg px-3 py-2 font-['Poppins'] text-sm bg-white"
          value={formData.card_token}
          onChange={(e) => setFormData({ ...formData, card_token: e.target.value })}
        >
          {SAMPLE_CARDS.map((card) => (
            <option key={card} value={card}>
              {card}
            </option>
          ))}
        </select>
        <input
          type="text"
          className="mt-2 w-full border-2 border-black rounded-lg px-3 py-2 font-['Poppins'] text-sm bg-white"
          placeholder="Or enter custom card token"
          value={formData.card_token}
          onChange={(e) => setFormData({ ...formData, card_token: e.target.value })}
        />
      </div>

      <div className="grid grid-cols-3 gap-4">
        <div>
          <label className="block text-xs md:text-sm font-semibold font-['Poppins'] uppercase mb-2">
            Amount (cents)
          </label>
          <input
            type="number"
            className="w-full border-2 border-black rounded-lg px-3 py-2 font-['Poppins'] text-sm bg-white"
            value={formData.amount}
            onChange={(e) => setFormData({ ...formData, amount: parseInt(e.target.value) || 0 })}
            min="1"
            required
          />
        </div>
        <div>
          <label className="block text-xs md:text-sm font-semibold font-['Poppins'] uppercase mb-2">
            Currency
          </label>
          <input
            className="w-full border-2 border-black rounded-lg px-3 py-2 font-['Poppins'] text-sm bg-white"
            value={formData.currency}
            onChange={(e) => setFormData({ ...formData, currency: e.target.value })}
            required
          />
        </div>
        <div>
          <label className="block text-xs md:text-sm font-semibold font-['Poppins'] uppercase mb-2">
            MCC
          </label>
          <input
            className="w-full border-2 border-black rounded-lg px-3 py-2 font-['Poppins'] text-sm bg-white"
            value={formData.mcc}
            onChange={(e) => setFormData({ ...formData, mcc: e.target.value })}
            placeholder="5411"
            required
          />
        </div>
      </div>

      <button
        type="submit"
        disabled={loading}
        className="mt-4 w-full bg-black px-6 py-3 text-white text-sm md:text-base font-bold font-['Poppins'] uppercase transition-transform disabled:opacity-50 disabled:cursor-not-allowed hover:scale-[1.02] active:scale-95"
      >
        {loading ? 'Processing...' : 'Send Authorization Request'}
      </button>
    </form>
  );
}
