import { createFileRoute } from "@tanstack/react-router";
import { useState } from "react";
import { Badge } from "@/components/ui/badge";
import { Spinner } from "@/components/ui/spinner";
import {
  useGetProject,
  useGetEnvironments,
  EnvironmentList,
  CreateEnvironmentModal,
} from "@/features/projects";

export const Route = createFileRoute("/app/projects/$projectId/")({
  component: RouteComponent,
});

function RouteComponent() {
  const { projectId } = Route.useParams();
  const [createDialogOpen, setCreateDialogOpen] = useState(false);

  const { data: project, isLoading: projectsLoading } =
    useGetProject(projectId);
  const { data: environments, isLoading: envsLoading } =
    useGetEnvironments(projectId);

  if (projectsLoading || envsLoading) {
    return (
      <div className="flex flex-col items-center justify-center h-[60vh] gap-4">
        <Spinner className="w-8 h-8 text-blue-600" />
        <p className="text-xs text-stone-400 animate-pulse">
          Loading infrastructure resources...
        </p>
      </div>
    );
  }

  return (
    <div className="flex min-w-full flex-col h-full">
      <header className="flex items-center justify-between py-6 border-b">
        <div className="flex flex-col gap-1">
          <div className="flex items-center gap-3">
            <h2 className="text-xl font-bold tracking-tighter text-stone-900 uppercase">
              {project?.name ?? "Unknown Project"}
            </h2>
          </div>

          <div className="flex items-center gap-3">
            <div className="flex items-center gap-1.5 text-[10px] font-bold text-green-600 uppercase tracking-tight">
              <div className="w-1.5 h-1.5 rounded-full bg-green-500 animate-pulse" />
              System Operational
            </div>
            <span className="text-stone-300 text-[10px]">•</span>
            <Badge
              variant="outline"
              className="text-[10px] border-stone-200 text-stone-500 h-5 px-1.5 font-medium uppercase"
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
