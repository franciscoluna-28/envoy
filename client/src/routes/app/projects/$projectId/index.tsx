import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'
import { 
  useGetProject, 
  useGetEnvironments, 
  EnvironmentList,
  CreateEnvironmentModal
} from '@/features/projects'

export const Route = createFileRoute('/app/projects/$projectId/')({
  component: RouteComponent,
})

function RouteComponent() {
  const { projectId } = Route.useParams()
  const [createDialogOpen, setCreateDialogOpen] = useState(false)

  const { data: project, isLoading: projectsLoading } = useGetProject(projectId)
  const { data: environments, isLoading: envsLoading } = useGetEnvironments(projectId)

  if (projectsLoading || envsLoading) {
    return <div className="p-8 text-center animate-pulse text-muted-foreground">Loading infra...</div>
  }

  return (
    <div className="flex flex-col gap-8 p-6 max-w-6xl mx-auto">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">{project?.data?.name}</h1>
          <p className="text-muted-foreground text-sm">
            Manage your database environments and schema migrations.
          </p>
        </div>
        <CreateEnvironmentModal 
          projectId={projectId}
          onCreated={() => setCreateDialogOpen(false)} 
          open={createDialogOpen} 
          onOpenChange={setCreateDialogOpen} 
        />
      </div>

      <EnvironmentList 
        environments={environments?.data || []} 
        projectId={projectId} 
      />
    </div>
  )
}
