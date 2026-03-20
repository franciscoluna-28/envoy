import { useState } from "react";
import { createFileRoute, Link } from "@tanstack/react-router";
import { Database } from "lucide-react";
import { Button } from "@/components/ui/button";
import { useGetEnvironment } from "@/features/environments/hooks/useEnvironments";
import { useGetEnvironmentMigrations } from "@/features/environments/hooks/useMigrations";
import { toast } from "sonner";
import { useEnvironmentSchema } from "@/features/environments/hooks/useEnvironmentSchema";
import { useGetProject } from "@/features/projects/hooks/useProjects";
import { LoadingState } from "@/components/shared/LoadingState";
import { EnvironmentMigrationsTable } from "@/features/environments/components/EnvironmentMigrationsTable";
import { DisplayEnvironmentSchemaModal } from "@/features/environments/components/DisplayEnvironmentSchemaModal";
import { EnvironmentMetrics } from "@/features/environments/components/EnvironmentMetrics";
import { Badge } from "@/components/ui/badge";
import { requireAuth } from "@/utils/guard";

export const Route = createFileRoute(
  "/app/projects/$projectId/environments/$envId/",
)({
  component: RouteComponent,
  beforeLoad: requireAuth,
});

function RouteComponent() {
  const params = Route.useParams();
  const { projectId, envId } = params;
  const [isModalOpen, setIsModalOpen] = useState(false);

  const { isPending: isLoadingProject } = useGetProject(projectId);
  const { data: environmentData, isPending: isLoadingEnvironment } =
    useGetEnvironment(projectId, envId);
  const { data: environmentSchema, isPending: isLoadingEnvironmentSchema } =
    useEnvironmentSchema(envId);
  const { data: migrations, isPending: isLoadingMigrations } =
    useGetEnvironmentMigrations(envId);

  if (isLoadingProject || isLoadingEnvironment) {
    return <LoadingState />;
  }

  return (
    <div className="flex flex-col h-full duration-500">
      <header className="flex flex-col gap-4 py-6 border-b bg-stone-50/30">
        <div className="flex items-end justify-between">
          <div className="space-y-1">
            <h1 className="text-xl font-bold tracking-tighter text-stone-900 flex items-center gap-3">
              <div className="p-2 bg-background border rounded-xl">
                <Database className="h-5 w-5 text-blue-600" />
              </div>
              {environmentData?.name}
            </h1>
            <div className="flex items-center gap-2 pl-12">
              <button
                onClick={() => {
                  navigator.clipboard.writeText(envId);
                  toast.success("Node ID copied to clipboard");
                }}
                className="text-[10px] font-medium text-muted-foreground hover:text-blue-600 transition-colors uppercase tracking-wide"
              >
                Node ID: #{envId.slice(0, 8)}
              </button>
            </div>
          </div>

          <div className="flex gap-3">
            <Button
              variant="outline"
              size="lg"
              onClick={() => setIsModalOpen(true)}
            >
              Preview Schema
            </Button>
            <Link
              to="/app/projects/$projectId/environments/$envId/migrations/new"
              params={{ projectId, envId }}
            >
              <Button size="lg">Create Migration</Button>
            </Link>
          </div>
        </div>
      </header>
      <div className="flex-1 py-8 space-y-12">
        <EnvironmentMetrics
          environmentData={environmentData || {}}
          totalMigrations={migrations?.length || 0}
        />

        <div className="space-y-6">
          <div className="flex items-center justify-between">
            <h3 className="text-sm font-semibold tracking-tight text-stone-900">
              Migration History
            </h3>

            <Badge variant="outline">{migrations?.length || 0} Total</Badge>
          </div>
          <EnvironmentMigrationsTable
            migrations={migrations || []}
            isLoading={isLoadingMigrations}
          />
        </div>
      </div>

      <DisplayEnvironmentSchemaModal
        onClose={() => {
          setIsModalOpen(false);
        }}
        open={isModalOpen}
        onOpenChange={setIsModalOpen}
        environmentData={environmentData || {}}
        environmentSchema={environmentSchema || []}
        isLoadingEnvironmentSchema={isLoadingEnvironmentSchema}
      />
    </div>
  );
}
