import { Button } from '@/components/ui/button'
import { Spinner } from '@/components/ui/spinner'
import { ShadcnReactFlowSchema, useEnvironmentSchema, useGetProject } from '@/features/projects'
import { useGetEnvironment } from '@/features/projects/hooks/useEnvironments'
import { createFileRoute, Link } from '@tanstack/react-router'
import { ChevronRight, Database } from 'lucide-react'
import { Separator } from '@/components/ui/separator'
import { Badge } from '@/components/ui/badge'

export const Route = createFileRoute(
  '/app/projects/$projectId/environments/$envId/',
)({
  component: RouteComponent,
})

function RouteComponent() {
  const params = Route.useParams()
  const { projectId, envId } = params

  const { data: projectData, isPending: isLoadingProject } = useGetProject(projectId)
  const { data: environmentData, isPending: isLoadingEnvironment } = useGetEnvironment(projectId, envId)
  const { data: environmentSchema, isPending: isLoadingEnvironmentSchema } = useEnvironmentSchema(envId)

  if (isLoadingProject || isLoadingEnvironment || isLoadingEnvironmentSchema) {
    return (
      <div className="flex h-[80vh] items-center justify-center">
        <Spinner className="h-8 w-8 text-black" />
      </div>
    )
  }

  return (
    <div className="flex flex-col h-full bg-white text-black font-sans antialiased">
      <header className="border-b px-6 py-4">
        <div className="flex items-center gap-2 text-xs text-muted-foreground mb-2">
          <Link to="/app" className="hover:underline">Projects</Link>
          <ChevronRight className="h-3 w-3" />
          <span className="font-medium text-black">{projectData?.name}</span>
          <ChevronRight className="h-3 w-3" />
          <span className="font-mono">{environmentData?.name}</span>
        </div>
        
        <div className="flex justify-between items-end">
          <div>
            <h1 className="text-2xl font-bold tracking-tight flex items-center gap-2">
              <Database className="h-5 w-5" />
              {environmentData?.name}
            </h1>
            <p className="text-xs font-mono text-muted-foreground mt-1 uppercase tracking-widest">
              ID: {envId}
            </p>
          </div>
          
          <div className="flex gap-2">
           
            <Button>
              
              Create Migration
            </Button>
          </div>
        </div>
      </header>

      <main className="flex-1 flex flex-col p-6 gap-6">
        
        <div className="grid grid-cols-3 gap-4">
          {[
            { label: 'Status', value: 'Connected', status: 'bg-black' },
            { label: 'Type', value: environmentData?.type || 'Development', status: 'bg-stone-200' },
            { label: 'Database', value: 'Postgres', status: 'bg-stone-100' },
          ].map((stat) => (
            <div key={stat.label} className="border p-3 rounded-md flex flex-col gap-1">
              <span className="text-[10px] font-semibold text-muted-foreground capitalize">{stat.label}</span>
              <span className="text-sm font-semibold tracking-tighter capitalize">{stat.value}</span>
            </div>
          ))}
        </div>

        <Separator className="bg-stone-100" />

        <div className="flex-1 min-h-[500px] relative bg-stone-50/30">
          <div className="absolute top-4 left-4 z-10">
            <Badge variant="outline" className="bg-white border-black text-[10px] font-mono shadow-sm">
              Schema Preview
            </Badge>
          </div>
          <ShadcnReactFlowSchema 
            schema={environmentSchema} 
            isLoading={isLoadingEnvironmentSchema} 
          />
        </div>
      </main>
    </div>
  )
}

