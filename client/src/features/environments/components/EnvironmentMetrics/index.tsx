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
            icon: <Server className="w-3 h-3 text-muted-foreground" />,
          },
          {
            label: "Database Engine",
            value: environmentData?.db_engine
              ? environmentData.db_engine.split(" on ")[0]
              : "PostgreSQL",
            icon: <Database className="w-3 h-3 text-muted-foreground" />,
          },
        ].map((stat) => (
          <div
            key={stat.label}
            className="border p-4 rounded-2xl bg-card flex flex-col gap-1.5"
          >
            <span className="text-[10px] font-medium text-muted-foreground">
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
      
    </>
  );
}
