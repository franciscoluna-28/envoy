import { Card, CardContent } from '@/components/ui/card'
import { Database } from 'lucide-react'

export function EmptyEnvironmentsState() {
  return (
    <Card className="border-dashed border-2 col-span-full">
      <CardContent className="flex flex-col items-center justify-center text-center p-12 space-y-4">
        <Database className="h-12 w-12 text-muted-foreground/50" />
        <div className="space-y-2">
          <h3 className="text-lg font-semibold">No environments yet</h3>
          <p className="text-muted-foreground text-sm">
            Create your first database environment to get started with managing connections and migrations.
          </p>
        </div>
      </CardContent>
    </Card>
  )
}
