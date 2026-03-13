import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'
import { Spinner } from '@/components/ui/spinner'
import { 
  ProjectCard, 
  CreateProjectModal,
  UpdateProjectModal,
  DeleteProjectModal,
  useGetAllProjects
} from '@/features/projects'

export const Route = createFileRoute('/app/')({
  component: RouteComponent
})

function RouteComponent() {
  const { data: projects, isLoading } = useGetAllProjects()
  const [createDialogOpen, setCreateDialogOpen] = useState(false)
  const [updateDialogOpen, setUpdateDialogOpen] = useState<string | null>(null)
  const [deleteDialogOpen, setDeleteDialogOpen] = useState<string | null>(null)

  const currentUpdateProject = projects?.data?.find(p => p.id === updateDialogOpen)
  const currentDeleteProject = projects?.data?.find(p => p.id === deleteDialogOpen)

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <Spinner/>
      </div>
    )
  }

  return (
    <div className="max-w-6xl mx-auto w-full px-4 py-8">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-10 border-b pb-6">
        <div>
          <h2 className="text-3xl font-bold tracking-tight text-foreground">Projects</h2>
          <p className="text-sm text-muted-foreground mt-1">
            Scale and manage your deployment infrastructure.
          </p>
        </div>
        <CreateProjectModal onCreated={() => setCreateDialogOpen(false)} open={createDialogOpen} onOpenChange={setCreateDialogOpen} />
      </div>

      {currentUpdateProject && (
        <UpdateProjectModal 
          project={currentUpdateProject}
          open={updateDialogOpen !== null}
          onUpdated={() => setUpdateDialogOpen(null)}
          onOpenChange={(isOpen) => {
            if (!isOpen) setUpdateDialogOpen(null)
          }}
        />
      )}

      {currentDeleteProject && (
        <DeleteProjectModal 
          project={currentDeleteProject}
          open={deleteDialogOpen !== null}
          onDeleted={() => setDeleteDialogOpen(null)}
          onOpenChange={(isOpen) => {
            if (!isOpen) setDeleteDialogOpen(null)
          }}
        />
      )}

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5">
        {projects?.data?.map((project) => (
          <ProjectCard 
            key={project.id || 'unknown'} 
            project={project}
            onUpdate={() => setUpdateDialogOpen(project.id || null)}
            onDelete={() => setDeleteDialogOpen(project.id || null)}
          />
        ))}
      </div>
    </div>
  )
}