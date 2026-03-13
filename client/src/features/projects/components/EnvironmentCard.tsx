import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Database, ExternalLink } from 'lucide-react'
import { Badge } from '@/components/ui/badge'
import { Link } from '@tanstack/react-router'

interface EnvironmentCardProps {
  env: {
    id?: string
    name?: string
    connection_status?: string
  }
  projectId: string
}

export function EnvironmentCard({ env, projectId }: EnvironmentCardProps) {
  return (
    <Card className="hover:border-primary/50 transition-colors">
      <CardHeader className="pb-2">
        <div className="flex justify-between items-start">
          <CardTitle className="text-lg font-bold flex items-center gap-2">
            <Database className="h-4 w-4 text-muted-foreground" />
            {env.name}
          </CardTitle>
          <Badge variant={env.connection_status === 'connected' ? 'default' : 'secondary'}>
            {env.connection_status || 'Unknown'}
          </Badge>
        </div>
        <div className="font-mono text-xs truncate text-muted-foreground">
          Database Environment
        </div>
      </CardHeader>
      <CardContent>
        <div className="flex items-center gap-2 text-sm">
          <div className="h-2 w-2 rounded-full bg-green-500 animate-pulse" />
          <span className="text-muted-foreground">Connected</span>
        </div>
      </CardContent>
      <div className="border-t pt-4 bg-muted/10 px-6 pb-4">
        <Button variant="ghost" className="w-full justify-between group">
          <Link to={`/app/projects/${projectId}/environments/${env.id}`} className="flex items-center justify-between w-full">
            <span>View Migrations</span>
            <ExternalLink className="h-4 w-4 opacity-0 group-hover:opacity-100 transition-opacity" />
          </Link>
        </Button>
      </div>
    </Card>
  )
}
