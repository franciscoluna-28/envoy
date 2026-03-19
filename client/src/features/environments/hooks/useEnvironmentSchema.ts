import { useQuery } from '@tanstack/react-query'
import client from '@/api/client'
import { ENVIRONMENT_SCHEMA_QUERY_KEYS } from '@/features/keys'
import type { EnvironmentId, DatabaseSchemaItem } from '@/features/types'

export function useEnvironmentSchema(envId: EnvironmentId) {
  return useQuery({
    queryKey: ENVIRONMENT_SCHEMA_QUERY_KEYS.detail(envId),
    queryFn: async () => {
      const response = await client.GET('/environments/{id}/schema', {
        params: { path: { id: envId } }
      })
      return response.data as unknown as DatabaseSchemaItem[]
    },
    enabled: !!envId
  })
}
