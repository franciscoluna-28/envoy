import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import client from '@/api/client'
import { toast } from 'sonner'

const ENVIRONMENTS_QUERY_KEYS = {
  all: ['environments'] as const,
  lists: () => [...ENVIRONMENTS_QUERY_KEYS.all, 'list'] as const,
  list: (projectId: string) => [...ENVIRONMENTS_QUERY_KEYS.lists(), { projectId }] as const,
}

export function useGetEnvironments(projectId: string) {
  return useQuery({
    queryKey: ENVIRONMENTS_QUERY_KEYS.list(projectId),
    queryFn: async () => {
      const response = await client.GET('/projects/{id}/environments', {
        params: { path: { id: projectId } }
      })
      return response.data
    },
    enabled: !!projectId
  })
}

export function useGetEnvironment(projectId: string, envId: string) {
  return useQuery({
    queryKey: [...ENVIRONMENTS_QUERY_KEYS.list(projectId), envId],
    queryFn: async () => {
      const response = await client.GET('/environments/{id}', {
        params: { path: { id: envId } }
      })
      return response.data
    },
    enabled: !!projectId && !!envId
    } )
}

export function useCreateEnvironment() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (data: { projectId: string; name: string; connection_string: string }) => {
      const response = await client.POST('/projects/{id}/environments', {
        params: { path: { id: data.projectId } },
        body: {
          name: data.name,
          connection_url: data.connection_string,
          project_id: data.projectId,
          type: data.type
        }
      })
      return response.data
    },
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ENVIRONMENTS_QUERY_KEYS.list(variables.projectId) })
      toast.success("Environment created successfully!")
    },
    onError: () => {
      toast.error("Failed to create environment")
    }
  })
}
