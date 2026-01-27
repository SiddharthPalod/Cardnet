"use client";

interface StatusBadgeProps {
  status: string;
}

export default function StatusBadge({ status }: StatusBadgeProps) {
  const getStatusColor = (status: string) => {
    switch (status) {
      case 'APPROVED':
        return 'bg-green-500 text-white';
      case 'SOFT_DECLINE':
        return 'bg-yellow-500 text-white';
      case 'HARD_DECLINE':
        return 'bg-red-500 text-white';
      case 'ERROR':
        return 'bg-red-600 text-white';
      default:
        return 'bg-gray-500 text-white';
    }
  };

  return (
    <span className={`px-3 py-1 rounded-full text-sm font-semibold ${getStatusColor(status)}`}>
      {status}
    </span>
  );
}
