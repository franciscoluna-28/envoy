import { useQuery } from '@tanstack/react-query'
import client from '@/api/client'

const ENVIRONMENT_SCHEMA_QUERY_KEYS = {
  all: ['environment-schema'] as const,
  details: () => [...ENVIRONMENT_SCHEMA_QUERY_KEYS.all, 'detail'] as const,
  detail: (envId: string) => [...ENVIRONMENT_SCHEMA_QUERY_KEYS.details(), { envId }] as const,
}

export function useEnvironmentSchema(envId: string) {
  return useQuery({
    queryKey: ENVIRONMENT_SCHEMA_QUERY_KEYS.detail(envId),
    queryFn: async () => {
      const response = await client.GET('/environments/{id}/schema', {
        params: { path: { id: envId } }
      })
      return response.data
    },
    enabled: !!envId
  })
}
