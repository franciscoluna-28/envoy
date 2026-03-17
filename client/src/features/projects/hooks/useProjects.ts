import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import client from '@/api/client'

const PROJECTS_QUERY_KEYS = {
  all: ['projects'] as const,
  lists: () => [...PROJECTS_QUERY_KEYS.all, 'list'] as const,
  list: (filter: string) => [...PROJECTS_QUERY_KEYS.lists(), { filter }] as const,
  details: () => [...PROJECTS_QUERY_KEYS.all, 'detail'] as const,
  detail: (id: string) => [...PROJECTS_QUERY_KEYS.details(), id] as const,
}

export function useGetAllProjects() {
  return useQuery({
    queryKey: PROJECTS_QUERY_KEYS.lists(),
    queryFn: async () => {
      const response = await client.GET('/projects')
      return response.data
    }
  })
}

export function useGetProject(projectId: string) {
  return useQuery({
    queryKey: PROJECTS_QUERY_KEYS.detail(projectId),
    queryFn: async () => {
      const response = await client.GET('/projects/{id}', {
        params: { path: { id: projectId } }
      })
      return response.data
    },
    enabled: !!projectId
  })
}

export function useCreateProject() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (data: { name: string }) => {
      const response = await client.POST('/projects', {
        body: data
      })
      return response.data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: PROJECTS_QUERY_KEYS.lists() })
    },
  })
}

export function useUpdateProject() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async ({ id, data }: { id: string; data: { name: string } }) => {
      const response = await client.PUT('/projects/{id}', {
        params: { path: { id } },
        body: { id, ...data }
      })
      return response.data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: PROJECTS_QUERY_KEYS.lists() })
      queryClient.invalidateQueries({ queryKey: PROJECTS_QUERY_KEYS.details() })
    },
  })
}

export function useDeleteProject() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (id: string) => {
      const response = await client.DELETE('/projects/{id}', {
        params: { path: { id } }
      })
      return response.data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: PROJECTS_QUERY_KEYS.lists() })
      queryClient.invalidateQueries({ queryKey: PROJECTS_QUERY_KEYS.details() })
    },
  })
}
