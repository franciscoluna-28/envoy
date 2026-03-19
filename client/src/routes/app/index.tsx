import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'
import { Badge } from '@/components/ui/badge'
import { CreateProjectModal } from '@/features/projects/components/CreateProjectModal'
import { DeleteProjectModal } from '@/features/projects/components/DeleteProjectModal'
import { ProjectCard } from '@/features/projects/components/ProjectCard'
import { UpdateProjectModal } from '@/features/projects/components/UpdateProjectModal'
import { useGetAllProjects } from '@/features/projects/hooks/useProjects'
import { LoadingState } from '@/components/shared/LoadingState'

export const Route = createFileRoute('/app/')({
  component: RouteComponent
})

function RouteComponent() {
  const { data: projects, isLoading } = useGetAllProjects()
  const [createDialogOpen, setCreateDialogOpen] = useState(false)
  const [updateDialogOpen, setUpdateDialogOpen] = useState<string | null>(null)
  const [deleteDialogOpen, setDeleteDialogOpen] = useState<string | null>(null)

  const currentUpdateProject = projects?.find(p => p.id === updateDialogOpen)
  const currentDeleteProject = projects?.find(p => p.id === deleteDialogOpen)

  if (isLoading) {
    return <LoadingState />;
  }

  return (
    <div className="flex flex-col h-full">
      <header className="flex items-center justify-between py-6 border-b">
        <div className="flex flex-col gap-1">
          <div className="flex items-center gap-3">
            <h2 className="text-xl font-bold tracking-tighter text-stone-900 uppercase">
              Infrastructure / Projects
            </h2>
          </div>
         <div className="flex items-center gap-3">
             <div className="flex items-center gap-1.5 text-[10px] font-bold text-green-600 uppercase tracking-tight">
                <div className="w-1.5 h-1.5 rounded-full bg-green-500 animate-pulse" />
                System Operational
             </div>
             <span className="text-stone-300 text-[10px]">•</span>
            <Badge variant="outline" className="text-[10px] border-stone-200 text-stone-500 h-5 px-1.5 font-medium uppercase">
              {projects?.length || 0} Project{projects?.length !== 1 ? 's' : ''} 
            </Badge>
          </div>
        </div>

        <CreateProjectModal 
          onCreated={() => setCreateDialogOpen(false)} 
          open={createDialogOpen} 
          onOpenChange={setCreateDialogOpen} 
        />
      </header>

      {currentUpdateProject && (
        <UpdateProjectModal 
          project={currentUpdateProject}
          open={updateDialogOpen !== null}
          onUpdated={() => setUpdateDialogOpen(null)}
          onOpenChange={(isOpen) => !isOpen && setUpdateDialogOpen(null)}
        />
      )}

      {currentDeleteProject && (
        <DeleteProjectModal 
          project={currentDeleteProject}
          open={deleteDialogOpen !== null}
          onDeleted={() => setDeleteDialogOpen(null)}
          onOpenChange={(isOpen) => !isOpen && setDeleteDialogOpen(null)}
        />
      )}

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6 py-6">
          {projects?.map((project) => (
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