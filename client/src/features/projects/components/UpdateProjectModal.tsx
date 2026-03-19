import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Field, FieldGroup, FieldLabel, FieldError } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { useUpdateProject } from '../hooks/useProjects'
import type { Project } from '@/features/types'
import { updateProjectSchema } from '@/features/schemas'

type UpdateProjectForm = z.infer<typeof updateProjectSchema>

interface UpdateProjectModalProps {
  project: Pick<Project, 'id' | 'name'>
  onUpdated?: () => void
  open?: boolean
  onOpenChange?: (open: boolean) => void
}

export function UpdateProjectModal({ project, onUpdated, open, onOpenChange }: UpdateProjectModalProps) {
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
            <Button type="button" variant="ghost" onClick={() => onOpenChange?.(false)}>
              Cancel
            </Button>
            <Button type="submit" disabled={isPending}>
              {isPending ? 'Updating...' : 'Update Project'}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
