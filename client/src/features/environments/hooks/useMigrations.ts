import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import client from '@/api/client'
import { toast } from 'sonner'
import { useNavigate } from '@tanstack/react-router'
import { MIGRATIONS_QUERY_KEYS } from '@/features/keys'
import type { EnvironmentId, SchemaColumn, PreviewError, TablePermission } from '@/features/types'

export function useGetEnvironmentMigrations(envId: EnvironmentId) {
  return useQuery({
    queryKey: MIGRATIONS_QUERY_KEYS.list(envId),
    queryFn: async () => {
      const response = await client.GET('/environments/{id}/migrations', {
        params: { path: { id: envId } }
      })

      if (response.error) {
        throw new Error(response.error.message || 'Failed to fetch migrations')
      }

      return response.data
    },
    enabled: !!envId,
    retry: (failureCount, error) => {
      if (error?.message?.includes('401') || error?.message?.includes('404') || error?.message?.includes('403')) {
        return false
      }
      return failureCount < 3
    }
  })
}

export function usePreviewSchemaChanges() {
  return useMutation<SchemaColumn[], PreviewError, { envId: EnvironmentId; sqlContent: string }>({
    mutationFn: async ({ envId, sqlContent }) => {
      const response = await client.POST('/environments/{id}/migrations/preview', {
        params: { path: { id: envId } },
        body: {
          sql_content: sqlContent
        }
      })
      
      if (response.error) {
        const error = new Error(response.error.message || 'Failed to preview schema changes') as PreviewError
        error.errors = Object.values(response.error.errors || {}).flat() as string[]
        throw error
      }
      
      // Don't display a toast as errors are displayed directly on the UI
      // The API returns { data: Array, message: string, success: boolean }
      return (response.data as any).data as SchemaColumn[]
    },
  })
}

export function useRunMigration() {
  const queryClient = useQueryClient()
  const navigate = useNavigate();

  return useMutation({
    mutationFn: async ({ 
      envId, 
      name, 
      description, 
      sqlContent, 
      clientId 
    }: { 
      envId: EnvironmentId; 
      name: string; 
      description: string; 
      sqlContent: string; 
      clientId: string;
    }) => {
      const response = await client.POST('/environments/{id}/migrations', {
        params: { path: { id: envId } },
        body: {
          name,
          description,
          sql_content: sqlContent,
          client_id: clientId,
          environment_id: envId
        }
      })

      if (response.error) {
  throw response.error;
      }

    if (response.data && (response.data as any).success === false) {
  throw response.data;
}
      
      return response.data
    },
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: MIGRATIONS_QUERY_KEYS.list(variables.envId) })
      queryClient.invalidateQueries({ queryKey: ['environments', 'list'] }) 
      toast.success("Migration executed successfully!")
      navigate({
        to: "/app/projects/$projectId/environments/$envId",
        params: {
          projectId: variables.clientId,
          envId: variables.envId
        }
      })
    },
    onError: (error: any) => {
      const errorMessage = error.details?.length 
        ? error.details.join('; ') 
        : error.message || 'Failed to run migration'
      
      toast.error(errorMessage)
    }
  })
}
    export function useValidateEnvironmentConnection() {
  return useMutation({
    mutationFn: async (envId: EnvironmentId) => {
      const response = await client.POST('/environments/{id}/validate', {
        params: { path: { id: envId } }
      })

      if (response.error) {
        const error = new Error(response.error.message || 'Environment validation failed')
        throw error
      }

      return response.data
    },
    onSuccess: () => {
      toast.success("Environment connection validated successfully!")
    },
    onError: (error: any) => {
      const errorMessage = error.message || 'Failed to validate environment connection'
      toast.error(errorMessage)
      console.error('Validation error:', error)
    }
  })
}

export function useTestPermissionsWithPreview() {
  return useMutation<TablePermission[], Error, { envId: EnvironmentId; databaseUser: string; sqlContent: string }>({
    mutationFn: async ({ envId, databaseUser, sqlContent }) => {
      const response = await client.POST('/environments/{id}/test-permissions-preview' as any, {
        params: { path: { id: envId } },
        body: {
          database_user: databaseUser,
          sql_content: sqlContent
        }
      } as any)
      
      if ((response as any).error) {
        throw new Error((response as any).error.message || 'Failed to test permissions with preview')
      }
      
      const results = ((response as any).data as any).data as TablePermission[]
      console.log('Permission test with preview results:', results);
      return results
    },
    onError: (error: any) => {
      const errorMessage = error.message || 'Failed to test permissions'
      toast.error(errorMessage)
      console.error('Permission test error:', error)
    }
  })
}

export function useTestPermissionsWithCurrentSchema() {
  return useMutation<TablePermission[], Error, { envId: EnvironmentId; databaseUser: string }>({
    mutationFn: async ({ envId, databaseUser }) => {
      const response = await client.POST('/environments/{id}/test-permissions-current' as any, {
        params: { path: { id: envId } },
        body: {
          database_user: databaseUser
        }
      } as any)
      
      if ((response as any).error) {
        throw new Error((response as any).error.message || 'Failed to test permissions with current schema')
      }
      
      const results = ((response as any).data as any).data as TablePermission[]
      console.log('Permission test with current schema results:', results);
      return results
    },
    onError: (error: any) => {
      const errorMessage = error.message || 'Failed to test permissions'
      toast.error(errorMessage)
      console.error('Permission test error:', error)
    }
  })
}
