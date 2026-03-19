import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Database } from 'lucide-react'
import { Badge } from '@/components/ui/badge'
import { Link } from '@tanstack/react-router'
import type { Environment } from '@/features/types'

interface EnvironmentCardProps {
  env: Environment
}

export function EnvironmentCard({ env }: EnvironmentCardProps) {
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
                className="bg-stone-100 text-[10px] px-1.5 py-0 h-4 font-bold border-none text-stone-500 uppercase tracking-wider"
              >
                {env.type || 'development'}
              </Badge>
              
              <div className="flex items-center gap-1.5 text-[11px] font-medium text-green-600">
                <div className="h-1.5 w-1.5 rounded-full bg-green-500 animate-pulse shadow-[0_0_8px_rgba(34,197,94,0.4)]" />
                Connected
              </div>
            </div>
          </div>
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