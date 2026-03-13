import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog'
import { Field, FieldGroup, FieldLabel, FieldDescription, FieldError } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import { DialogClose } from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Plus } from 'lucide-react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { useCreateEnvironment } from '../hooks/useEnvironments'

const createEnvironmentSchema = z.object({
  name: z.string().min(1, 'Environment name is required').max(100, 'Environment name must be less than 100 characters'),
  connection_string: z.string().min(1, 'Connection string is required')
})

type CreateEnvironmentForm = z.infer<typeof createEnvironmentSchema>

interface CreateEnvironmentModalProps {
  projectId: string
  onCreated?: () => void
  open?: boolean
  onOpenChange?: (open: boolean) => void
}

export function CreateEnvironmentModal({ projectId, onCreated, open, onOpenChange }: CreateEnvironmentModalProps) {
  const form = useForm<CreateEnvironmentForm>({
    resolver: zodResolver(createEnvironmentSchema),
    defaultValues: {
      name: '',
      connection_string: ''
    }
  })

  const { mutate, isPending } = useCreateEnvironment()

  const onSubmit = (data: CreateEnvironmentForm) => {
    mutate({ projectId, ...data }, {
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
          <Plus className="mr-2 h-4 w-4" />
          New Environment
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Create New Environment</DialogTitle>
          <DialogDescription>
            Add a new database environment to your project.
          </DialogDescription>
        </DialogHeader>
        
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
          <FieldGroup>
            <Field>
              <FieldLabel htmlFor="name">Environment Name</FieldLabel>
              <Input
                id="name"
                placeholder="e.g. Production Database"
                {...form.register('name')}
              />
              <FieldDescription>
                This will be used to identify your database environment.
              </FieldDescription>
              <FieldError errors={form.formState.errors.name ? [form.formState.errors.name] : []} />
            </Field>
            
            <Field>
              <FieldLabel htmlFor="connection_string">Connection String</FieldLabel>
              <Input
                id="connection_string"
                placeholder="postgresql://user:password@localhost:5432/dbname?sslmode=require"
                {...form.register('connection_string')}
              />
              <FieldDescription>
                The database connection string. Include SSL mode in the URL parameters.
              </FieldDescription>
              <FieldError errors={form.formState.errors.connection_string ? [form.formState.errors.connection_string] : []} />
            </Field>
          </FieldGroup>
            
          <DialogFooter className="gap-2 sm:gap-0">
            <DialogClose asChild>
              <Button type="button" variant="ghost">
                Cancel
              </Button>
            </DialogClose>
            <Button type="submit" disabled={isPending}>
              {isPending ? 'Creating...' : 'Create Environment'}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
