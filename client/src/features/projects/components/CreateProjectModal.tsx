import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog'
import { Field, FieldGroup, FieldLabel, FieldDescription, FieldError } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import { DialogClose } from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Plus } from 'lucide-react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { useCreateProject } from '../hooks/useProjects'

const createProjectSchema = z.object({
  name: z.string().min(1, 'Project name is required').max(100, 'Project name must be less than 100 characters')
})

type CreateProjectForm = z.infer<typeof createProjectSchema>

interface CreateProjectModalProps {
  onCreated?: () => void
  open?: boolean
  onOpenChange?: (open: boolean) => void
}

export function CreateProjectModal({ onCreated, open, onOpenChange }: CreateProjectModalProps) {
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
        onOpenChange?.(false)
      }
    })
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogTrigger asChild>
        <Button>
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
