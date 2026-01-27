"use client";

interface MetricTableProps {
  title: string;
  data: Array<{ label: string; value: string | number }>;
}

export default function MetricTable({ title, data }: MetricTableProps) {
  return (
    <div className="bg-white border-2 border-black rounded-2xl shadow-[10px_10px_0_0_#000] p-6 md:p-8">
      <h3 className="text-xl md:text-2xl font-extrabold font-['Poppins'] uppercase mb-6">{title}</h3>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {data.map((item, idx) => (
          <div key={idx} className="p-4 bg-[#B5FF37] border-2 border-black rounded-lg">
            <div className="text-xs md:text-sm font-semibold font-['Poppins'] uppercase mb-2 text-black/60">
              {item.label}
            </div>
            <div className="text-2xl md:text-3xl font-extrabold font-['Poppins'] text-black">
              {item.value}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
