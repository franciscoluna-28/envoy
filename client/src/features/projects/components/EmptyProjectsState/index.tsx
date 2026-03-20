import { Card, CardContent } from '@/components/ui/card'
import { LayoutDashboard } from 'lucide-react'

export function EmptyProjectsState() {
  return (
    <Card className='my-6'>
      <CardContent className="flex flex-col items-center justify-center text-center p-12 space-y-4">
        <LayoutDashboard className="h-12 w-12 text-muted-foreground/20" />
        <div className="space-y-2">
          <h3 className="text-lg font-semibold">No projects yet</h3>
          <p className="text-muted-foreground text-sm">
            Create your first project to get started with managing your database environments.
          </p>
        </div>
      </CardContent>
    </Card>
  )
}
