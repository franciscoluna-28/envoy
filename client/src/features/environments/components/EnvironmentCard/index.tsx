import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Calendar, Database, MoreHorizontal } from 'lucide-react'
import { Badge } from '@/components/ui/badge'
import { Link } from '@tanstack/react-router'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import type { Environment } from '@/features/types'
import { formatDate } from '@/utils/date'

interface EnvironmentCardProps {
  env: Environment
  onUpdate?: () => void
  onDelete?: () => void 
}

export function EnvironmentCard({ env, onUpdate }: EnvironmentCardProps) {
  return (
    <Card className="group relative min-w-[350px] shadow-sm hover:shadow-md transition-all duration-300 border-stone-200 p-0 overflow-hidden">
      <CardContent className="p-6">
        <div className="flex items-start justify-between mb-8">
          <div className="space-y-1.5">
            <h3 className="text-lg font-semibold tracking-tight text-stone-900 flex items-center gap-2">
              <Database className="w-4 h-4 text-stone-400 transition-colors" />
              {env.name || "Unnamed Environment"}
            </h3>
            
            <div className="flex items-center gap-2">
              <Badge 
                variant="secondary" 
                
              >
                {env.type || 'development'}
              </Badge>
              <span className="text-stone-300 text-[10px]">•</span>
               <div className="flex items-center gap-1 text-xs">
                <Calendar className="w-3 h-3 text-muted-foreground" />
                <span>
                  {env.created_at
                    ? formatDate(env.created_at)
                    : "Just now"}
                </span>
              </div>
            </div>
          </div>

          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button
                variant="ghost"
                className="h-8 w-8 p-0 hover:bg-stone-100 text-stone-400 data-[state=open]:bg-stone-100 data-[state=open]:text-stone-900"
              >
                <MoreHorizontal className="h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>

            <DropdownMenuContent
              align="end"
              className="w-[180px] p-1.5 shadow-xl border-stone-200"
            >
              <DropdownMenuItem
                onClick={onUpdate}
                className="text-xs font-semibold cursor-pointer px-2.5 py-2 rounded-md transition-colors"
              >
                Update Environment
              </DropdownMenuItem>

            {/*   <Separator className="my-1.5" />
 */}
          {/*     <DropdownMenuItem
                className="text-destructive text-xs font-semibold cursor-pointer px-2.5 py-2 rounded-md hover:bg-red-50 hover:text-red-600 transition-all"
                onClick={onDelete}
              >
                Delete Environment
              </DropdownMenuItem> */}
            </DropdownMenuContent>
          </DropdownMenu>
        </div>

        <div className="flex items-center justify-between gap-3">
          <Link 
            className="w-full"
            to="/app/projects/$projectId/environments/$envId"
            params={{
              projectId: env.project_id || '',
              envId: env.id || ''
            }}
          >
            <Button 
              variant="secondary" 
              size="lg"
              className="w-full"
            >
              Configure Environment
            </Button>
          </Link>
        </div>
      </CardContent>
    </Card>
  )
}