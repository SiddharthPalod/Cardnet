"use client";

import { useState } from 'react';

interface JsonViewerProps {
  data: any;
  title?: string;
}

export default function JsonViewer({ data, title = 'JSON' }: JsonViewerProps) {
  const [isOpen, setIsOpen] = useState(false);

  return (
    <div className="mt-4">
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="text-xs md:text-sm font-semibold font-['Poppins'] uppercase text-black hover:opacity-80"
      >
        {isOpen ? '▼' : '▶'} {title} (click to {isOpen ? 'collapse' : 'expand'})
      </button>
      {isOpen && (
        <pre className="mt-3 p-4 bg-white border-2 border-black rounded-xl overflow-auto text-xs font-mono">
          {JSON.stringify(data, null, 2)}
        </pre>
      )}
    </div>
  );
}
