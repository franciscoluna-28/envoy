import { Badge } from "@/components/ui/badge";
import type { Environment } from "@/features/types";
import { Database, Server, Clock } from "lucide-react";

type Props = {
  environmentData: Environment;
  totalMigrations: number;
};

export function EnvironmentMetrics({
  environmentData,
  totalMigrations = 0,
}: Props) {
  return (
    <>
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {[
          {
            label: "Network Status",
            value: "Connected",
            icon: ( 
              <div className="h-1.5 w-1.5 rounded-full bg-green-500 animate-pulse" />
            ),
          },
          {
            label: "Node Type",
            value: environmentData?.type || "Development",
            icon: <Server className="w-3 h-3 text-stone-400" />,
          },
          {
            label: "Database Engine",
            value: environmentData?.db_engine
              ? environmentData.db_engine.split(" on ")[0]
              : "PostgreSQL",
            icon: <Database className="w-3 h-3 text-stone-400" />,
          },
        ].map((stat) => (
          <div
            key={stat.label}
            className="border p-4 rounded-2xl bg-card flex flex-col gap-1.5"
          >
            <span className="text-[10px] font-bold text-stone-400 uppercase">
              {stat.label}
            </span>
            <div className="flex items-center gap-2">
              {stat.icon}
              <span className="text-sm font-medium text-stone-900 capitalize tracking-tight">
                {stat.value}
              </span>
            </div>
          </div>
        ))}
      </div>
      <div className="space-y-6">
        <div className="flex items-center justify-between border-b border-stone-100 pb-4">
          <div className="flex items-center gap-2">
            <Clock className="h-4 w-4 text-stone-400" />
            <h3 className="text-[11px] font-bold text-stone-500 uppercase tracking-[0.2em]">
              Migration History
            </h3>
          </div>
          <Badge
            variant="outline"
            className="text-[10px] border-stone-200 text-stone-500 font-bold uppercase"
          >
            {totalMigrations} Total
          </Badge>
        </div>
      </div>
    </>
  );
}
