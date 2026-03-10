import { createFileRoute, Link } from '@tanstack/react-router'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Plus, ArrowRight, MoreHorizontal } from 'lucide-react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import client from '@/api/client'
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog'
import { Field, FieldGroup, FieldLabel, FieldDescription, FieldError } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import { DialogClose } from '@/components/ui/dialog'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Spinner } from '@/components/ui/spinner'
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from '@/components/ui/dropdown-menu'
import { useState } from 'react'


export const Route = createFileRoute('/app/')({
  component: RouteComponent
})

const createProjectSchema = z.object({
  name: z.string().min(1, 'Project name is required').max(100, 'Project name must be less than 100 characters')
})

const updateProjectSchema = z.object({
  name: z.string().min(1, 'Project name is required').max(100, 'Project name must be less than 100 characters')
})

type CreateProjectForm = z.infer<typeof createProjectSchema>
type UpdateProjectForm = z.infer<typeof updateProjectSchema>

const PROJECTS_QUERY_KEYS = {
  all: ['projects'] as const,
  lists: () => [...PROJECTS_QUERY_KEYS.all, 'list'] as const,
  list: (filter: string) => [...PROJECTS_QUERY_KEYS.lists(), { filter }] as const,
  details: () => [...PROJECTS_QUERY_KEYS.all, 'detail'] as const,
  detail: (id: string) => [...PROJECTS_QUERY_KEYS.details(), id] as const,
}

function useCreateProject() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (data: CreateProjectForm) => {
      const response = await client.POST('/api/v1/projects', {
        body: data
      })
      return response.data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: PROJECTS_QUERY_KEYS.lists() })
    },
  })
}

function useUpdateProject() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async ({ id, data }: { id: string; data: UpdateProjectForm }) => {
      const response = await client.PUT('/api/v1/projects/{id}', {
        params: { path: { id } },
        body: data
      })
      return response.data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: PROJECTS_QUERY_KEYS.lists() })
    },
  })
}

function useDeleteProject() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (id: string) => {
      const response = await client.DELETE('/api/v1/projects/{id}', {
        params: { path: { id } }
      })
      return response.data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: PROJECTS_QUERY_KEYS.lists() })
    },
  })
}

function useGetAllProjects() {
  return useQuery({
    queryKey: PROJECTS_QUERY_KEYS.lists(),
    queryFn: async () => {
      const response = await client.GET('/api/v1/projects')
      return response.data
    }
  })
}



type CreateProjectDialogProps = {
  onCreated?: () => void
  open?: boolean
  onOpenChange?: (open: boolean) => void
}

type UpdateProjectDialogProps = {
  project: { id?: string; name?: string }
  onUpdated?: () => void
  open?: boolean
  onOpenChange?: (open: boolean) => void
}

type DeleteProjectDialogProps = {
  project: { id?: string; name?: string }
  onDeleted?: () => void
  open?: boolean
  onOpenChange?: (open: boolean) => void
}

function CreateProjectDialog({ onCreated, open, onOpenChange }: CreateProjectDialogProps) {



  const form = useForm<CreateProjectForm>({
    resolver: zodResolver(createProjectSchema),
    defaultValues: {
      name: ''
    }
  })

  const { mutate, isPending } = useCreateProject()

  const onSubmit = (data: CreateProjectForm) => {
    mutate(data, {
      onSuccess: () => {
        form.reset()
        onCreated?.()
      }
    })
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogTrigger asChild>
        <Button variant="outline" className="font-semibold">
          <Plus className="w-4 h-4 mr-2" />
          New Project
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Create New Project</DialogTitle>
          <DialogDescription>
            Create a logical container for your environments.
          </DialogDescription>
        </DialogHeader>
        
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
          <FieldGroup>
            <Field>
              <FieldLabel htmlFor="name">Project Name</FieldLabel>
              <Input
                id="name"
                placeholder="e.g. Phoenix-Core-DB"
                {...form.register('name')}
              />
              <FieldDescription>
                This will be used to identify your infrastructure clusters.
              </FieldDescription>
              <FieldError errors={form.formState.errors.name ? [form.formState.errors.name] : []} />
            </Field>
          </FieldGroup>
            
          <DialogFooter className="gap-2 sm:gap-0">
            <DialogClose asChild>
              <Button type="button" variant="ghost">
                Cancel
              </Button>
            </DialogClose>
            <Button type="submit" disabled={isPending}>
              {isPending ? 'Creating...' : 'Create Project'}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}

function UpdateProjectDialog({ project, onUpdated, open, onOpenChange }: UpdateProjectDialogProps & { open?: boolean, onOpenChange?: (open: boolean) => void }) {
  const form = useForm<UpdateProjectForm>({
    resolver: zodResolver(updateProjectSchema),
    defaultValues: {
      name: project.name || ''
    }
  })

  const { mutate, isPending } = useUpdateProject()

  const onSubmit = (data: UpdateProjectForm) => {
    if (!project.id) return
    mutate({ id: project.id as string, data }, {
      onSuccess: () => {
        onUpdated?.()
        onOpenChange?.(false)
      }
    })
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Update Project</DialogTitle>
          <DialogDescription>
            Update the project name.
          </DialogDescription>
        </DialogHeader>
        
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
          <FieldGroup>
            <Field>
              <FieldLabel htmlFor="name">Project Name</FieldLabel>
              <Input
                id="name"
                placeholder="e.g. Phoenix-Core-DB"
                {...form.register('name')}
              />
              <FieldError errors={form.formState.errors.name ? [form.formState.errors.name] : []} />
            </Field>
          </FieldGroup>
            
          <DialogFooter className="gap-2 sm:gap-0">
            <DialogClose asChild>
              <Button type="button" variant="ghost">
                Cancel
              </Button>
            </DialogClose>
            <Button type="submit" disabled={isPending}>
              {isPending ? 'Updating...' : 'Update Project'}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}

function DeleteProjectDialog({ project, onDeleted, open, onOpenChange }: DeleteProjectDialogProps & { open?: boolean, onOpenChange?: (open: boolean) => void }) {
  const { mutate, isPending } = useDeleteProject()

  const handleDelete = () => {
    if (!project.id) return
    mutate(project.id, { onSuccess: () => {
      onDeleted?.()
      onOpenChange?.(false)
    }})
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Delete Project</DialogTitle>
          <DialogDescription>
            This will permanently delete <strong>{project.name}</strong>.
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <Button variant="ghost" onClick={() => onOpenChange?.(false)}>Cancel</Button>
          <Button variant="destructive" onClick={handleDelete} disabled={isPending}>
            {isPending ? 'Deleting...' : 'Delete Project'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}


function RouteComponent() {
  const { data: projects, isLoading: projectsLoading } = useGetAllProjects();
  const [updateDialogOpen, setUpdateDialogOpen] = useState<string | null>(null)
  const [deleteDialogOpen, setDeleteDialogOpen] = useState<string | null>(null)
  const [createDialogOpen, setCreateDialogOpen] = useState(false)


  const currentUpdateProject = projects?.data?.find(p => p.id === updateDialogOpen)
  const currentDeleteProject = projects?.data?.find(p => p.id === deleteDialogOpen)

  if (projectsLoading) {
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
      <CreateProjectDialog onCreated={() => setCreateDialogOpen(false)} open={createDialogOpen} onOpenChange={setCreateDialogOpen} />
    </div>

    {currentUpdateProject && (
      <UpdateProjectDialog 
        project={currentUpdateProject}
        open={updateDialogOpen !== null}
        onUpdated={() => setUpdateDialogOpen(null)}
        onOpenChange={(isOpen) => {
          if (!isOpen) setUpdateDialogOpen(null)
        }}
      />
    )}

   
    {currentDeleteProject && (
      <DeleteProjectDialog 
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
        <Card 
          key={project.id || 'unknown'} 
          className="group relative flex flex-col justify-between border-input hover:border-foreground/30 transition-all duration-200 bg-card shadow-sm hover:shadow-md"
        >
          <CardHeader className="pb-4">
            <div className="flex items-start justify-between">
              <CardTitle className="text-lg font-medium tracking-tight group-hover:underline underline-offset-4">
                {project.name || 'Untitled Project'}
              </CardTitle>
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button variant="ghost" className="h-8 w-8 p-0">
                    <MoreHorizontal className="h-4 w-4" />
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end" className="w-[160px]">
                  <DropdownMenuItem 
                    onClick={() => setUpdateDialogOpen(project.id || null)}
                  >
                    Update Project
                  </DropdownMenuItem>
                  <DropdownMenuItem 
                    className="text-destructive"
                    onClick={() => setDeleteDialogOpen(project.id || null)}
                  >
                    Delete Project
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            </div>
          </CardHeader>

          <CardContent>
            <div className="flex items-center gap-2 mb-6">
              <div className="h-1.5 w-1.5 rounded-full bg-foreground shadow-[0_0_5px_rgba(0,0,0,0.2)]" />
              <span className="text-xs font-medium text-foreground">Production Active</span>
            </div>
            
            <div className="flex items-center justify-between border-t pt-4">
              <div className="text-[10px] font-mono text-muted-foreground uppercase flex gap-3">
                <span>0 ENVS</span>
                <span>•</span>
                <span>Updated 2d ago</span>
              </div>
              
              <Link 
                to="/app/projects/$projectId" 
                params={{ projectId: project.id || '' }} 
                className="text-xs font-semibold flex items-center gap-1 hover:gap-2 transition-all"
              >
                DETAILS
                <ArrowRight className="w-3 h-3" />
              </Link>
            </div>
          </CardContent>
        </Card>
      ))}
    </div>
  </div>
)
}