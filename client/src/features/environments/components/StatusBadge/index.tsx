import { CheckCircle2, XCircle, AlertCircle, Loader2 } from "lucide-react";

interface StatusBadgeProps {
  status?: string;
}

export function StatusBadge({ status }: StatusBadgeProps) {
  const configs = {
    completed: {
      icon: CheckCircle2,
      color: "text-emerald-700",
      bg: "bg-emerald-500/15",
      border: "border-emerald-500/20",
      label: "Success",
      animate: "",
    },
    failed: {
      icon: XCircle,
      color: "text-red-700",
      bg: "bg-red-500/15",
      border: "border-red-500/20",
      label: "Failed",
      animate: "",
    },
    running: {
      icon: Loader2,
      color: "text-blue-700",
      bg: "bg-blue-500/15",
      border: "border-blue-500/20",
      label: "Running",
      animate: "animate-spin", 
    },
  };

  const config = configs[status as keyof typeof configs] || {
    icon: AlertCircle,
    color: "text-zinc-500",
    bg: "bg-zinc-500/10",
    border: "border-zinc-500/10",
    label: "Unknown",
    animate: "",
  };

  const Icon = config.icon;

  return (
    <div
      className={`
        inline-flex items-center gap-1.5 px-2 py-0.5 rounded-full 
        ${config.bg} ${config.border} border
        transition-all duration-200 select-none
      `}
    >
      <Icon 
        size={12} 
        strokeWidth={2.5} 
        className={`${config.color} ${config.animate}`} 
      />
      <span className={`
        text-[10px] font-bold uppercase
        ${config.color} leading-none
      `}>
        {config.label}
      </span>
    </div>
  );
}