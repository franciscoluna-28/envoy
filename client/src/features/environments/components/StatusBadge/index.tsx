import { CheckCircle, XCircle, AlertCircle, Activity } from "lucide-react";

interface StatusBadgeProps {
  status?: string;
}

export function StatusBadge({ status }: StatusBadgeProps) {
  const configs = {
    completed: {
      icon: CheckCircle,
      color: "text-green-600",
      bg: "bg-green-50",
      label: "Success",
    },
    failed: {
      icon: XCircle,
      color: "text-red-600",
      bg: "bg-red-50",
      label: "Failed",
    },
    running: {
      icon: Activity,
      color: "text-blue-600",
      bg: "bg-blue-50",
      label: "Running",
    },
  };

  const config = configs[status as keyof typeof configs] || {
    icon: AlertCircle,
    color: "text-stone-400",
    bg: "bg-stone-50",
    label: "Unknown",
  };

  const Icon = config.icon;

  return (
    <div
      className={`inline-flex items-center gap-1.5 px-2 py-0.5 rounded-full ${config.bg} border border-transparent transition-colors`}
    >
      <Icon className={`h-3 w-3 ${config.color}`} />
      <span
        className={`text-[11px] font-bold uppercase tracking-wide ${config.color}`}
      >
        {config.label}
      </span>
    </div>
  );
}