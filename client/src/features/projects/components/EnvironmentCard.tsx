import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Database } from 'lucide-react'
import { Badge } from '@/components/ui/badge'
import { useState } from 'react'
import { SchemaPreviewDialog } from './SchemaPreviewDialog'
import { CreateMigrationDialog } from './CreateMigrationDialog'
import { Link } from '@tanstack/react-router'

interface EnvironmentCardProps {
  env: {
    id?: string
    name?: string
    connection_status?: string
    type?: string
    project_id?: string
  }
}

export function EnvironmentCard({ env }: EnvironmentCardProps) {
  const [schemaDialogOpen, setSchemaDialogOpen] = useState(false)
  const [migrationDialogOpen, setMigrationDialogOpen] = useState(false)

  return (
    <Card className="hover:border-primary/50 transition-colors min-w-[300px]">
      <CardHeader className="pb-2">
        <div className="flex justify-between items-start">
          <CardTitle className="text-lg font-bold flex items-center gap-2">
            <Database className="h-4 w-4 text-muted-foreground" />
            {env.name}
          </CardTitle>
          <Badge className='capitalize' variant="secondary">
            {env.type || 'development'}
          </Badge>
        </div>
        <div className="font-mono text-xs truncate text-muted-foreground">
          Database Environment
        </div>
      </CardHeader>
      <CardContent>
        <div className="flex items-center gap-2 text-sm">
          <div className="h-2 w-2 rounded-full bg-green-500 animate-pulse" />
          <span className="text-black font-medium text-xs">Connected</span>
        </div>
      </CardContent>
      <div className="border-t pt-4 bg-muted/10 px-6 pb-4 space-y-2">
      <Link to={`/app/projects/${env.project_id}/environments/${env.id}`}>
      <Button variant="secondary" className="w-full gap-2">
        View Details
      </Button>
      </Link>
      {/*   <Button 
          variant="outline" 
          className="w-full justify-start gap-2"
          onClick={() => setSchemaDialogOpen(true)}
        >
          <Eye className="h-4 w-4" />
          Preview Schema
        </Button>
        <Button 
          className="w-full justify-start gap-2"
          onClick={() => setMigrationDialogOpen(true)}
        >
          <FileText className="h-4 w-4" />
          Create Migration
        </Button> */}
      </div>
      
      <SchemaPreviewDialog
        environmentId={env.id || ""}
        environmentName={env.name || ""}
        open={schemaDialogOpen}
        onOpenChange={setSchemaDialogOpen}
      />

      <CreateMigrationDialog
        environmentId={env.id || ""}
        open={migrationDialogOpen}
        onOpenChange={setMigrationDialogOpen}
      />
    </Card>
  )
}
