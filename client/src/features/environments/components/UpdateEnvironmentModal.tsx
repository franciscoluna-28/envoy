import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Field, FieldGroup, FieldLabel, FieldError } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { useUpdateEnvironment } from '../hooks/useEnvironments'
import type { Environment } from '@/features/types'
import { updateEnvironmentSchema } from '@/features/schemas'

type UpdateEnvironmentForm = z.infer<typeof updateEnvironmentSchema>

interface UpdateEnvironmentModalProps {
  environment: Pick<Environment, 'id' | 'name'>
  onUpdated?: () => void
  open?: boolean
  onOpenChange?: (open: boolean) => void
}

export function UpdateEnvironmentModal({ environment, onUpdated, open, onOpenChange }: UpdateEnvironmentModalProps) {
  const form = useForm<UpdateEnvironmentForm>({
    resolver: zodResolver(updateEnvironmentSchema),
    defaultValues: {
      name: environment.name || ''
    }
  })

  const { mutate, isPending } = useUpdateEnvironment()

  const onSubmit = (data: UpdateEnvironmentForm) => {
    if (!environment.id) return
    mutate({ id: environment.id as string, name: data.name }, {
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
          <DialogTitle>Update Environment</DialogTitle>
          <DialogDescription>
            Update the environment name.
          </DialogDescription>
        </DialogHeader>
        
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
          <FieldGroup>
            <Field>
              <FieldLabel htmlFor="name">Environment Name</FieldLabel>
              <Input
                id="name"
                placeholder="e.g. Development Database"
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
              {isPending ? 'Updating...' : 'Update Environment'}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
