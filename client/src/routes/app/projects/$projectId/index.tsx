import { createFileRoute } from '@tanstack/react-router'
import { useQuery } from '@tanstack/react-query'
import client from '@/api/client'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'

export const Route = createFileRoute('/app/projects/$projectId/')({
  component: RouteComponent,
})

function RouteComponent() {
  const { projectId } = Route.useParams()

  const { data: project, isLoading, error } = useQuery({
    queryKey: ['project', projectId],
    queryFn: async () => {
      const response = await client.GET('/api/v1/projects/{id}', {
        params: {
          path: { id: projectId }
        }
      })
      return response.data
    }
  })

  if (isLoading) {
    return (
      <div className="flex items-center justify-center p-8">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
      </div>
    )
  }

  if (error) {
    return (
      <Card className="w-full max-w-md">
        <CardHeader>
          <CardTitle className="text-xl text-destructive">Error</CardTitle>
          <CardDescription>
            Failed to load project
          </CardDescription>
        </CardHeader>
        <CardContent>
          <p className="text-sm text-muted-foreground">
            {error.message || 'An unknown error occurred'}
          </p>
        </CardContent>
      </Card>
    )
  }

  if (!project?.data) {
    return (
      <Card className="w-full max-w-md">
        <CardHeader>
          <CardTitle>Project Not Found</CardTitle>
          <CardDescription>
            The requested project could not be found.
          </CardDescription>
        </CardHeader>
      </Card>
    )
  }

  return (
    <Card className="w-full max-w-2xl">
      <CardHeader>
        <CardTitle className="text-2xl">{project.data.name}</CardTitle>
        <CardDescription>Project Details</CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        <div>
          <h3 className="text-lg font-semibold mb-2">Description</h3>
          <p className="text-muted-foreground">
            Description will be available after API types are regenerated
          </p>
        </div>
        
        <div>
          <h3 className="text-lg font-semibold mb-2">Project Information</h3>
          <div className="space-y-2">
            <div className="flex justify-between">
              <span className="text-sm font-medium">Project ID:</span>
              <span className="text-sm font-mono">{project.data.id}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-sm font-medium">Project Name:</span>
              <span className="text-sm">{project.data.name}</span>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}
