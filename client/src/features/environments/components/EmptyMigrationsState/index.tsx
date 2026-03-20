import { Card, CardContent } from '@/components/ui/card'
import { Server } from 'lucide-react'

export function EmptyMigrationsState() {
  return (
    <Card className="overflow-hidden">
      <CardContent className="flex flex-col items-center justify-center text-center p-12 space-y-4">
        <Server className="h-12 w-12 text-muted-foreground/20" />
        <div className="space-y-2">
          <h3 className="text-lg font-semibold">No migrations yet</h3>
          <p className="text-muted-foreground text-sm">
            This environment doesn't have any migrations yet. 
          </p>
        </div>
      </CardContent>
    </Card>
  )
}
