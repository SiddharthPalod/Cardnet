"use client";

import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';

interface TimeSeriesChartProps {
  data: Array<{ time: string; value: number; label?: string }>;
  title: string;
  dataKey?: string;
  color?: string;
}

export default function TimeSeriesChart({
  data,
  title,
  dataKey = 'value',
  color = '#000000',
}: TimeSeriesChartProps) {
  if (data.length === 0) {
    return (
      <div className="bg-white border-2 border-black rounded-2xl shadow-[10px_10px_0_0_#000] p-6 md:p-8">
        <h3 className="text-xl md:text-2xl font-extrabold font-['Poppins'] uppercase mb-6">{title}</h3>
        <div className="h-64 flex items-center justify-center text-black/60 font-['Poppins']">
          No data available
        </div>
      </div>
    );
  }

  return (
    <div className="bg-white border-2 border-black rounded-2xl shadow-[10px_10px_0_0_#000] p-6 md:p-8">
      <h3 className="text-xl md:text-2xl font-extrabold font-['Poppins'] uppercase mb-6">{title}</h3>
      <ResponsiveContainer width="100%" height={300}>
        <LineChart data={data}>
          <CartesianGrid strokeDasharray="3 3" stroke="#000000" opacity={0.1} />
          <XAxis
            dataKey="time"
            tick={{ fill: '#000000', fontFamily: 'Poppins', fontSize: 12 }}
          />
          <YAxis tick={{ fill: '#000000', fontFamily: 'Poppins', fontSize: 12 }} />
          <Tooltip
            contentStyle={{
              backgroundColor: '#ffffff',
              border: '2px solid #000000',
              borderRadius: '12px',
              fontFamily: 'Poppins',
            }}
          />
          <Legend wrapperStyle={{ fontFamily: 'Poppins' }} />
          <Line
            type="monotone"
            dataKey={dataKey}
            stroke={color}
            strokeWidth={3}
            dot={{ r: 5, fill: color }}
            name={title}
          />
        </LineChart>
      </ResponsiveContainer>
    </div>
  );
}
