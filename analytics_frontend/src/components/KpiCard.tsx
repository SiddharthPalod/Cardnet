"use client";

interface KpiCardProps {
  title: string;
  value: string | number;
  subtitle?: string;
  color?: 'default' | 'green' | 'blue' | 'red' | 'yellow';
}

export default function KpiCard({ title, value, subtitle, color = 'default' }: KpiCardProps) {
  const colorClasses = {
    default: 'text-black',
    green: 'text-green-600',
    blue: 'text-blue-600',
    red: 'text-red-600',
    yellow: 'text-yellow-600',
  };

  return (
    <div className="bg-white border-2 border-black rounded-2xl shadow-[10px_10px_0_0_#000] p-6 md:p-8">
      <h3 className="text-xs md:text-sm font-semibold font-['Poppins'] uppercase mb-3 text-black/60">
        {title}
      </h3>
      <p className={`text-4xl md:text-5xl font-extrabold font-['Poppins'] ${colorClasses[color]}`}>
        {value}
      </p>
      {subtitle && (
        <p className="text-sm font-['Poppins'] text-black/60 mt-3">{subtitle}</p>
      )}
    </div>
  );
}
