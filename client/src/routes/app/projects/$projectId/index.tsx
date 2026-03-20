import { createFileRoute } from "@tanstack/react-router";
import { useState } from "react";
import { Badge } from "@/components/ui/badge";
import { CreateEnvironmentModal } from "@/features/environments/components/CreateEnvironmentModal";
import { EnvironmentList } from "@/features/environments/components/EnvironmentList";
import { useGetEnvironments } from "@/features/environments/hooks/useEnvironments";
import { useGetProject } from "@/features/projects/hooks/useProjects";
import { LoadingState } from "@/components/shared/LoadingState";
import { SystemOperationalBadge } from "@/features/projects/components/SystemOperationalBadge";
import { requireAuth } from "@/utils/guard";

export const Route = createFileRoute("/app/projects/$projectId/")({
  component: RouteComponent,
  beforeLoad: requireAuth,
});

function RouteComponent() {
  const { projectId } = Route.useParams();
  const [createDialogOpen, setCreateDialogOpen] = useState(false);

  const { data: project, isLoading: projectsLoading } =
    useGetProject(projectId);
  const { data: environments, isLoading: envsLoading } =
    useGetEnvironments(projectId);

  if (projectsLoading || envsLoading) {
    return <LoadingState />;
  }

  return (
    <div className="flex min-w-full flex-col h-full">
      <header className="flex items-center justify-between py-6 border-b">
        <div className="flex flex-col gap-1">
          <div className="flex items-center gap-3">
            <h2 className="text-xl font-bold tracking-tighter text-stone-900">
              {project?.name ?? "Unknown Project"}
            </h2>
          </div>

          <div className="flex items-center gap-3">
           <SystemOperationalBadge />
            <span className="text-stone-300 text-[10px]">•</span>
            <Badge
              variant="outline"
              
            >
              {environments?.length || 0} Environment
              {environments?.length !== 1 ? "s" : ""}
            </Badge>
          </div>
        </div>

        <div className="flex items-center gap-3">
          <CreateEnvironmentModal
            projectId={projectId}
            onCreated={() => setCreateDialogOpen(false)}
            open={createDialogOpen}
            onOpenChange={setCreateDialogOpen}
          />
        </div>
      </header>

      <div className="flex flex-col gap-8 py-6">
        <div className="max-w-full">
          <EnvironmentList environments={environments || []} />
        </div>
      </div>
    </div>
  );
}
