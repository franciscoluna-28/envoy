import { Button } from "@/components/ui/button"
import { Database, FileText, Eye } from "lucide-react"
import { useState } from "react"
import { SchemaPreviewDialog } from "./SchemaPreviewDialog"
import { CreateMigrationDialog } from "./CreateMigrationDialog"

interface ProjectActionsProps {
  environments: any[]
}

export function ProjectActions({ environments }: ProjectActionsProps) {
  const [schemaDialogOpen, setSchemaDialogOpen] = useState(false)
  const [migrationDialogOpen, setMigrationDialogOpen] = useState(false)
  const [selectedEnvironment, setSelectedEnvironment] = useState<any>(null)

  const handlePreviewSchema = (environment: any) => {
    setSelectedEnvironment(environment)
    setSchemaDialogOpen(true)
  }

  const handleCreateMigration = (environment: any) => {
    setSelectedEnvironment(environment)
    setMigrationDialogOpen(true)
  }

  if (!environments || environments.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center py-12 text-center">
        <Database className="h-12 w-12 text-muted-foreground mb-4" />
        <h3 className="text-lg font-semibold mb-2">No environments available</h3>
        <p className="text-muted-foreground mb-4">
          Create an environment first to start managing your database schema.
        </p>
      </div>
    )
  }

  return (
    <div className="flex flex-col items-center justify-center py-12 text-center space-y-6">
      <div className="flex flex-col sm:flex-row gap-4">
        <Button
          onClick={() => handlePreviewSchema(environments[0])}
          className="flex items-center gap-2"
          variant="outline"
        >
          <Eye className="h-4 w-4" />
          Preview Schema
        </Button>
        <Button
          onClick={() => handleCreateMigration(environments[0])}
          className="flex items-center gap-2"
        >
          <FileText className="h-4 w-4" />
          Create Migration
        </Button>
      </div>
      
      {environments.length > 1 && (
        <div className="text-sm text-muted-foreground">
          {environments.length} environments available - actions use: {environments[0].name}
        </div>
      )}

      <SchemaPreviewDialog
        environmentId={selectedEnvironment?.id || ""}
        environmentName={selectedEnvironment?.name || ""}
        open={schemaDialogOpen}
        onOpenChange={setSchemaDialogOpen}
      />

      <CreateMigrationDialog
        environmentId={selectedEnvironment?.id || ""}
        open={migrationDialogOpen}
        onOpenChange={setMigrationDialogOpen}
      />
    </div>
  )
}
