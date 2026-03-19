import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { useDeleteProject } from '../hooks/useProjects'
import type { Project } from '@/features/types'

interface DeleteProjectModalProps {
  project: Pick<Project, 'id' | 'name'>
  onDeleted?: () => void
  open?: boolean
  onOpenChange?: (open: boolean) => void
}

export function DeleteProjectModal({ project, onDeleted, open, onOpenChange }: DeleteProjectModalProps) {
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
